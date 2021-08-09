package wallpapers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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

// Solution to a cloudproxy request
type CloudProxySolution struct {
	// URL requested
	URL string `json:"url"`
	// HTTP status from request
	Status int `json:"status"`
	// Headers in response
	Headers map[string]string `json:"headers"`
	// Raw text response of page
	Response string `json:"response"`
	// Our user-agent
	Useragent string `json:"userAgent"`
	// Cookies returned
	Cookies []Cookie `json:"cookies"`
}

type CloudProxyResponse struct {
	Solution *CloudProxySolution `json:"solution"`
	// "ok" or "error"
	Status string `json:"status"`
	// Error message
	Message        string `json:"message"`
	Starttimestamp int64  `json:"startTimestamp"`
	Endtimestamp   int64  `json:"endTimestamp"`
	// CloudProxy version
	Version string `json:"version"`
}

// A browser cookie
type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	Size     int     `json:"size"`
	Httponly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
	Session  bool    `json:"session"`
	Samesite string  `json:"sameSite"`
}

// Cookies created in cloudproxy sessions
var CloudProxyCookies []Cookie

// Uses cloudproxy (if availible on local port 8191) to get the content of a website.
// Does not solve captchas, only bypasses CloudFlare IUAM.
//
// CloudProxy: https://github.com/NoahCardoza/CloudProxy
func CloudProxyGetContent(Url string) (content *CloudProxyResponse, err error) {
	req, err := http.NewRequest("POST", "http://localhost:8191/v1",
		strings.NewReader(fmt.Sprintf(`{
			"cmd": "request.get",
			"url":"%s",
			"userAgent": "%s",
			"maxTimeout": 20000
		  }`, Url, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36"),
		))
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &content)
	if err != nil {
		return nil, err
	}
	if content.Status != "ok" {
		return nil, fmt.Errorf("cloudproxy: %s", content.Message)
	}
	return content, nil
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
