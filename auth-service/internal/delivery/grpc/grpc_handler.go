package grpc

import (
	"auth-service/internal/domain"
	"context"
	pb "grpc/pb/auth"
)

type GRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	authUseCase domain.AuthUseCase
}

func NewGRPCHandler(authUseCase domain.AuthUseCase) *GRPCHandler {
	return &GRPCHandler{
		authUseCase: authUseCase,
	}
}

func (h *GRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	user, token, err := h.authUseCase.Register(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		Token: token,
		User: &pb.UserData{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (h *GRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	user, token, err := h.authUseCase.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.AuthResponse{
		Token: token,
		User: &pb.UserData{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (h *GRPCHandler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	user, err := h.authUseCase.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateResponse{
			Valid: false,
			User:  nil,
		}, nil
	}

	return &pb.ValidateResponse{
		Valid: true,
		User: &pb.UserData{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}

func (h *GRPCHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := h.authUseCase.Logout(req.Token)
	if err != nil {
		return &pb.LogoutResponse{
			Success: false,
		}, err
	}

	return &pb.LogoutResponse{
		Success: true,
	}, nil
}

// serve start gRPC server
func (h *GRPCHandler) Serve(address string) error {
	server := NewGRPCServer(address)
	server.RegisterGRPCServices(h)
	return server.Start()
}
