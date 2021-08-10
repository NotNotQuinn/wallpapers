package wallpapers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
	Cookies []CloudProxyCookie `json:"cookies"`
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
type CloudProxyCookie struct {
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
var CloudProxyCookies []CloudProxyCookie

// Adds cookies to the cookie cache, but only ones that arent already there.
func addCookies(cookies []CloudProxyCookie) {
	for _, cookie := range CloudProxyCookies {
		for i, newCookie := range cookies {
			if cookie == newCookie {
				// remove matching cookie
				cookies[i] = cookies[len(cookies)-1]
				cookies = cookies[:len(cookies)-1]
			}
		}
	}
	CloudProxyCookies = append(CloudProxyCookies, cookies...)
}

// Uses cloudproxy (if availible on local port 8191) to get the content of a website.
// Does not solve captchas, only bypasses CloudFlare IUAM.
//
// CloudProxy: https://github.com/NoahCardoza/CloudProxy
func CloudProxyGetContent(Url string) (response string, err error) {
	requestJSON, err := json.Marshal(struct {
		Cmd        string             `json:"cmd"`
		Url        string             `json:"url"`
		UserAgent  string             `json:"userAgent"`
		MaxTimeout int                `json:"maxTimeout"`
		Cookies    []CloudProxyCookie `json:"cookies"`
	}{
		"request.get",
		Url,
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36",
		20000,
		CloudProxyCookies,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "http://localhost:8191/v1", bytes.NewReader(requestJSON))
	if err != nil {
		return "", err
	}

	var cloudproxy bool = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// Likely that cloudproxy is not running, and we couldnt connect.
		// Try a direct GET
		resp, err = http.Get(Url)
		if err != nil {
			return "", err
		}
		cloudproxy = false
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if !cloudproxy {
		return string(bytes), nil
	}

	var content *CloudProxyResponse
	err = json.Unmarshal(bytes, &content)
	if err != nil {
		return "", err
	}

	if content.Status != "ok" {
		return "", fmt.Errorf("cloudproxy: %s", content.Message)
	}

	addCookies(content.Solution.Cookies)
	return content.Solution.Response, nil
}
