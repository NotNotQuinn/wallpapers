package wallpapers

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// Maps playlist names to lists of URLs
type mappedPlaylists map[string]([]string)

// Playlists created by user
var Playlists = make(mappedPlaylists)
var playlistsFile string

func init() {
	playlistsFile = filepath.Join(wallpaperDir, "playlists.json")
}

var playlistsLoaded bool = false

func LoadPlaylists() (mappedPlaylists, error) {
	if playlistsLoaded {
		return Playlists, nil
	}

	bytes, err := os.ReadFile(playlistsFile)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		file, err := os.Create(playlistsFile)
		if err != nil {
			return nil, err
		}
		file.Close()
		bytes = []byte("")
	}

	err = json.Unmarshal(bytes, &Playlists)
	if err != nil {
		return nil, err
	}

	playlistsLoaded = true
	return Playlists, nil
}

// remove duplicates from playlists, and remove empty playlists
func SanitizePlaylists() {
	var tmp = make(mappedPlaylists)
	for key, list := range Playlists {
		// filter duplicates
		var seen = make(map[string]bool)
		for _, value := range list {
			seen[value] = true
		}
		var tmp2 []string
		for value := range seen {
			tmp2 = append(tmp2, value)
		}
		list = tmp2

		// remove empty
		if len(list) != 0 {
			tmp[key] = list
		}
	}
	Playlists = tmp
}

// Get the playlists that the file is in.
func GetPlaylistsByPath(Path string) ([]string, error) {
	allPlaylists, err := LoadPlaylists()
	if err != nil {
		return nil, err
	}
	playlists := []string{}
	for name, list := range allPlaylists {
		for _, fp := range list {
			if fp == Path {
				playlists = append(playlists, name)
			}
		}
	}

	return playlists, nil
}

func SavePlaylists() error {
	SanitizePlaylists()
	bytes, err := json.MarshalIndent(Playlists, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(playlistsFile, bytes, 0644)
}
