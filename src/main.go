package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/subosito/gotenv"
)

type Product struct {
	Data []struct {
		ID              string      `json:"id"`
		ProductID       string      `json:"productId"`
		AvailableFrom   time.Time   `json:"availableFrom"`
		ExpiresAt       time.Time   `json:"expiresAt"`
		LimitedAmount   bool        `json:"limitedAmount"`
		AmountAvailable interface{} `json:"amountAvailable"`
		Details         interface{} `json:"details"`
		Priority        int         `json:"priority"`
		CreatedAt       time.Time   `json:"createdAt"`
		UpdatedAt       time.Time   `json:"updatedAt"`
		DeletedAt       interface{} `json:"deletedAt"`
	} `json:"data"`
	Count int `json:"count"`
}

type ProductSummary struct {
	Title           string    `json:"title"`
	Type            string    `json:"type"`
	CoverImage      string    `json:"coverImage"`
	ProductID       string    `json:"productId"`
	Isbn13          string    `json:"isbn13"`
	OneLiner        string    `json:"oneLiner"`
	Pages           int       `json:"pages"`
	PublicationDate time.Time `json:"publicationDate"`
	Length          string    `json:"length"`
	About           string    `json:"about"`
	Learn           string    `json:"learn"`
	Features        string    `json:"features"`
	Authors         []string  `json:"authors"`
	ShopURL         string    `json:"shopUrl"`
	ReadURL         string    `json:"readUrl"`
	Category        string    `json:"category"`
	EarlyAccess     bool      `json:"earlyAccess"`
	Available       bool      `json:"available"`
}

func init() {
	gotenv.Load()
}

func scrape() {

	// GENERATE URL
	productID := getProductID()

	// GET Product Details
	productDetial := getProductDetails(productID)

	// Prepare Message

	message := fmt.Sprintf(`{
		"@type": "MessageCard",
		"@context": "https://schema.org/extensions",
		"summary": "Issue 176715375",
		"themeColor": "0078D7",
		"sections": [
			{
				"activityTitle": "__%s__",
				"activityImage": "%s",
				"activitySubtitle": "Published %s"
			},{
				"text": "%s"
			},{
				"text" :"[Pakt Free Learning](%s)"
			}
		],"markdown": true
	}`,
		productDetial.Title,
		productDetial.CoverImage,
		parsePublicationDate(productDetial.PublicationDate),
		productDetial.OneLiner,
		"https://www.packtpub.com/packt/offers/free-learning",
	)

	//	Send to TEAMS
	temsWebHook := os.Getenv("WEBHOOK")
	req, err := http.NewRequest("POST", temsWebHook, bytes.NewBuffer([]byte(message)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print response Status
	fmt.Println("response Status:", resp.Status)
}

func parsePublicationDate(date time.Time) string {
	year, month, day := date.Date()
	newDate := fmt.Sprintf("%d %s %d", day, month, year)
	return newDate
}

func getProductDetails(productID string) ProductSummary {
	SummaryURL := fmt.Sprintf("https://static.packt-cdn.com/products/%s/summary", productID)
	summaryRes, err := http.Get(SummaryURL)
	if err != nil {
		log.Fatal(err)
	}
	if summaryRes.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", summaryRes.StatusCode, summaryRes.Status)
	}
	productDetails := &ProductSummary{}
	if err := json.NewDecoder(summaryRes.Body).Decode(&productDetails); err != nil {
		log.Panic(err.Error())
	}
	defer summaryRes.Body.Close()
	return *productDetails
}

func getProductID() string {
	today := time.Now()
	year, month, day := today.Date()
	dateFrom := fmt.Sprintf("%d-%d-%dT00:00:00.000Z", year, month, day)
	dateTo := fmt.Sprintf("%d-%d-%dT00:00:00.000Z", year, month, day+1)
	ProductURL := fmt.Sprintf("https://services.packtpub.com/free-learning-v1/offers?dateFrom=%s&dateTo=%s", dateFrom, dateTo)
	fmt.Println(ProductURL)

	// GET Product
	prosuctRes, err := http.Get(ProductURL)
	if err != nil {
		log.Fatal(err)
	}
	if prosuctRes.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", prosuctRes.StatusCode, prosuctRes.Status)
	}
	products := &Product{}
	if err := json.NewDecoder(prosuctRes.Body).Decode(&products); err != nil {
		log.Panic(err.Error())
	}
	defer prosuctRes.Body.Close()
	return products.Data[0].ProductID
}

func main() {
	scrape()
}
