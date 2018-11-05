package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var version string = fmt.Sprintf("%q", os.Getenv("VERSION"))
	log.Println("Expected version is:", version)
	i := 0
	var actualVersion string
	for {
		actualVersion = requestVersion()
		if i >= 10 || actualVersion == version {
			break
		}
		time.Sleep(2 * time.Second)
		log.Printf("Seems actual version doesn't match to expected, retrying...")
		i++
	}
	if version != actualVersion {
		log.Fatalf("Actual version doesn't match to expected. Actual is: %v", actualVersion)
	}

	requestEndpoint("http://vpered.su:8080/dev/?translate=Hello+World&group=Romanic")
	log.Printf("Translation OK")

	requestEndpoint("http://vpered.su:8080/dev/groups")
	log.Printf("Groups OK")
}

func requestVersion() string {
	httpClient := &http.Client{Timeout: time.Second * 15}
	versionUrl := "http://vpered.su:8080/dev/version"

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

func requestEndpoint(e string) {
	httpClient := &http.Client{Timeout: time.Second * 15}
	url := e
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
