package requests

type GetRequest interface {
	Execute(s string) ([]byte, error)
}

type getRequest struct {
	Endpoint string
}

func NewGetRequest(s string) GetRequest {
	return getRequest{Endpoint: s}
}

func (r *getRequest) Execute(s string) ([]byte, eror) {
	url := s
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	httpClient := getClient()

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
