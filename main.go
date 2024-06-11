package main

import (
	"fmt"
	"github.com/GRPCPractice/proto/proto/helloworld"
	"github.com/GRPCPractice/proto/proto/user"
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
	user.RegisterUserServiceServer(s, &userServer{
		users: make(map[string]User),
	})
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}
