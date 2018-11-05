package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	httpClient := &http.Client{Timeout: time.Second * 15}
	url := "http://vpered.su:8080/dev/?translate=Hello+World&group=Romanic"
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Invalid response status")
	}
}
