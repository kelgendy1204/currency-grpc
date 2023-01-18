package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

func scrapeCurrencyConverter() {
	res, err := http.Get("https://www.currency.me.uk/convert/usd/egp")

	// Request the HTML page.
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

	currencyValue, _ := doc.Find("input#answer").Attr("value")
	fmt.Println(currencyValue);
}

func main() {
	scrapeCurrencyConverter()
}
