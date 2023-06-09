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

// GiveTakeService - service contains methods for loading and unloading user data.
type GiveTakeService struct {
	pb.UnimplementedGiveTakeServiceServer
	Rep        storage.Storage
	TokenTools *tokentools.JWTTools
	Logger     *logrus.Logger
}

// NewGiveTakeService - constructor GiveTakeService.
func NewGiveTakeService(r storage.Storage, tt *tokentools.JWTTools, l *logrus.Logger) *GiveTakeService {
	return &GiveTakeService{Rep: r, TokenTools: tt, Logger: l}
}

// InsLogPwd - send the login and password data to the server.
func (g *GiveTakeService) InsLogPwd(ctx context.Context, in *pb.InsertLoginPasswordReq) (*pb.InsertResp, error) {
	resp := &pb.InsertResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "ins_log_pwd",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelTechData := models.ReqTechDataModel{Title: in.Title, Tag: in.Tag, Comment: in.Comment, Type: in.Type}
	modelLogPwd := models.LogPwdModel{Login: in.Login, Password: in.Password}
	modelInsLogPwd := models.ReqLogPwdModel{UUID: uuid, Data: modelLogPwd, TechData: modelTechData}

	res, err := g.Rep.InsertLogPwdPair(ctx, modelInsLogPwd)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "ins_log_pwd",
			"err":     err,
			"from":    "storage.insert_log_pwd_pair",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.ID
	resp.Title = res.Title
	return resp, nil
}

// InsCard - send the card data to the server.
func (g *GiveTakeService) InsCard(ctx context.Context, in *pb.InsertCardReq) (*pb.InsertResp, error) {
	resp := &pb.InsertResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "ins_card",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelTechData := models.ReqTechDataModel{Title: in.Title, Tag: in.Tag, Comment: in.Comment, Type: in.Type}
	modelCard := models.CardModel{Num: in.Num, Date: in.Date, CVC: in.Cvc, FirstName: in.FirstName, LastName: in.LastName}
	modelInsCard := models.ReqCardModel{UUID: uuid, Data: modelCard, TechData: modelTechData}

	res, err := g.Rep.InsertCardData(ctx, modelInsCard)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "ins_card",
			"err":     err,
			"from":    "storage.insert_card_data",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.ID
	resp.Title = res.Title
	return resp, nil
}

// InsText - send the text data to the server.
func (g *GiveTakeService) InsText(ctx context.Context, in *pb.InsertTextReq) (*pb.InsertResp, error) {
	resp := &pb.InsertResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "ins_text",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelTechData := models.ReqTechDataModel{Title: in.Title, Tag: in.Tag, Comment: in.Comment, Type: in.Type}
	modelText := models.TextDataModel{Data: in.Data}
	modelInsText := models.ReqTextModel{UUID: uuid, Data: modelText, TechData: modelTechData}

	res, err := g.Rep.InsertTextData(ctx, modelInsText)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "ins_text",
			"err":     err,
			"from":    "storage.insert_text_data",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.ID
	resp.Title = res.Title
	return resp, nil
}

// InsBinary - send the binary data to the server.
func (g *GiveTakeService) InsBinary(ctx context.Context, in *pb.InsertBinaryReq) (*pb.InsertResp, error) {
	resp := &pb.InsertResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "ins_binary",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelTechData := models.ReqTechDataModel{Title: in.Title, Tag: in.Tag, Comment: in.Comment, Type: in.Type}
	strData := hex.EncodeToString(in.Data)
	modelBinary := models.BinaryDataModel{Data: strData}
	modelInsBinary := models.ReqBinaryModel{UUID: uuid, Data: modelBinary, TechData: modelTechData}

	res, err := g.Rep.InsertBinaryData(ctx, modelInsBinary)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "ins_binary",
			"err":     err,
			"from":    "storage.insert_binary_data",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.ID
	resp.Title = res.Title
	return resp, nil
}

// GetLogPwd - get the login and password data from server.
func (g *GiveTakeService) GetLogPwd(ctx context.Context, in *pb.GetItemReq) (*pb.GetLoginPasswordResp, error) {
	resp := &pb.GetLoginPasswordResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "get_log_pwd",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelGetItem := models.IDModel{ID: in.Id, UUID: uuid}
	res, err := g.Rep.SelectLogPwdPair(ctx, modelGetItem)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "get_log_pwd",
			"err":     err,
			"from":    "storage.select_log_pwd_pair",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.TechData.ID
	resp.Title = res.TechData.Title
	resp.Login = res.Data.Login
	resp.Password = res.Data.Password
	resp.Tag = res.TechData.Tag
	resp.Comment = res.TechData.Comment
	return resp, nil
}

// GetCard - get the card data from server.
func (g *GiveTakeService) GetCard(ctx context.Context, in *pb.GetItemReq) (*pb.GetCardResp, error) {
	resp := &pb.GetCardResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "get_card",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelGetItem := models.IDModel{ID: in.Id, UUID: uuid}
	res, err := g.Rep.SelectCardData(ctx, modelGetItem)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "get_card",
			"err":     err,
			"from":    "storage.select_card_data",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.TechData.ID
	resp.Title = res.TechData.Title
	resp.Num = res.Data.Num
	resp.Date = res.Data.Date
	resp.Cvc = res.Data.CVC
	resp.FirstName = res.Data.FirstName
	resp.LastName = res.Data.LastName
	resp.Tag = res.TechData.Tag
	resp.Comment = res.TechData.Comment
	return resp, nil
}

// GetText - get the text data from server.
func (g *GiveTakeService) GetText(ctx context.Context, in *pb.GetItemReq) (*pb.GetTextResp, error) {
	resp := &pb.GetTextResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "get_text",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelGetItem := models.IDModel{ID: in.Id, UUID: uuid}
	res, err := g.Rep.SelectTextData(ctx, modelGetItem)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "get_text",
			"err":     err,
			"from":    "storage.select_text_data",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}

	resp.Id = res.TechData.ID
	resp.Title = res.TechData.Title
	resp.Data = res.Data.Data
	resp.Tag = res.TechData.Tag
	resp.Comment = res.TechData.Comment
	return resp, nil
}

// GetBinary - get the binary data from server.
func (g *GiveTakeService) GetBinary(ctx context.Context, in *pb.GetItemReq) (*pb.GetBinaryResp, error) {
	resp := &pb.GetBinaryResp{}
	uuid := ctx.Value(UUIDKey).(string)
	if uuid == "" {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "get_binary",
			"err":     customerror.ErrMissingToken.Error(),
		}).Trace("Token error")
		resp.Error = customerror.ErrMissingToken.Error()
		return resp, status.Error(codes.Unauthenticated, customerror.ErrMissingToken.Error())
	}

	modelGetItem := models.IDModel{ID: in.Id, UUID: uuid}
	res, err := g.Rep.SelectBinaryData(ctx, modelGetItem)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "get_binary",
			"err":     err,
			"from":    "storage.select_binary_data",
		}).Error("Storage error")
		resp.Error = customerror.ErrInternalServer.Error()
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}
	byteData, err := hex.DecodeString(res.Data.Data)
	if err != nil {
		g.Logger.WithFields(logrus.Fields{
			"service": "give_take_service",
			"handler": "get_binary",
			"err":     err,
			"from":    "hex.decode_string",
		}).Error("Decode error")
		return resp, status.Error(codes.Internal, customerror.ErrInternalServer.Error())
	}
	resp.Id = res.TechData.ID
	resp.Title = res.TechData.Title
	resp.Data = byteData
	resp.Tag = res.TechData.Tag
	resp.Comment = res.TechData.Comment
	return resp, nil
}
