package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/kelgendy1204/currency-converter/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	currencyconverter.CurrencyServer
}

func (s *server) Convert(ctx context.Context, in *emptypb.Empty) (*currencyconverter.ConvertValue, error) {
	currencyValue := scrapeCurrencyConverter()
	return &currencyconverter.ConvertValue{Value: currencyValue}, nil
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

func setupServer() {
	lis, err := net.Listen("tcp", ":9000")

	if err != nil {
		panic(err)
	}

	defer lis.Close()

	s := server{}

	grpcServer := grpc.NewServer()
	currencyconverter.RegisterCurrencyServer(grpcServer, &s)

	log.Printf("server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func main() {
	setupServer()
}
