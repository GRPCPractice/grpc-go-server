package main

import (
	"fmt"
	"github.com/GRPCPractice/proto/proto/helloworld"
	"google.golang.org/grpc"
	"net"
)

func main() {
	fmt.Println("Hello, World!")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		panic(err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &hellowordServer{})
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}
