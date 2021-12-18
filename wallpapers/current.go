package wallpapers

import (
	"encoding/json"
	"errors"
	"io/fs"
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
func must(err error) {
	if err != nil {
		panic(err)
	}
}

type gistResponse struct {
	Files map[string]struct {
		Language string `json:"language"`
		Content  string `json:"content"`
	} `json:"files"`
}

const (
	// Gist to look for a JSON file to parse as info.
	gistId = "6c1acc57b33cdb88b720637d3d4d2af5"
)

var urlDataPath string

func init() {
	urlDataPath = filepath.Join(wallpaperDir, "urldata.json")
	must(wp.SetMode(wp.Crop))
	// Should be pretty random
	rand.Seed(int64(time.Now().Local().Hour()*time.Now().Nanosecond() + time.Now().Day()*time.Now().Hour() + time.Now().Second()))
	go func() {
		var content []byte
		urldataInfo, err := os.Stat(urlDataPath)
		// if the file doesnt exist or wasnt modified in the last 30 mins, fetch the data
		if errors.Is(err, fs.ErrNotExist) || time.Since(urldataInfo.ModTime()).Minutes() > 30 {
			// Load data!
			resp, err := http.Get("https://api.github.com/gists/" + gistId)
			must(err)
			bytes, err := ioutil.ReadAll(resp.Body)
			must(err)
			var gist gistResponse
			must(json.Unmarshal(bytes, &gist))
			for _, file := range gist.Files {
				if file.Language == "JSON" {
					if content != nil {
						log.Fatal("Gist " + gistId + " has 2+ JSON files, unable to find proper data.")
					}
					content = []byte(file.Content)
				}
			}
			os.WriteFile(urlDataPath, content, 0)
		} else if err != nil {
			panic(err)
		} else {
			content, err = os.ReadFile(urlDataPath)
			if err != nil {
				panic(err)
			}
		}
		must(json.Unmarshal(content, &OnlineWallpapers))
		// Wallpaper data is now loaded.
		close(WallpaperDataIsValid)
	}()
}
