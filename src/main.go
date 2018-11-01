package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Core core info
type Core struct {
	Name string `json:"name" binding:"required"`
	File string `json:"file" binding:"required"`
	URL  string `json:"url" binding:"required"`
}

var client *github.Client
var ctx = context.Background()

func main() {

	token := flag.String("t", "", "Github Personal API Token")
	output := flag.String("o", "/tmp/mister", "Output Directory")
	flag.Parse()

	if *token == "" {
		fmt.Println("Github Personal API Token Required, use -t <API Token>")
		os.Exit(-1)
	}

	// Setup Github Client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)

	// Setup Output Directory
	os.MkdirAll(*output, os.ModePerm)

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	api := router.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		api.GET("/cores", GetCoreList)
	}

	router.Run(":3000")
}

// GetCoreList return list of latest cores.
func GetCoreList(c *gin.Context) {
	response := make([]Core, 0)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 500},
	}

	org := "MiSTer-devel"
	path := "releases"

	repos, _, err := client.Repositories.ListByOrg(ctx, org, opt)

	if err != nil {
		panic(err)
	}

	for _, repo := range repos {
		_, dir, _, err := client.Repositories.GetContents(ctx, org, *repo.Name, path, nil)

		if err != nil {
			continue
		}

		hasCoresTest := func(s string) bool { return strings.HasSuffix(s, ".rbf") }
		cores := FilterCores(dir, hasCoresTest)

		sort.Slice(cores, func(i, j int) bool {
			return strings.ToLower(*cores[i].Name) < strings.ToLower(*cores[j].Name)
		})

		if len(cores) < 1 {
			continue
		}

		latest := cores[len(cores)-1]

		core := Core{
			Name: *repo.Name,
			File: *latest.Name,
			URL:  *latest.DownloadURL,
		}

		response = append(response, core)
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
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
