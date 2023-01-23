package service

import (
	context "context"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Server struct {
	CurrencyServer
}

func scrapeCurrencyConverter(convertInput *ConvertInput) string {
	fromCurrency := convertInput.From
	toCurrency := convertInput.To

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

func (s *Server) Convert(ctx context.Context, convertInput *ConvertInput) (*ConvertValue, error) {
	currencyValue := scrapeCurrencyConverter(convertInput)
	return &ConvertValue{Value: currencyValue}, nil
}
