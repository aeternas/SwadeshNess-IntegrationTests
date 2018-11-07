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
	time.Sleep(5 * time.Second)
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

	requestEndpoint("http://vpered.su/?translate=Hello+World&group=Romanic")
	log.Printf("Translation OK")

	requestEndpoint("http://vpered.su/groups")
	log.Printf("Groups OK")
}

func requestVersion() string {
	httpClient := &http.Client{Timeout: time.Second * 15}
	versionUrl := "http://vpered.su/version"

	req, err := http.NewRequest("GET", versionUrl, nil)

	if err != nil {
		log.Printf("Error during initializing request")
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		log.Printf("Error during executing request")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return string(resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error during unmarshalling, trying to repeat")
		return "Error!"
	}

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
