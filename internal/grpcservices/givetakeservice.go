package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/storage"
	pb "github.com/BillyBones007/pwdm_service_api/api"
)

// GiveTakeService - service contains methods for loading and unloading user data.
type GiveTakeService struct {
	pb.UnimplementedGiveTakeServiceServer
	Rep *storage.Storage
}

// NewGiveTakeService - constructor GiveTakeService.
func NewGiveTakeService(r *storage.Storage) *GiveTakeService {
	return &GiveTakeService{Rep: r}
}

// InsLogPwd - send the login and password data to the server.
func (g *GiveTakeService) InsLogPwd(ctx context.Context, in *pb.InsertLoginPasswordReq) (*pb.InsertResp, error) {
	// TODO: realization
	return nil, nil
}

// InsCard - send the card data to the server.
func (g *GiveTakeService) InsCard(ctx context.Context, in *pb.InsertCardReq) (*pb.InsertResp, error) {
	// TODO: realization
	return nil, nil
}

// InsText - send the text data to the server.
func (g *GiveTakeService) InsText(ctx context.Context, in *pb.InsertTextReq) (*pb.InsertResp, error) {
	// TODO: realization
	return nil, nil
}

// InsBinary - send the binary data to the server.
func (g *GiveTakeService) InsBinary(ctx context.Context, in *pb.InsertBinaryReq) (*pb.InsertResp, error) {
	// TODO: realization
	return nil, nil
}

// GetLogPwd - get the login and password data from server.
func (g *GiveTakeService) GetLogPwd(ctx context.Context, in *pb.GetItemReq) (*pb.GetLoginPasswordResp, error) {
	// TODO: realization
	return nil, nil
}

// GetCard - get the card data from server.
func (g *GiveTakeService) GetCard(ctx context.Context, in *pb.GetItemReq) (*pb.GetCardResp, error) {
	// TODO: realization
	return nil, nil
}

// GetText - get the text data from server.
func (g *GiveTakeService) GetText(ctx context.Context, in *pb.GetItemReq) (*pb.GetTextResp, error) {
	// TODO: realization
	return nil, nil
}

// GetBinary - get the binary data from server.
func (g *GiveTakeService) GetBinary(ctx context.Context, in *pb.GetItemReq) (*pb.GetBinaryResp, error) {
	// TODO: realization
	return nil, nil
}
