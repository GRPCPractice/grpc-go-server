package main

import (
	"context"
	"fmt"
	"github.com/GRPCPractice/proto/proto/chat"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type chatServer struct {
	chat.UnimplementedChatServiceServer
	chatChanel map[string]chan *chat.ChatMessage
	mu         sync.Mutex
}

func (s *chatServer) Connect(in *chat.ConnectRequest, stream chat.ChatService_ConnectServer) error {
	fmt.Println("Connect: ", in.GetUserId())
	s.mu.Lock()
	if _, ok := s.chatChanel[in.GetUserId()]; !ok {
		s.chatChanel[in.GetUserId()] = make(chan *chat.ChatMessage)
	}
	s.mu.Unlock()

	for {
		select {
		case msg := <-s.chatChanel[in.GetUserId()]:
			if err := stream.Send(msg); err != nil {
				fmt.Println("Send Fail: ", in.GetUserId())
				return err
			}
		case <-stream.Context().Done():
			fmt.Println("Disconnect: ", in.GetUserId())
			s.mu.Lock()
			delete(s.chatChanel, in.GetUserId())
			s.mu.Unlock()
			return nil
		}
	}

	return nil
}

func (s *chatServer) Send(ctx context.Context, msg *chat.ChatMessage) (*emptypb.Empty, error) {
	s.mu.Lock()
	for _, ch := range s.chatChanel {
		ch <- msg
	}
	s.mu.Unlock()

	return &emptypb.Empty{}, nil
}
