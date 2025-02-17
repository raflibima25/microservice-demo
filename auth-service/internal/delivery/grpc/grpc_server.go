package grpc

import (
	"fmt"
	"log"
	"net"

	pb "grpc/pb/auth"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	address string
	server  *grpc.Server
}

func NewGRPCServer(address string) *GRPCServer {
	// create new server
	server := grpc.NewServer()

	return &GRPCServer{
		address: address,
		server:  server,
	}
}

func (s *GRPCServer) RegisterGRPCServices(authHandler pb.AuthServiceServer) {
	pb.RegisterAuthServiceServer(s.server, authHandler)
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Printf("gRPC server is listening on %s", s.address)

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}
