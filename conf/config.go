package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	// path to wallpaper data duh!
	WallpaperData string `json:"wallpaperData"`
	// Directory to cache wallpapers in. (default $TMP/wallpapers)
	CacheDir string `json:"overwriteCacheDir"`
}

// config :)
var Conf *Config

// Load/Reload the config - must be called once at the start
func Load(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &Conf)
	if err != nil {
		return err
	}
	if Conf.CacheDir == "" {
		Conf.CacheDir = filepath.Join(os.TempDir(), "wallpapers")
	}
	return nil
}
