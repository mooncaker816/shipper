package main

import (
	"encoding/json"
	"errors"
	"log"

	"golang.org/x/net/context"

	"golang.org/x/crypto/bcrypt"

	"github.com/micro/go-micro/broker"
	pb "github.com/mooncaker816/shipper/user-service/proto/user"
)

const topic = "user.created"

type service struct {
	repo   Repository
	auth   Authable
	PubSub broker.Broker
}

func (s *service) Create(ctx context.Context, user *pb.User, res *pb.Response) error {
	hashPSW, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPSW)
	err = s.repo.Create(user)
	if err != nil {
		return err
	}
	res.User = user
	if err := s.publishEvent(user); err != nil {
		return err
	}
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

func (s *service) Auth(ctx context.Context, req *pb.User, tok *pb.Token) error {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return err
	}
	token, err := s.auth.Encode(user)
	if err != nil {
		return err
	}
	tok.Token = token
	return nil
}

func (s *service) ValidateToken(ctx context.Context, tok1 *pb.Token, tok2 *pb.Token) error {
	// Decode token
	log.Printf("decode tokenstr:%s", tok1.Token)
	claims, err := s.auth.Decode(tok1.Token)
	if err != nil {
		return err
	}

	log.Println(claims)

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	tok2.Valid = true

	return nil
}

// Create(context.Context, *User, *Response) error
// Get(context.Context, *User, *Response) error
// GetAll(context.Context, *Request, *Response) error
// Auth(context.Context, *User, *Token) error
// ValidateToken(context.Context, *Token, *Token) error

func (s *service) publishEvent(user *pb.User) error {
	// Marshal to JSON string
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Create a broker message
	msg := &broker.Message{
		Header: map[string]string{
			"id": user.Id,
		},
		Body: body,
	}

	// Publish message to broker
	if err := s.PubSub.Publish(topic, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	}

	return nil
}
