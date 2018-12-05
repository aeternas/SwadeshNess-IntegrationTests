package requests

import (
	httpClient "github.com/aeternas/SwadeshNess-IntegrationTests/httpClient"
	"io/ioutil"
	"net/http"
)

type GetRequest interface {
	Execute(s string) (int, []byte)
}

type getRequest struct {
	Endpoint string
}

func NewGetRequest(s string) GetRequest {
	return &getRequest{Endpoint: s}
}

func (r *getRequest) Execute(s string) (int, []byte) {
	urlString := s
	req, err := http.NewRequest("GET", urlString, nil)

	if err != nil {
		panic(err)
	}

	httpClient := httpClient.NewHttpClient()

	resp, err := httpClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return resp.StatusCode, bodyBytes
}
