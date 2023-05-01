package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/storage"
	pb "github.com/BillyBones007/pwdm_service_api/api"
)

// DeleteService - service contains methods for deleting user data.
type DeleteService struct {
	pb.UnimplementedDeleteServiceServer
	Rep *storage.Storage
}

// NewDeleteService - constructor DeleteService.
func NewDeleteService(r *storage.Storage) *DeleteService {
	return &DeleteService{Rep: r}
}

// DelItem - delete one item.
func (d *DeleteService) DelItem(ctx context.Context, in *pb.DeleteItemReq) (*pb.DeleteResp, error) {
	// TODO: realization
	return nil, nil
}

// DelAll - delete all data for current user.
func (d *DeleteService) DelAll(ctx context.Context, in *pb.DeleteAllReq) (*pb.DeleteResp, error) {
	// TODO: realization
	return nil, nil
}
