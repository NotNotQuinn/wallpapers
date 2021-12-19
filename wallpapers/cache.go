package wallpapers

import (
	"errors"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Keep It Simple Stupid
const wallpaperDir = "d:\\wallpapers\\"

var cacheDir string

func init() {
	cacheDir = filepath.Join(wallpaperDir, "cached")
}

// type ImgCache struct {
// 	// Maps url to file path
// 	Images map[string]string `json:"images"`
// }

// Cache of URLs to paths to downloaded images

// downloads file and saves it, returning path if successful
func AddUrl(Url string) (string, error) {
	path, err := CalculatePath(Url)
	if err != nil {
		return "", err
	}

	err = DownloadFile(Url, path)
	if err != nil {
		if !errors.Is(err, fs.ErrExist) {
			return "", err
		}
	}
	return path, nil
}

// saves a file as specific url, but copies the file at `path` rather than downloading
func AddFile(Url, originalPath string) error {
	path, err := CalculatePath(Url)
	if err != nil {
		return err
	}
	bytes, err := os.ReadFile(originalPath)
	if err != nil {
		return err
	}
	// create parent directories
	err = os.MkdirAll(filepath.Dir(path), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// delete all cached images
func ClearCache() error {
	return os.Remove(cacheDir)
}

// Calculates the origin URL for the given cache path
func CalculateURL(Path string) (string, error) {
	Path = filepath.Clean(Path)
	if !strings.HasPrefix(Path, cacheDir) {
		return "", errors.New("path is not within cache directory, unable to determine origin url")
	}
	Path = strings.TrimPrefix(Path, cacheDir)
	// left with something like "\\example.com\\url%20stuff\\thing.jpg" (or using '/' on non-windows systems)
	Path = strings.Replace(Path, "\\", "/", -1)
	Path = strings.TrimPrefix(Path, "/")
	return "https://" + Path, nil
}

// Calculates the cache path for the given url
func CalculatePath(Url string) (string, error) {
	parsed, err := url.Parse(Url)
	if err != nil {
		return "", err
	}

	return filepath.Clean(filepath.Join(cacheDir, parsed.Hostname()+"/"+parsed.EscapedPath())), nil
}
