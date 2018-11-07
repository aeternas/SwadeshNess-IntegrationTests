package main

import (
	"encoding/json"
	"fmt"
	dto "github.com/aeternas/SwadeshNess-packages/dto"
	languages "github.com/aeternas/SwadeshNess-packages/language"
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
		time.Sleep(10 * time.Second)
		log.Printf("Seems actual version doesn't match to expected, retrying...")
		i++
	}
	if version != actualVersion {
		log.Fatalf("Actual version doesn't match to expected. Actual is: %v", actualVersion)
	}

	requestTranslation()
	log.Printf("Translation OK")

	requestGroups()
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

func requestGroups() {
	code, body := requestEndpoint("http://vpered.su/groups")
	if code != 200 {
		log.Fatalf("Groups response code is not 200")
	}

	var data []languages.LanguageGroup

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Error unmarshalling body")
	}

	if data[0].Name != "Turkic" {
		log.Fatalf("Group is not Turkic")
	}
}

func requestTranslation() {
	code, body := requestEndpoint("http://vpered.su/?translate=me&group=turkic")
	if code != 200 {
		log.Fatalf("Translation response code is not 200")
	}

	var data dto.SwadeshTranslation

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Error unmarshalling body")
	}

	if data.Results[0].Results[3].Translation != "bana" {
		log.Fatalf("Result translation doesn't match expected one")
	}
}

func requestEndpoint(e string) (int, []byte) {
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

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return resp.StatusCode, bodyBytes
}
