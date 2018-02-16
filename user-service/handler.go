package main

import (
	"context"

	pb "github.com/mooncaker816/shipper/user-service/proto/user"
)

type service struct {
	repo Repository
}

func (s *service) Create(ctx context.Context, user *pb.User, res *pb.Response) error {
	err := s.repo.Create(user)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (s *service) Get(ctx context.Context, user *pb.User, res *pb.Response) error {
	user, err := s.repo.Get(user.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (s *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

func (s *service) Auth(ctx context.Context, user *pb.User, tok *pb.Token) error {
	return nil
}

func (s *service) ValidateToken(ctx context.Context, tok1 *pb.Token, tok2 *pb.Token) error {
	return nil
}

// Create(context.Context, *User, *Response) error
// Get(context.Context, *User, *Response) error
// GetAll(context.Context, *Request, *Response) error
// Auth(context.Context, *User, *Token) error
// ValidateToken(context.Context, *Token, *Token) error
