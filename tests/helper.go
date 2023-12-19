package tests

import (
	"io"
	"net/http"
)

func ResponseBodyToBytes(resp *http.Response) []byte {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return bodyBytes
}

func ResponseBodyToString(resp *http.Response) string {
	return string(ResponseBodyToBytes(resp))
}
