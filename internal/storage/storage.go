package storage

import (
	"context"

	"github.com/BillyBones007/pwdm_server/internal/storage/models"
)

// Storage - interface for working with the storage.
type Storage interface {
	CreateUser(ctx context.Context, model models.UserModel) error
	ValidUser(ctx context.Context, model models.UserModel) (bool, error)
	GetUUID(ctx context.Context, model models.UserModel) (string, error)
	DeleteUser(ctx context.Context, uuid string) error
	InsertLogPwdPair(ctx context.Context, model models.ReqLogPwdModel) (models.InsertRespModel, error)
	InsertCardData(ctx context.Context, model models.ReqCardModel) (models.InsertRespModel, error)
	InsertTextData(ctx context.Context, model models.ReqTextModel) (models.InsertRespModel, error)
	InsertBinaryData(ctx context.Context, model models.ReqBinaryModel) (models.InsertRespModel, error)
	SelectLogPwdPair(ctx context.Context, model models.IDModel) (models.RespLogPwdModel, error)
	SelectCardData(ctx context.Context, model models.IDModel) (models.RespCardModel, error)
	SelectTextData(ctx context.Context, model models.IDModel) (models.RespTextModel, error)
	SelectBinaryData(ctx context.Context, model models.IDModel) (models.RespBinaryModel, error)
	SelectAllInfoUser(ctx context.Context, uuid string) ([]models.DataRecordModel, error)
	DeleteRecord(ctx context.Context, model models.IDModel) error
	DeleteAllRecords(ctx context.Context, model models.ListRecordsModel) error
	Close()
}
