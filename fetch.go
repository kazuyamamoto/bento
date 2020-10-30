package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// fetch は location からデータを得る.
// location は URL やファイルパスである.
// 戻り値の io.ReadCloser.Close を呼ぶのは、fetch を呼んだ側の責務.
type fetch func(location string) (io.ReadCloser, error)

func fetchHTTPBody(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("issueing a GET request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}

	return res.Body, nil
}

// fetchFile は file をオープンする.
func fetchFile(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %v", path, err)
	}

	return f, nil
}
