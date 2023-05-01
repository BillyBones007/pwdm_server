package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/storage"
	pb "github.com/BillyBones007/pwdm_service_api/api"
)

// ShowInfoService - service contains methods for getting information about uploaded user data.
type ShowInfoService struct {
	pb.UnimplementedShowInfoServiceServer
	Rep *storage.Storage
}

// NewShowInfoService - constructor ShowInfoService.
func NewShowInfoService(r *storage.Storage) *ShowInfoService {
	return &ShowInfoService{Rep: r}
}

// GetInfo - get information for current user.
func (s *ShowInfoService) GetInfo(ctx context.Context, in *pb.Empty) (*pb.ShowItemsResp, error) {
	// TODO: realization
	return nil, nil
}
