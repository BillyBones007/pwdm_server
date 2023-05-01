package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/storage"
	pb "github.com/BillyBones007/pwdm_service_api/api"
)

// Authentication service.
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	Rep *storage.Storage
}

// NewAuthService - constructor AuthService.
func NewAuthService(r *storage.Storage) *AuthService {
	return &AuthService{Rep: r}
}

// Create - create new user.
func (a *AuthService) Create(ctx context.Context, in *pb.AuthReq) (*pb.AuthResp, error) {
	// TODO: realization
	return nil, nil
}

// Enter - authorization the user.
func (a *AuthService) Enter(ctx context.Context, in *pb.AuthReq) (*pb.AuthResp, error) {
	// TODO: realization
	return nil, nil
}
