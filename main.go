package main

import (
	"encoding/json"
	"fmt"
	httpClient "github.com/aeternas/SwadeshNess-IntegrationTests/httpClient"
	dto "github.com/aeternas/SwadeshNess-packages/dto"
	languages "github.com/aeternas/SwadeshNess-packages/language"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	host   string
	client *http.Client
)

const (
	DEV_HOST      = "DEV_HOST"
	PROD_HOST     = "PROD_HOST"
	BRANCH        = "BRANCH"
	MASTER_BRANCH = "master"
	VERSION       = "VERSION"
)

func init() {
	if branch := os.Getenv(BRANCH); branch == MASTER_BRANCH {
		host = os.Getenv(PROD_HOST)
	} else {
		host = os.Getenv(DEV_HOST)
	}

	client = httpClient.NewHttpClient()
}

func main() {
	var version string = fmt.Sprintf("%q", os.Getenv(VERSION))
	log.Println("Expected version is:", version)
	log.Print("Sleeping to ensure app is deployed...")
	time.Sleep(5 * time.Second)
	log.Print("Woke!")
	i := 0
	var actualVersion string
	for {
		actualVersion = requestVersionV1()
		if i >= 10 || actualVersion == version {
			log.Printf("Version matches, OK")
			break
		}
		time.Sleep(10 * time.Second)
		log.Printf("Seems actual version doesn't match to expected, retrying...")
		i++
	}

	if version != actualVersion {
		log.Fatalf("Actual version doesn't match to expected. Actual is: %v", actualVersion)
	}

	requestTranslationDeterminedV1()
	log.Printf("Determined Translation OK")

	requestTranslationRandomizedV1()
	log.Printf("Randomized Translation OK")

	requestGroupsV1()
	log.Printf("Groups OK")
}

func requestVersionV1() string {
	versionUrl := fmt.Sprintf("%v/v1/version", host)

	log.Println("versionUrl is", versionUrl)

	req, err := http.NewRequest("GET", versionUrl, nil)

	if err != nil {
		log.Fatalln("Error during initializing request", err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln("Error during executing request", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return string(resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("Error during unmarshalling, trying to repeat")
	}

	bodyString := string(bodyBytes)

	return bodyString
}

func requestGroupsV1() {
	endpoint := fmt.Sprintf("%v/v1/groups", host)
	code, body := requestEndpointV1(endpoint)
	if code != http.StatusOK {
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

func requestTranslationDeterminedV1() {
	endpoint := fmt.Sprintf("%v/v1/?translate=Hello,+World&group=turkic", host)
	code, body := requestEndpointV1(endpoint)
	if code != http.StatusOK {
		log.Fatalf("Translation response code is not 200")
	}

	var data dto.SwadeshTranslation

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Error unmarshalling body")
	}

	translationResult := data.Results[0].Results[4].Translation

	var containsTranslatedWords = strings.Contains(translationResult, "Merhaba") && strings.Contains(translationResult, "Dünya")

	if !containsTranslatedWords {
		log.Fatalf("Result translation doesn't match expected one: %v", translationResult)
	}
}

func requestTranslationRandomizedV1() {
	r := rand.New(rand.NewSource(99999))
	num := r.Int31()
	word := fmt.Sprintf("Hello,+World+%v", num)
	endpoint := fmt.Sprintf("%v/v1/?translate=%v&group=turkic", host, word)
	code, body := requestEndpointV1(endpoint)
	if code != http.StatusOK {
		log.Fatalf("Translation response code is not 200")
	}

	var data dto.SwadeshTranslation

	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Error unmarshalling body")
	}

	translationResult := data.Results[0].Results[4].Translation

	if translationResult != fmt.Sprintf("Merhaba, Dünya %v", num) {
		log.Fatalf("Result translation doesn't match expected one: %v", translationResult)
	}
}

func requestEndpointV1(e string) (int, []byte) {
	req, err := http.NewRequest(http.MethodGet, e, nil)

	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)

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
