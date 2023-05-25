package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/storage/models"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	pb "github.com/BillyBones007/pwdm_service_api/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteService - service contains methods for deleting user data.
type DeleteService struct {
	pb.UnimplementedDeleteServiceServer
	Rep        storage.Storage
	TokenTools *tokentools.JWTTools
}

// NewDeleteService - constructor DeleteService.
func NewDeleteService(r storage.Storage, tt *tokentools.JWTTools) *DeleteService {
	return &DeleteService{Rep: r, TokenTools: tt}
}

// DelItem - delete one item.
func (d *DeleteService) DelItem(ctx context.Context, in *pb.DeleteItemReq) (*pb.DeleteResp, error) {
	resp := &pb.DeleteResp{}
	// uuid := metadatatools.GetUUIDFromMetadata(ctx)
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		resp.Error = customerror.ErrMissingToken
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken)
	}

	modelDelItem := models.IDModel{ID: in.Id, UUID: uuid, Type: in.Type}
	if err := d.Rep.DeleteRecord(ctx, modelDelItem); err != nil {
		resp.Error = err.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer)
	}
	return resp, nil
}
