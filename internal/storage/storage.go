package storage

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/storage/models"
)

// Storage - interface for working with the storage.
type Storage interface {
	CreateUser(ctx context.Context, model models.UserModel) error
	DeleteUser(ctx context.Context, uuid string) error
	InsertLogPwdPair(ctx context.Context, model models.ReqLogPwdModel) error
	InsertCardData(ctx context.Context, model models.ReqCardModel) error
	InsertTextData(ctx context.Context, model models.ReqTextModel) error
	InsertBinaryData(ctx context.Context, model models.ReqBinaryModel) error
	SelectLogPwdPair(ctx context.Context, model models.IDModel) models.RespLogPwdModel
	SelectCardData(ctx context.Context, model models.IDModel) models.RespCardModel
	SelectTextData(ctx context.Context, model models.IDModel) models.RespTextModel
	SelectBinaryData(ctx context.Context, model models.IDModel) models.RespBinaryModel
	SelectAllInfoUser(ctx context.Context, uuid string) []models.DataRecordModel
	DeleteRecord(ctx context.Context, model models.IDModel) error
	DeleteAllRecords(ctx context.Context, model models.ListRecordsModel) error
	Close()
}
