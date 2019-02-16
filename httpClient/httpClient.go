package httpClient

import (
	"net/http"
	"time"
)

func NewHttpClient() *http.Client {
	return &http.Client{Timeout: time.Second * 30}
}
