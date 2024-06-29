package main

import (
	"context"
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
	s.mu.Lock()
	if _, ok := s.chatChanel[in.GetUserId()]; !ok {
		s.chatChanel[in.GetUserId()] = make(chan *chat.ChatMessage)
	}
	s.mu.Unlock()

	for {
		select {
		case msg := <-s.chatChanel[in.GetUserId()]:
			if err := stream.Send(msg); err != nil {
				return err
			}
		case <-stream.Context().Done():
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
