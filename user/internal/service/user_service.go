package service

import (
	"context"
	"errors"

	"github.com/gauss2302/testcommm/user/internal/domain/entity"
	"github.com/gauss2302/testcommm/user/internal/repository"
	pb "github.com/gauss2302/testcommm/user/proto"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &pb.User{
		Id:    uint64(user.ID),
		Email: user.Email,
	}, nil
}

func (s *UserService) VerifyUser(ctx context.Context, req *pb.VerifyUserRequest) (*pb.User, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &pb.User{
		Id:    uint64(user.ID),
		Email: user.Email,
	}, nil
}
