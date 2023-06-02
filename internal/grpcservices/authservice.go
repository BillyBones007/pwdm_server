package grpcservices

import (
	"context"
	"time"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/storage/models"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	pb "github.com/BillyBones007/pwdm_service_api/api"
	"github.com/sirupsen/logrus"
)

// Authentication service.
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	Rep        storage.Storage
	TokenTools *tokentools.JWTTools
	Logger     *logrus.Logger
}

// NewAuthService - constructor AuthService.
func NewAuthService(r storage.Storage, tt *tokentools.JWTTools, l *logrus.Logger) *AuthService {
	return &AuthService{Rep: r, TokenTools: tt, Logger: l}
}

// Create - create new user.
func (a *AuthService) Create(ctx context.Context, in *pb.AuthReq) (*pb.AuthResp, error) {
	var uuid string
	user := models.UserModel{Login: in.Login, Password: in.Password}
	resp := &pb.AuthResp{}

	// checks login a new user
	existFlag, err := a.Rep.UserIsExists(ctx, user)
	if err != nil {
		a.Logger.WithFields(logrus.Fields{
			"service": "auth_service",
			"handler": "create",
			"err":     err,
			"from":    "storage.user_is_exists",
		}).Error("Storage error")
		resp.Error = customerror.ErrCreateUser.Error()
		return resp, customerror.ErrInternalServer
	}

	if existFlag {
		resp.Error = customerror.ErrCreateUser.Error()
		return resp, customerror.ErrUserIsExists
	}

	// create a new user
	uuid, err = a.Rep.CreateUser(ctx, user)
	if err != nil {
		a.Logger.WithFields(logrus.Fields{
			"service": "auth_service",
			"handler": "create",
			"err":     err,
			"from":    "storage.create_user",
		}).Error("Storage error")
		resp.Error = customerror.ErrCreateUser.Error()
		return resp, err
	}

	// creare a new token with an expiration time and uuid in claim field
	expAt := time.Now().Add(time.Hour * 1).Unix()
	token, err := a.TokenTools.CreateToken(expAt, uuid)
	if err != nil {
		a.Logger.WithFields(logrus.Fields{
			"service": "auth_service",
			"handler": "create",
			"err":     err,
			"from":    "token_tools.create_token",
		}).Error("TokenTools error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, customerror.ErrInternalServer
	}

	resp.Token = token
	return resp, nil
}

// Enter - authorization the user.
func (a *AuthService) Enter(ctx context.Context, in *pb.AuthReq) (*pb.AuthResp, error) {
	user := models.UserModel{Login: in.Login, Password: in.Password}
	resp := &pb.AuthResp{}

	ok, err := a.Rep.ValidUser(ctx, user)
	if !ok {
		a.Logger.WithFields(logrus.Fields{
			"service": "auth_service",
			"handler": "enter",
			"err":     err,
			"from":    "storage.valid_user",
		}).Error("Storage error")
		resp.Error = customerror.ErrLogIn.Error()
		return resp, err
	}

	uuid, err := a.Rep.GetUUID(ctx, user)
	if err != nil {
		a.Logger.WithFields(logrus.Fields{
			"service": "auth_service",
			"handler": "enter",
			"err":     err,
			"from":    "storage.get_uuid",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, customerror.ErrInternalServer
	}

	expAt := time.Now().Add(time.Hour * 1).Unix()
	token, err := a.TokenTools.CreateToken(expAt, uuid)
	if err != nil {
		a.Logger.WithFields(logrus.Fields{
			"service": "auth_service",
			"handler": "create",
			"err":     err,
			"from":    "token_tools.create_token",
		}).Error("TokenTools error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, customerror.ErrInternalServer
	}

	resp.Token = token
	return resp, nil
}
