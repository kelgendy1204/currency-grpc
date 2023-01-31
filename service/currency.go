package service

import (
	context "context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Server struct {
	CurrencyServer
}

func scrapeCurrencyConverter(convertInput *ConvertInput) string {
	fromCurrency := convertInput.From
	toCurrency := convertInput.To

	fmt.Println("Fetch from " + fromCurrency + " to " + toCurrency)

	res, err := http.Get("https://www.currency.me.uk/convert/" + fromCurrency + "/" + toCurrency)

	// Request the HTML page.
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Print(err)
	}

	currencyValue, _ := doc.Find("input#answer").Attr("value")

	return currencyValue
}

type currencyData struct {
	time time.Time
	value string
}

func GetCurrencyFetcher() func(*ConvertInput) string {
	cache := make(map[string]currencyData)

	return func(convertInput *ConvertInput) string {
		currencyKey := convertInput.From + "-" + convertInput.To
		val, found := cache[currencyKey]
		currentTime := time.Now()

		if found {
			durationHours := currentTime.Sub(val.time).Hours()

			if durationHours > 6 {
				val.time = currentTime
				val.value = scrapeCurrencyConverter(convertInput)
			}

			return val.value
		}

		currencyValue := scrapeCurrencyConverter(convertInput)

		data := currencyData { time: currentTime, value: currencyValue }

		cache[currencyKey] = data

		return currencyValue
	}
}

var fetchCurrency = GetCurrencyFetcher()

func (s *Server) Convert(ctx context.Context, convertInput *ConvertInput) (*ConvertValue, error) {
	currencyValue := fetchCurrency(convertInput)
	return &ConvertValue{Value: currencyValue}, nil
}
