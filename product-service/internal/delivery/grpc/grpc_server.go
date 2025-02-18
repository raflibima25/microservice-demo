package grpc

import (
	"fmt"
	pb "grpc/pb/product"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	address string
	server  *grpc.Server
}

func NewGRPCProductServer(address string) *Server {
	// create a new gRPC server
	server := grpc.NewServer()

	return &Server{
		address: address,
		server:  server,
	}
}

func (s *Server) RegisterServices(productHandler pb.ProductServiceServer) {
	pb.RegisterProductServiceServer(s.server, productHandler)
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Printf("gRPC server is running on %s", s.address)

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
}
