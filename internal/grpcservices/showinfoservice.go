package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/tools/metadatatools"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	pb "github.com/BillyBones007/pwdm_service_api/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShowInfoService - service contains methods for getting information about uploaded user data.
type ShowInfoService struct {
	pb.UnimplementedShowInfoServiceServer
	Rep        storage.Storage
	TokenTools *tokentools.JWTTools
}

// NewShowInfoService - constructor ShowInfoService.
func NewShowInfoService(r storage.Storage, tt *tokentools.JWTTools) *ShowInfoService {
	return &ShowInfoService{Rep: r, TokenTools: tt}
}

// GetInfo - get information for current user.
func (s *ShowInfoService) GetInfo(ctx context.Context, in *pb.Empty) (*pb.ShowItemsResp, error) {
	resp := &pb.ShowItemsResp{}
	uuid := metadatatools.GetUUIDFromMetadata(ctx)
	if uuid == "" {
		return nil, status.Error(codes.Unauthenticated, customerror.MissingTokenErr)
	}

	// result list from database
	listResult, err := s.Rep.SelectAllInfoUser(ctx, uuid)
	if err != nil {
		resp.Error = err.Error()
		return resp, status.Error(codes.Internal, customerror.InternalServerErr)
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
