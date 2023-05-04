package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/storage/models"
	"github.com/BillyBones007/pwdm_server/internal/tools/metadatatools"
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
	uuid := metadatatools.GetUUIDFromMetadata(ctx)
	if uuid == "" {
		resp.Error = customerror.MissingTokenErr
		return resp, status.Error(codes.Unauthenticated, customerror.MissingTokenErr)
	}

	modelDelItem := models.IDModel{ID: in.Id, UUID: uuid}
	if err := d.Rep.DeleteRecord(ctx, modelDelItem); err != nil {
		resp.Error = err.Error()
		return resp, status.Error(codes.Internal, customerror.InternalServerErr)
	}
	return resp, nil
}

// DelAll - delete all data for current user.
func (d *DeleteService) DelAll(ctx context.Context, in *pb.DeleteAllReq) (*pb.DeleteResp, error) {
	resp := &pb.DeleteResp{}
	uuid := metadatatools.GetUUIDFromMetadata(ctx)
	if uuid == "" {
		resp.Error = customerror.MissingTokenErr
		return resp, status.Error(codes.Unauthenticated, customerror.MissingTokenErr)
	}

	modelDelAll := models.ListRecordsModel{ListID: in.ListId, UUID: uuid}
	if err := d.Rep.DeleteAllRecords(ctx, modelDelAll); err != nil {
		resp.Error = err.Error()
		return resp, status.Error(codes.Internal, customerror.InternalServerErr)
	}
	return resp, nil
}
