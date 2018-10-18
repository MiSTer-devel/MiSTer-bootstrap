package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

func main() {
	token := flag.String("t", "", "Github Personal API Token")
	output := flag.String("o", "~/MiSTer-bootstrap", "Output Directory")

	flag.Parse()

	os.MkdirAll(*output, os.ModePerm)

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

		fmt.Println(*latest.Name)

		err = DownloadCore(fmt.Sprintf("%s/%s", *output, *latest.Name), *latest.DownloadURL)

		if err != nil {
			panic(err)
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
func DownloadCore(filepath string, url string) error {
	out, err := os.Create(filepath)
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
