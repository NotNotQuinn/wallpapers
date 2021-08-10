package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type confType struct {
	// path to wallpaper data duh!
	WallpaperData string `json:"wallpaperData"`
	// Directory to cache wallpapers in. (default $TMP/wallpapers)
	CacheDir string `json:"overwriteCacheDir"`
}

// config :)
var conf *confType
var confPath string = DefaultPath
var loadedChan = make(chan struct{})
var hasLoaded = false

// Default config path, used if path is never set.
const DefaultPath = "./wallpaperconf.json"

// Get the config, ensure it is loaded.
func Get() *confType {
	return conf
}

// Set the config path, and reload the config.
func SetPath(newPath string) error {
	confPath = newPath
	return load()
}

// Get the current file path for the config, not guaranteed to exists.
func GetPath() string {
	return confPath
}

// Force reload the config.
func Reload() error {
	return load()
}

// Returns a channel that will be closed when the config is first loaded.
func FirstLoad() <-chan struct{} {
	return loadedChan
}

// Load/Reload the config - must be called once before the config is used.
func load() error {
	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return err
	}
	if conf.CacheDir == "" {
		conf.CacheDir = filepath.Join(os.TempDir(), "wallpapers")
	}

	// Broadcast that it has loaded.
	if !hasLoaded {
		hasLoaded = true
		close(loadedChan)
	}

	return nil
}

func init() {
	// ignore error, but attempt to load in case there is a reliance on the default path.
	_ = load()
}
