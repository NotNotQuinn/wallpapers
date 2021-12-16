package wallpapers

import (
	"net/url"
	"os"
	"path/filepath"
)

// Keep It Simple Stupid
const cacheDir = "d:\\wallpapers\\cached"

// type ImgCache struct {
// 	// Maps url to file path
// 	Images map[string]string `json:"images"`
// }

// Cache of URLs to paths to downloaded images

// downloads file and saves it
func AddUrl(Url string) error {
	path, err := CalculatePath(Url)
	if err != nil {
		return err
	}

	err = DownloadFile(Url, path)
	return err
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
	err = os.MkdirAll(filepath.Dir(path), 0)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0)
	if err != nil {
		return err
	}

	return nil
}

// delete all cached images
func ClearCache() error {
	return os.Remove(cacheDir)
}

func CalculatePath(Url string) (string, error) {
	parsed, err := url.Parse(Url)
	if err != nil {
		return "", err
	}

	return filepath.Clean(filepath.Join(cacheDir, parsed.Hostname()+"/"+parsed.EscapedPath())), nil
}
