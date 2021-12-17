package wallpapers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	wp "github.com/reujab/wallpaper"
)

var currentWallpaperFile string

func init() {
	currentWallpaperFile = filepath.Join(wallpaperDir, "current-url.txt")
}

func SetFromURL(Url string) error {
	err := os.WriteFile(currentWallpaperFile, []byte(Url), 0)
	if err != nil {
		return err
	}

	path, err := AddUrl(Url)
	if err != nil {
		return err
	}
	return wp.SetFromFile(path)
}

func CurrentFilePath() (string, error) {
	return wp.Get()
}

func CurrentURL() (string, error) {
	bytes, err := os.ReadFile(currentWallpaperFile)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func init() {
	must(wp.SetMode(wp.Crop))
	// Should be pretty random
	rand.Seed(int64(time.Now().Local().Hour()*time.Now().Nanosecond() + time.Now().Day()*time.Now().Hour() + time.Now().Second()))
	go func() {
		// Load data!
		resp, err := http.Get("https://api.github.com/gists/" + GistId)
		must(err)
		bytes, err := ioutil.ReadAll(resp.Body)
		must(err)
		var gist gistResponse
		must(json.Unmarshal(bytes, &gist))
		var content []byte
		for _, file := range gist.Files {
			if file.Language == "JSON" {
				if content != nil {
					log.Fatal("Gist " + GistId + " has 2+ JSON files, unable to find proper data.")
				}
				content = []byte(file.Content)
			}
		}
		must(json.Unmarshal(content, &OnlineWallpapers))
		// Wallpaper data is now loaded.
		close(WallpaperDataIsValid)
	}()
}
