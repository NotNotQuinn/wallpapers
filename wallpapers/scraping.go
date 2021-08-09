package wallpapers

import (
	"fmt"
	"net/url"
	"regexp"
)

var ErrUnrecognisedWallpaperRepoType = fmt.Errorf("unrecognised wallpaper repo type")

// Will get all download links to files located in the directory
// Ignores subdirectories.
func GetFiles(repoType WallpaperRepoType, Url string) ([]string, error) {
	switch repoType {
	case Eyy_Indexer:
		return getEyyIndexerFiles(Url)
	default:
		return nil, ErrUnrecognisedWallpaperRepoType
	}
}

func getEyyIndexerFiles(Url string) (urls []string, err error) {
	findFiles := regexp.MustCompile(`<td data-raw="image" class="download"><a href="([^"]*)"`)

	content, err := CloudProxyGetContent(Url)
	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(content.Solution.URL)
	if err != nil {
		return nil, err
	}

	matches := findFiles.FindAllStringSubmatch(content.Solution.Response, -1)
	baseURL := "https://" + parsed.Host
	for _, match := range matches {
		urls = append(urls, baseURL+match[1])
	}
	return urls, nil
}
