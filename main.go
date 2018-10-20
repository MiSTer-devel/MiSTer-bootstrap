package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	token := flag.String("t", "", "Github Personal API Token")
	output := flag.String("o", "/tmp/mister", "Output Directory")

	flag.Parse()

	os.MkdirAll(*output, os.ModePerm)

	dbPath := path.Join(*output, "cores.db")

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("CoreInfo"))
		if err != nil {
			return fmt.Errorf("Create bolt bucket failed: %s", err)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Could not create bucket: %s", err)
	}

	if *token == "" {
		fmt.Println("Github Personal API Token Required, use -t <API Token>")
		os.Exit(-1)
	}

	org := "MiSTer-devel"
	path := "releases"

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 500},
	}

	repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)

	if err != nil {
		panic(err)
	}

	hasCoresTest := func(s string) bool { return strings.HasSuffix(s, ".rbf") }

	for _, repo := range repos {
		_, dir, _, err := client.Repositories.GetContents(ctx, org, *repo.Name, path, nil)

		if err != nil {
			continue
		}

		cores := FilterCores(dir, hasCoresTest)

		sort.Slice(cores, func(i, j int) bool {
			return strings.ToLower(*cores[i].Name) < strings.ToLower(*cores[j].Name)
		})

		if len(cores) < 1 {
			continue
		}

		latest := cores[len(cores)-1]

		sha, _ := GetLatestCoreInfo(db, repo)

		if sha != nil {
			fmt.Printf("Latest SHA found for %s: %s\n", *repo.Name, sha)
		}

		if string(sha) == *latest.SHA {
			fmt.Printf("Already have latest version for %s: %s\n", *repo.Name, *latest.SHA)
		} else {
			fmt.Printf("Found newer version for %s: %s\n", *repo.Name, *latest.SHA)

			err := DownloadCore(fmt.Sprintf("%s/%s", *output, *latest.Name), *latest.DownloadURL)

			if err != nil {
				panic(err)
			}

			err = SaveCoreInfo(db, repo, latest)

			if err != nil {
				panic(err)
			}
		}
	}
}

// SaveCoreInfo save core info to kv store.
func SaveCoreInfo(db *bolt.DB, repo *github.Repository, core *github.RepositoryContent) error {
	tx, err := db.Begin(false)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = tx.Bucket([]byte("CoreInfo")).Put([]byte(*repo.Name), []byte(*core.SHA))
	if err != nil {
		return fmt.Errorf("Put error: %s", err)
	}
	return nil
}

// GetLatestCoreInfo get core info from kv store.
func GetLatestCoreInfo(db *bolt.DB, repo *github.Repository) (sha []byte, err error) {
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return tx.Bucket([]byte("CoreInfo")).Get([]byte(*repo.Name)), nil
}

// FilterCores applies a test function to each element in the list of cores and returns passing cores.
func FilterCores(cores []*github.RepositoryContent, test func(string) bool) (ret []*github.RepositoryContent) {
	for _, core := range cores {
		if test(*core.Name) {
			ret = append(ret, core)
		}
	}
	return
}

// DownloadCore downloads core to local filesystem.
func DownloadCore(path string, url string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	out, err := os.Create(abs)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
