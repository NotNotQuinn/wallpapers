package wallpapers

import (
	"fmt"
	"io/ioutil"
	"net/http"
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

// Wraps http.NewRequest (with GET and no body) passed to http.DefaultClient.Do with user agent set to a browser's
func httpGetWithFakeUserAgent(Url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil, err
	}
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36")
	return http.DefaultClient.Do(req)
}

func getEyyIndexerFiles(Url string) (urls []string, err error) {
	findFiles := regexp.MustCompile(`<td data-raw="image" class="download"><a href="([^"]*)"`)
	resp, err := httpGetWithFakeUserAgent(Url)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := string(bytes)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println(body)
		panic(fmt.Errorf("status code is %v; body is in stdout", resp.StatusCode))
	}
	matches := findFiles.FindAllStringSubmatch(body, -1)
	baseURL := "https://" + resp.Request.URL.Host
	for _, match := range matches {
		urls = append(urls, baseURL+match[1])
	}
	return urls, nil
}
