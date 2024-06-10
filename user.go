package main

import (
	"context"
	"fmt"
	"github.com/GRPCPractice/proto/proto/user"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sync"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type userServer struct {
	user.UnimplementedUserServiceServer
	userSeq int
	users   map[string]User
	mu      sync.Mutex
}

func (s *userServer) GetUser(ctx context.Context, in *user.UserID) (*user.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if u, ok := s.users[in.GetId()]; ok {
		return &user.User{
			Id:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
		}, nil
	}

	return nil, fmt.Errorf("user not found")
}

func (s *userServer) CreateUser(ctx context.Context, in *user.CreateUserRequest) (*user.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	u := User{
		ID:        s.NewUserID(),
		Name:      in.GetName(),
		Email:     in.GetEmail(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	timestamp := timestamppb.New(now)
	return &user.User{
		Id:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}, nil
}

func (s *userServer) UpdateUser(ctx context.Context, in *user.UpdateUserRequest) (*user.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if u, ok := s.users[in.GetId()]; ok {
		u.Name = in.GetName()
		u.Email = in.GetEmail()
		u.UpdatedAt = time.Now()
		s.users[in.GetId()] = u

		return &user.User{
			Id:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
		}, nil
	}

	return nil, fmt.Errorf("user not found")
}

func (s *userServer) DeleteUser(ctx context.Context, in *user.UserID) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[in.GetId()]; ok {
		delete(s.users, in.GetId())
	}

	return &emptypb.Empty{}, nil
}

func (s *userServer) ListUsers(ctx context.Context, in *emptypb.Empty) (*user.UserList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var userList []*user.User
	for _, u := range s.users {
		userList = append(userList, &user.User{
			Id:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
		})
	}

	return &user.UserList{Users: userList}, nil
}

func (s *userServer) NewUserID() string {
	for {
		s.userSeq++
		if _, ok := s.users[string(s.userSeq)]; !ok {
			return string(s.userSeq)
		}
	}
}
