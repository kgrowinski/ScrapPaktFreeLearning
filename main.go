package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jasonlvhit/gocron"
)

func getImageURL(image string) (URL string) {
	imageURLStart := strings.Index(image, "src") + 5
	imageURLEnds := strings.Index(image, "alt") - 2
	URL = image[imageURLStart:imageURLEnds]
	return
}

func scrape() {
	// Request the HTML page.
	URL := "https://www.packtpub.com/packt/offers/free-learning"
	temsWebHook := "https://outlook.office.com/webhook/1e016c31-31d6-431a-9596-4ce52e73c184@d8b0675b-bea4-4002-9289-1da6ecc1ab6c/IncomingWebhook/8580c5d7a5894b0e8dc55f133a9d86aa/4cde4289-a034-4f8b-bee9-5dcf62170dec"
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	image := doc.Find("#deal-of-the-day").Find("noscript").Text()
	title := doc.Find("#title-bar-title").Text()
	imageURL := getImageURL(image)
	message := fmt.Sprintf(`{
		"@type": "MessageCard",
		"@context": "http://schema.org/extensions",
		"themeColor": "0076D7",
		"summary": "Gopher Recomends:",
		"sections": [{
			"activityTitle": "Gopher Recomends for todays read:",
			"activitySubtitle": "%s",
			"activityImage": "%s",
			"facts": [{
				"name": "Check it out on",
				"value": "[Pakt Free Learning](%s)"
			}],
			"markdown": true
		}]
	}`, title, imageURL, URL)
	fmt.Println(message)

	req, err := http.NewRequest("POST", temsWebHook, bytes.NewBuffer([]byte(message)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func main() {
	s := gocron.NewScheduler()
	s.Every(24).Hours().Do(scrape)
	<-s.Start()
}
