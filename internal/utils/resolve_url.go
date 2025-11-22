package utils

import (
	"net/http"
	"strings"
	"time"
)

func ResolveURL(url string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()
	// remove /sent/?... part => ex: https://www.pinterest.com/pin/xxxxx/sent/?... => https://www.pinterest.com/pin/xxxxx/
	finalURL = strings.Split(finalURL, "/sent/")[0] + "/"

	return finalURL, nil
}
