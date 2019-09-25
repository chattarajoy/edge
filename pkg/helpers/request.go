package helpers

import (
	"fmt"
	"net/http"
)

func GetHeaders(req *http.Request) []string {
	var headers []string
	for name, value := range req.Header {
		headers = append(headers, fmt.Sprintf("%v: %v", name, value))
	}
	return headers
}
