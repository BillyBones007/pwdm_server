package grpcservices

import (
	"context"
	"fmt"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	"github.com/sirupsen/logrus"
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
	Logger     *logrus.Logger
}

// NewInterceptorsService - constructor.
func NewInterceptorsService(tt *tokentools.JWTTools, l *logrus.Logger) *InterceptorsService {
	return &InterceptorsService{tokenTools: tt, Logger: l}
}

// AuthInterceptor - middleware for checking the token when contacting grpc.
func (i *InterceptorsService) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Printf("INFO: Called method: %v\n", info.FullMethod)
	if info.FullMethod == "/pwdm.AuthService/Create" || info.FullMethod == "/pwdm.AuthService/Enter" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		i.Logger.WithFields(logrus.Fields{
			"service": "interceptors_service",
			"handler": "auth_interceptor",
			"err":     customerror.ErrMissingMD.Error(),
		}).Trace("Metadata error")
		return nil, status.Error(codes.Unauthenticated, customerror.ErrMissingMD.Error())
	}

	values := md.Get("token")
	if len(values) == 0 {
		i.Logger.WithFields(logrus.Fields{
			"service": "interceptors_service",
			"handler": "auth_interceptor",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		return nil, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	token := values[0]
	uuid, err := i.tokenTools.ParseUUID(token)
	if err != nil {
		i.Logger.WithFields(logrus.Fields{
			"service": "interceptors_service",
			"handler": "auth_interceptor",
			"err":     err,
			"from":    "token_tools.parse_uuid",
		}).Error("TokenTools error")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	newctx := context.WithValue(context.Background(), UUIDKey, uuid)
	return handler(newctx, req)
}
