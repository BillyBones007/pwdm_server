package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	pb "github.com/BillyBones007/pwdm_service_api/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShowInfoService - service contains methods for getting information about uploaded user data.
type ShowInfoService struct {
	pb.UnimplementedShowInfoServiceServer
	Rep        storage.Storage
	TokenTools *tokentools.JWTTools
	Logger     *logrus.Logger
}

// NewShowInfoService - constructor ShowInfoService.
func NewShowInfoService(r storage.Storage, tt *tokentools.JWTTools, l *logrus.Logger) *ShowInfoService {
	return &ShowInfoService{Rep: r, TokenTools: tt, Logger: l}
}

// GetInfo - get information for current user.
func (s *ShowInfoService) GetInfo(ctx context.Context, in *pb.Empty) (*pb.ShowItemsResp, error) {
	resp := &pb.ShowItemsResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		s.Logger.WithFields(logrus.Fields{
			"service": "show_info_service",
			"handler": "get_info",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		return nil, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	// result list from database
	listResult, err := s.Rep.SelectAllInfoUser(ctx, uuid)
	if err != nil {
		s.Logger.WithFields(logrus.Fields{
			"service": "show_info_service",
			"handler": "get_info",
			"err":     err,
			"from":    "storage.select_all_info_user",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	listItems := make([]*pb.ShowItemsResp_ItemModel, 0)

	for _, record := range listResult {
		item := &pb.ShowItemsResp_ItemModel{
			Id:      record.ID,
			Type:    record.Type,
			Title:   record.Title,
			Tag:     record.Tag,
			Comment: record.Comment,
		}
		listItems = append(listItems, item)
	}

	resp.Items = listItems
	return resp, nil
}
