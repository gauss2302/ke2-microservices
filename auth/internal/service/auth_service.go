package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gauss2302/testcommm/auth/internal/pkg/jwt"
	pb_auth "github.com/gauss2302/testcommm/auth/proto/auth"
	pb_user "github.com/gauss2302/testcommm/auth/proto/user"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type AuthService struct {
	pb_auth.UnimplementedAuthServiceServer
	redis      *redis.Client
	userClient pb_user.UserServiceClient
	jwtMaker   *jwt.JWTMaker
}

func NewAuthService(redis *redis.Client, userClient pb_user.UserServiceClient, jwtMaker *jwt.JWTMaker) *AuthService {
	return &AuthService{
		redis:      redis,
		userClient: userClient,
		jwtMaker:   jwtMaker,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*jwt.TokenPair, error) {
	user, err := s.userClient.CreateUser(ctx, &pb_user.CreateUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	tokens, err := s.jwtMaker.CreateTokenPair(user.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tokens")
	}

	err = s.redis.Set(ctx,
		fmt.Sprintf("refresh_token:%d", user.Id),
		tokens.RefreshToken,
		time.Hour*24*7,
	).Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to save refresh token")
	}

	return tokens, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*jwt.TokenPair, error) {
	user, err := s.userClient.VerifyUser(ctx, &pb_user.VerifyUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, errors.Wrap(err, "invalid credentials")
	}

	tokens, err := s.jwtMaker.CreateTokenPair(user.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tokens")
	}

	err = s.redis.Set(ctx,
		fmt.Sprintf("refresh_token:%d", user.Id),
		tokens.RefreshToken,
		time.Hour*24*7,
	).Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to save refresh token")
	}

	return tokens, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req *pb_auth.VerifyTokenRequest) (*pb_auth.VerifyTokenResponse, error) {
	log.Printf("Received verification request for token: %s", req.Token)

	userID, err := s.jwtMaker.VerifyToken(req.Token)
	if err != nil {
		log.Printf("Failed to verify token: %v", err)
		return nil, err
	}

	log.Printf("Successfully verified token for user ID: %d", userID)
	return &pb_auth.VerifyTokenResponse{
		UserId: userID,
	}, nil
}

func (s *AuthService) GetSession(ctx context.Context, token string) (*pb_user.User, error) {
	// Верифицируем токен
	userID, err := s.jwtMaker.VerifyToken(token)
	if err != nil {
		return nil, errors.Wrap(err, "invalid token")
	}

	// Получаем информацию о пользователе через user service
	user, err := s.userClient.GetUserByID(ctx, &pb_user.GetUserByIDRequest{
		Id: userID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	return user, nil
}
