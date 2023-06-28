package translator

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Jeffail/gabs/v2"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const tranlsateUrl string = "https://translate.googleapis.com/translate_a/single"

/**
 * RequestTranslation function
 * @param {RequestBody} body
 * @param {chan string} ch
 * @param {sync.WaitGroup} wg
 */
func RequestTranslation(body *RequestBody, ch chan string, wg *sync.WaitGroup) {
	// Create the HTTP client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// Create the request
	req, err := http.NewRequest("GET", tranlsateUrl, nil)
	if err != nil {
		log.Fatalf("There was an error: %v", err)
	}
	// Add the query parameters
	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)
	req.URL.RawQuery = query.Encode()
	// Send the request
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Oops! There was an error with Google API: %v", err)
	}
	// Close the response body
	defer res.Body.Close()
	// Check the response status code
	if res.StatusCode == http.StatusTooManyRequests {
		ch <- "Oops! You have exceeded the Google API quota."
		wg.Done()
		return
	}
	if res.StatusCode != http.StatusOK {
		ch <- "Oops! There was an error with Google API."
		wg.Done()
		return
	}
	// Decode the response
	parsedJson, err := gabs.ParseJSONBuffer(res.Body)
	if err != nil {
		log.Fatalf("Oops! There was an error parsing the response: %v", err)
	}
	// Get the translation
	nestOne, err := parsedJson.ArrayElement(0)
	if err != nil {
		log.Fatalf("Oops! There was a problem with the response el(0): %v", err)
	}
	nestTwo, err := nestOne.ArrayElement(0)
	if err != nil {
		log.Fatalf("Oops! There was a problem with the response el(1): %v", err)
	}
	nestThree, err := nestTwo.ArrayElement(0)
	if err != nil {
		log.Fatalf("Oops! There was a problem with the response el(2): %v", err)
	}
	// Send the translation to the channel
	ch <- nestThree.Data().(string)
	// Remove 1 from the WaitGroup
	wg.Done()
}
