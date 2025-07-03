package fetch

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http2.Transport{
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: false, // Enable gzip compression
	},
}

func Fetch[T any](url string) (*T, error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	pkg := new(T)
	err = json.NewDecoder(res.Body).Decode(pkg)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}
