package main

import (
	"github.com/jinzhu/gorm"
	pb "github.com/mooncaker816/shipper/user-service/proto/user"
)

type Repository interface {
	Create(*pb.User) error
	Get(string) (*pb.User, error)
	GetAll() ([]*pb.User, error)
	GetByEmail(string) (*pb.User, error)
}

// Create(context.Context, *User, *Response) error
// Get(context.Context, *User, *Response) error
// GetAll(context.Context, *Request, *Response) error
// Auth(context.Context, *User, *Token) error
// ValidateToken(context.Context, *Token, *Token) error
type UserRepository struct {
	db *gorm.DB
}

func (repo *UserRepository) Create(user *pb.User) error {
	err := repo.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) Get(id string) (*pb.User, error) {
	var user pb.User
	user.Id = id
	err := repo.db.First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetAll() ([]*pb.User, error) {
	var users []*pb.User
	err := repo.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *UserRepository) GetByEmail(email string) (*pb.User, error) {
	var user pb.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
