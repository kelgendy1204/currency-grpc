package service

import (
	context "context"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	CurrencyServer
}

func scrapeCurrencyConverter() string {
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

	return currencyValue
}

func (s *Server) Convert(ctx context.Context, in *emptypb.Empty) (*ConvertValue, error) {
	currencyValue := scrapeCurrencyConverter()
	return &ConvertValue{Value: currencyValue}, nil
}
