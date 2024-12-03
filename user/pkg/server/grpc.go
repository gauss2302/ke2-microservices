package server

import (
	"fmt"
	"github.com/gauss2302/testcommm/user/config"
	"github.com/gauss2302/testcommm/user/internal/service"
	pb "github.com/gauss2302/testcommm/user/proto"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	grpcServer  *grpc.Server
	userService *service.UserService
}

func NewServer(userService *service.UserService) *Server {
	return &Server{
		userService: userService,
		grpcServer:  grpc.NewServer(),
	}
}

func (s *Server) Start(cfg config.GRPCConfig) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	pb.RegisterUserServiceServer(s.grpcServer, s.userService)
	return s.grpcServer.Serve(lis)

}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}
