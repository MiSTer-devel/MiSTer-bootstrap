package main

import "C"

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Repo is the default repo structure.
type Repo struct {
	File	string `json:"file"`
	Name	string `json:"name"`
	URL		string `json:"url"`
}

func main() {
	Update()
}

//export Update
func Update() {
	repo := flag.String("r", "https://raw.githubusercontent.com/OpenVGS/MiSTer-repository/master/repo.json", "Repo URL")
	output := flag.String("o", ".", "Output Directory")

	flag.Parse()

	err := os.MkdirAll(*output, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating output directory:", err)
		return
	}

	req, err := http.NewRequest("GET", *repo, nil)
	if err != nil {
		log.Fatal("Building repo request:", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Repo request:", err)
		return
	}

	defer Close(resp.Body)

	var repos []Repo

	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		log.Println(err)
	}

	for _, repo := range repos {

		if _, err := os.Stat(fmt.Sprintf("%s/%s", *output, repo.File)); os.IsNotExist(err) {
			log.Printf("Downloading updated %s core...\n", repo.Name)
			err := DownloadCore(fmt.Sprintf("%s/%s", *output, repo.File), repo.URL)

			if err != nil {
				panic(err)
			}
		} else {
			log.Printf("Core %s already latest version, skipping\n", repo.Name)
		}
	}
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
	defer Close(out)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer Close(resp.Body)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Close is a generic io Closer with error handling
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal("IO close error:", err)
	}
}
