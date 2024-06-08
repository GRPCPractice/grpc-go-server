package main

import (
	"context"
	"fmt"
	"github.com/GRPCPractice/proto/proto/helloworld"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Println("Hello Request: ", in.GetName())

	return &helloworld.HelloReply{Message: "Hello, World!"}, nil
}

func main() {
	fmt.Println("Hello, World!")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		panic(err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}
