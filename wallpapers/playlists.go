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

func LoadPlaylists() (*mappedPlaylists, error) {
	if playlistsLoaded {
		return &Playlists, nil
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
	return &Playlists, nil
}

func SavePlaylists() error {
	bytes, err := json.MarshalIndent(Playlists, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(playlistsFile, bytes, 0)
}
