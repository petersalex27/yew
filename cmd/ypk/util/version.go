package util

import (
	"io"
	"net/http"
)

type Version struct {
	string
	latest bool
}

func GetLatest(url string) (v Version, err error) {
	var resp *http.Response

	resp, err = http.Get(url)
	if err != nil {
		return v, err
	}

	defer resp.Body.Close()

	// read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return v, err
	}


}

// (v([\d+]\.)*[\d+])|(@latest)
