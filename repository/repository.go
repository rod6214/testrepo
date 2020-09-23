package main

import (
	"log"
	"net"

	"github.com/southworks/gnalog/demo/repository/items"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting Repository...")
	lis, err := net.Listen("tcp", ":9000")
	log.Println("Listening to port 9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %s", err)
	}
	grpcServer := grpc.NewServer()
	listServer := items.Server{}
	items.RegisterItemServiceServer(grpcServer, &listServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %s", err)
	}
}
