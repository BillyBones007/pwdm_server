package grpcservices

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// InterceptorsService - interceptors struct.
type InterceptorsService struct {
	tokenTools *tokentools.JWTTools
}

// NewInterceptorsService - constructor.
func NewInterceptorsService(tt *tokentools.JWTTools) *InterceptorsService {
	return &InterceptorsService{tokenTools: tt}
}

// AuthInterceptor - middleware for checking the token when contacting grpc.
func (i *InterceptorsService) AuthInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		token := values[0]
		uuid, err := i.tokenTools.ParseUUID(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		ctx = metadata.AppendToOutgoingContext(ctx, "uuid", uuid)
		return handler(ctx, req)

	}
	return nil, status.Error(codes.Unauthenticated, "missing metadata")
}
