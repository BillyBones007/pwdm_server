package grpcservices

import (
	"context"
	"fmt"

	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Key string

const UUIDKey Key = "uuid"

// InterceptorsService - interceptors struct.
type InterceptorsService struct {
	tokenTools *tokentools.JWTTools
}

// NewInterceptorsService - constructor.
func NewInterceptorsService(tt *tokentools.JWTTools) *InterceptorsService {
	return &InterceptorsService{tokenTools: tt}
}

// AuthInterceptor - middleware for checking the token when contacting grpc.
func (i *InterceptorsService) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Printf("INFO: Called method: %v\n", info.FullMethod)
	if info.FullMethod == "/pwdm.AuthService/Create" || info.FullMethod == "/pwdm.AuthService/Enter" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("token")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	token := values[0]
	uuid, err := i.tokenTools.ParseUUID(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	newctx := context.WithValue(context.Background(), UUIDKey, uuid)
	// newMD := metadata.Pairs("uuid", uuid)
	// newctx := metadata.NewOutgoingContext(context.Background(), newMD)
	return handler(newctx, req)
}
