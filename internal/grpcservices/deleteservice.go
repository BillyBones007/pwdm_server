package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/storage/models"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	pb "github.com/BillyBones007/pwdm_service_api/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteService - service contains methods for deleting user data.
type DeleteService struct {
	pb.UnimplementedDeleteServiceServer
	Rep        storage.Storage
	TokenTools *tokentools.JWTTools
	Logger     *logrus.Logger
}

// NewDeleteService - constructor DeleteService.
func NewDeleteService(r storage.Storage, tt *tokentools.JWTTools, l *logrus.Logger) *DeleteService {
	return &DeleteService{Rep: r, TokenTools: tt, Logger: l}
}

// DelItem - delete one item.
func (d *DeleteService) DelItem(ctx context.Context, in *pb.DeleteItemReq) (*pb.DeleteResp, error) {
	resp := &pb.DeleteResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		d.Logger.WithFields(logrus.Fields{
			"service": "delete_service",
			"handler": "del_item",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelDelItem := models.IDModel{ID: in.Id, UUID: uuid, Type: in.Type}
	if err := d.Rep.DeleteRecord(ctx, modelDelItem); err != nil {
		d.Logger.WithFields(logrus.Fields{
			"service": "delete_service",
			"handler": "del_item",
			"err":     err,
			"from":    "storage.delete_record",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}
	return resp, nil
}
