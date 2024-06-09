package main

import (
	"context"
	"fmt"
	"github.com/GRPCPractice/proto/proto/helloworld"
	"io"
)

type hellowordServer struct {
	helloworld.UnimplementedGreeterServer
}

func (s *hellowordServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		fmt.Println("Hello Request: ", in.GetName())
		return &helloworld.HelloReply{Message: "Hello, World!"}, nil
	}
}

func (s *hellowordServer) StreamHelloRequests(stream helloworld.Greeter_StreamHelloRequestsServer) error {
	var names []string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("End of stream")
			message := "Hello"
			for _, name := range names {
				message += ", " + name
			}
			return stream.SendAndClose(&helloworld.HelloReply{Message: message})
		}

		if err != nil {
			return err
		}
		fmt.Println("Hello Request: ", in.GetName())
		names = append(names, in.GetName())
	}
}

func (s *hellowordServer) StreamHelloReplies(in *helloworld.HelloRequest, stream helloworld.Greeter_StreamHelloRepliesServer) error {
	for i := 0; i < 10; i++ {
		message := "Hello, " + in.GetName() + " " + fmt.Sprint(i)
		fmt.Println(message)
		if err := stream.Send(&helloworld.HelloReply{Message: message}); err != nil {
			return err
		}
	}

	return nil
}

func (s *hellowordServer) SayHelloChat(stream helloworld.Greeter_SayHelloChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Println("Hello Request: ", in.GetName())
		message := "Hello, " + in.GetName()
		if err := stream.Send(&helloworld.HelloReply{Message: message}); err != nil {
			return err
		}
	}
}
