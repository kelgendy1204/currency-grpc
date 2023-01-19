package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"github.com/kelgendy1204/currency-converter/service"
)

func setupServer() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	defer lis.Close()
	s := service.Server{}
	grpcServer := grpc.NewServer()
	service.RegisterCurrencyServer(grpcServer, &s)
	log.Printf("server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func main() {
	setupServer()
}
