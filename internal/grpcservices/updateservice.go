package grpcservices

import (
	"context"
	"encoding/hex"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/storage"
	"github.com/BillyBones007/pwdm_server/internal/storage/models"
	"github.com/BillyBones007/pwdm_server/internal/tools/tokentools"
	pb "github.com/BillyBones007/pwdm_service_api/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateService - service contains methods for updating user data.
type UpdateService struct {
	pb.UnimplementedUpdateServiceServer
	Rep        storage.Storage
	TokenTools *tokentools.JWTTools
	Logger     *logrus.Logger
}

// NewUpdateService - constructor UpdateService.
func NewUpdateService(r storage.Storage, tt *tokentools.JWTTools, l *logrus.Logger) *UpdateService {
	return &UpdateService{Rep: r, TokenTools: tt, Logger: l}
}

// UpdateLogPwd - update the login and password data on the server.
func (u *UpdateService) UpdateLogPwd(ctx context.Context, in *pb.UpdateLoginPasswordReq) (*pb.UpdateResp, error) {
	resp := &pb.UpdateResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		u.Logger.WithFields(logrus.Fields{
			"service": "update_service",
			"handler": "update_log_pwd",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelTechData := models.ReqTechDataModel{Title: in.Title, Tag: in.Tag, Comment: in.Comment, Type: in.Type}
	modelLogPwd := models.LogPwdModel{ID: in.Id, Login: in.Login, Password: in.Password}
	modelUpdLogPwd := models.ReqLogPwdModel{UUID: uuid, Data: modelLogPwd, TechData: modelTechData}

	res, err := u.Rep.UpdateLogPwdPair(ctx, modelUpdLogPwd)
	if err != nil {
		u.Logger.WithFields(logrus.Fields{
			"service": "update_service",
			"handler": "update_pwd",
			"err":     err,
			"from":    "storage.update_log_pwd_pair",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.ID
	resp.Title = res.Title
	return resp, nil

}

// UpdateCard - update the card data on the server.
func (u *UpdateService) UpdateCard(ctx context.Context, in *pb.UpdateCardReq) (*pb.UpdateResp, error) {
	resp := &pb.UpdateResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		u.Logger.WithFields(logrus.Fields{
			"service": "update_service",
			"handler": "update_card",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelTechData := models.ReqTechDataModel{Title: in.Title, Tag: in.Tag, Comment: in.Comment, Type: in.Type}
	modelCard := models.CardModel{ID: in.Id, Num: in.Num, Date: in.Date, CVC: in.Cvc, FirstName: in.FirstName, LastName: in.LastName}
	modelUpdCard := models.ReqCardModel{UUID: uuid, Data: modelCard, TechData: modelTechData}

	res, err := u.Rep.UpdateCardData(ctx, modelUpdCard)
	if err != nil {
		u.Logger.WithFields(logrus.Fields{
			"service": "update_service",
			"handler": "update_card",
			"err":     err,
			"from":    "storage.update_card_data",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.ID
	resp.Title = res.Title
	return resp, nil
}

// UpdateText - update the text data on the server.
func (u *UpdateService) UpdateText(ctx context.Context, in *pb.UpdateTextReq) (*pb.UpdateResp, error) {
	resp := &pb.UpdateResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		u.Logger.WithFields(logrus.Fields{
			"service": "update_service",
			"handler": "update_text",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelTechData := models.ReqTechDataModel{Title: in.Title, Tag: in.Tag, Comment: in.Comment, Type: in.Type}
	modelText := models.TextDataModel{ID: in.Id, Data: in.Data}
	modelUpdText := models.ReqTextModel{UUID: uuid, Data: modelText, TechData: modelTechData}

	res, err := u.Rep.UpdateTextData(ctx, modelUpdText)
	if err != nil {
		u.Logger.WithFields(logrus.Fields{
			"service": "update_service",
			"handler": "update_text",
			"err":     err,
			"from":    "storage.upadate_text_data",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.ID
	resp.Title = res.Title
	return resp, nil
}

// UpdateBinary - update the binary data on the server.
func (u *UpdateService) UpdateBinary(ctx context.Context, in *pb.UpdateBinaryReq) (*pb.UpdateResp, error) {
	resp := &pb.UpdateResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		u.Logger.WithFields(logrus.Fields{
			"service": "upadate_service",
			"handler": "update_binary",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelTechData := models.ReqTechDataModel{Title: in.Title, Tag: in.Tag, Comment: in.Comment, Type: in.Type}
	strData := hex.EncodeToString(in.Data)
	modelBinary := models.BinaryDataModel{ID: in.Id, Data: strData}
	modelUpdBinary := models.ReqBinaryModel{UUID: uuid, Data: modelBinary, TechData: modelTechData}

	res, err := u.Rep.UpdateBinaryData(ctx, modelUpdBinary)
	if err != nil {
		u.Logger.WithFields(logrus.Fields{
			"service": "update_service",
			"handler": "update_binary",
			"err":     err,
			"from":    "storage.update_binary_data",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.ID
	resp.Title = res.Title
	return resp, nil
}
