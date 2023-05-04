package grpcservices

import (
	"context"
	"fmt"
	"time"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/storage/models"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	pb "github.com/BillyBones007/pwdm_service_api/api"
)

// Authentication service.
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	Rep        storage.Storage
	TokenTools *tokentools.JWTTools
}

// NewAuthService - constructor AuthService.
func NewAuthService(r storage.Storage, tt *tokentools.JWTTools) *AuthService {
	return &AuthService{Rep: r, TokenTools: tt}
}

// Create - create new user.
func (a *AuthService) Create(ctx context.Context, in *pb.AuthReq) (*pb.AuthResp, error) {
	user := models.UserModel{Login: in.Login, Password: in.Password}
	resp := &pb.AuthResp{}

	// valid new user
	exist, _ := a.Rep.ValidUser(ctx, user)
	if exist {
		resp.Error = customerror.CreateUserErr
		return resp, fmt.Errorf(customerror.UserIsExist)
	}

	// create a new user
	if err := a.Rep.CreateUser(ctx, user); err != nil {
		resp.Error = err.Error()
		return resp, err
	}

	// get uuid from new record database
	uuid, err := a.Rep.GetUUID(ctx, user)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}

	// creare a new token with an expiration time and uuid in claim field
	expAt := time.Now().Add(time.Hour * 1).Unix()
	token, err := a.TokenTools.CreateToken(expAt, uuid)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
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
		resp.Error = customerror.SignInErr
		return resp, err
	}

	uuid, err := a.Rep.GetUUID(ctx, user)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}

	expAt := time.Now().Add(time.Hour * 1).Unix()
	token, err := a.TokenTools.CreateToken(expAt, uuid)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}

	resp.Token = token
	return resp, nil
}
