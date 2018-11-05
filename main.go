package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	version := os.Getenv("VERSION")
	fmt.Println("Expected version is: ", version)
	if actualVersion := requestVersion(); version != actualVersion {
		panic("Actual version is not equals to expected")
	}
	httpClient := &http.Client{Timeout: time.Second * 15}
	url := "http://vpered.su/?translate=Hello+World&group=Romanic"
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

func requestVersion() string {
	httpClient := &http.Client{Timeout: time.Second * 15}
	versionUrl := "http://vpered.su:8080/version"
	fmt.Println("VERSIONURL:>", versionUrl)

	req, err := http.NewRequest("GET", versionUrl, nil)

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
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	return bodyString
}
