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
	"runtime"
	"sort"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	org            = "MiSTer-devel"
	releasedirpath = "releases"
)

func getDefaultPath() string {
	home := ""
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}

	} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		home = "/tmp"
	}
	return path.Join(home, "mister")
}

func main() {
	token := flag.String("t", "", "Github Personal API Token")
	output := flag.String("o", getDefaultPath(), "Output Directory")

	flag.Parse()

	os.MkdirAll(*output, os.ModePerm)

	ctx := context.Background()

	client := github.NewClient(nil)

	if *token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 500},
	}

	repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)

	if err != nil {
		panic(err)
	}

	hasCoresTest := func(s string) bool { return strings.HasSuffix(s, ".rbf") }

	for _, repo := range repos {
		_, dir, _, err := client.Repositories.GetContents(ctx, org, *repo.Name, releasedirpath, nil)

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

		if _, err := os.Stat(fmt.Sprintf("%s/%s", *output, *latest.Name)); os.IsNotExist(err) {
			log.Printf("Downloading %s core...\n", *latest.Name)
			err := DownloadCore(fmt.Sprintf("%s/%s", *output, *latest.Name), *latest.DownloadURL)

			if err != nil {
				panic(err)
			}
		} else {
			log.Printf("Core %s already downloaded, skipping\n", *latest.Name)
		}
	}
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
