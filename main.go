package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"ShopServer/controllers/v1"
	"ShopServer/constant"
	"ShopServer/postgresql"
	"ShopServer/proto/ShopServer"
)

func init() {
	constant.ReadConfig(".env")
	postgresql.Initialize()
}

func main() {
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	ShopServer.RegisterShopServerServer(s, &v1.ShopServe{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
