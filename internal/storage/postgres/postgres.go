package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/BillyBones007/pwdm_server/internal/storage/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ClientPostgres - type for working with PostgreSQL.
type ClientPostgres struct {
	Pool     *pgxpool.Pool
	ConfigCP *pgxpool.Config
}

// NewClientPostgres - returns a pointer to the ClientPostgres.
func NewClientPostgres(dst string) *ClientPostgres {
	config, err := pgxpool.ParseConfig(dst)
	if err != nil {
		log.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	cp := ClientPostgres{Pool: pool, ConfigCP: config}
	cp.createTable()
	return &cp
}

// (c *ClientPostgres) createTable - creates a table when the server starts, if one does not already exist.
func (c *ClientPostgres) createTable() error {
	m, err := migrate.New("file://migrations/postgres", c.ConfigCP.ConnString())
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != migrate.ErrNoChange {
		log.Printf("ERROR: migration up: %v\n", err)
	}
	return nil
}

// (c *ClientPostgres) Close - close the pool connections.
func (c *ClientPostgres) Close() {
	c.Pool.Close()
	fmt.Println("Pool connections is closed")
}

// CreateUser - creating a new user in database.
func (c *ClientPostgres) CreateUser(ctx context.Context, model models.UserModel) error {
	return nil
}

// ValidUser - user validation.
func (c *ClientPostgres) ValidUser(ctx context.Context, model models.UserModel) (bool, error) {
	return true, nil
}

// GetUUID - get uuid current user from database.
func (c *ClientPostgres) GetUUID(ctx context.Context, model models.UserModel) (string, error) {
	return "", nil
}

// DeleteUser - delete user from database.
func (c *ClientPostgres) DeleteUser(ctx context.Context, uuid string) error {
	return nil
}

// INsertLogPwdPair - writes the login/password pair in database.
func (c *ClientPostgres) InsertLogPwdPair(ctx context.Context, model models.ReqLogPwdModel) (models.InsertRespModel, error) {
	res := models.InsertRespModel{}
	return res, nil
}

// INsertCardData - writes the card data in database.
func (c *ClientPostgres) InsertCardData(ctx context.Context, model models.ReqCardModel) (models.InsertRespModel, error) {
	res := models.InsertRespModel{}
	return res, nil
}

// InsertTextData - writes the some text data in database.
func (c *ClientPostgres) InsertTextData(ctx context.Context, model models.ReqTextModel) (models.InsertRespModel, error) {
	res := models.InsertRespModel{}
	return res, nil
}

// InsertBinaryData - writes the some binary data in database.
func (c *ClientPostgres) InsertBinaryData(ctx context.Context, model models.ReqBinaryModel) (models.InsertRespModel, error) {
	res := models.InsertRespModel{}
	return res, nil
}

// SelectLogPwdPair - get a login/password pair from database.
func (c *ClientPostgres) SelectLogPwdPair(ctx context.Context, model models.IDModel) (models.RespLogPwdModel, error) {
	respModel := models.RespLogPwdModel{}
	return respModel, nil
}

// SelectCardData - get a card data from database.
func (c *ClientPostgres) SelectCardData(ctx context.Context, model models.IDModel) (models.RespCardModel, error) {
	respModel := models.RespCardModel{}
	return respModel, nil
}

// SelectTextData - get some text data from database.
func (c *ClientPostgres) SelectTextData(ctx context.Context, model models.IDModel) (models.RespTextModel, error) {
	respModel := models.RespTextModel{}
	return respModel, nil
}

// SelectBinaryData - get some binary data from database.
func (c *ClientPostgres) SelectBinaryData(ctx context.Context, model models.IDModel) (models.RespBinaryModel, error) {
	respModel := models.RespBinaryModel{}
	return respModel, nil
}

// SelectAllInfoUser - get all info by current user.
func (c *ClientPostgres) SelectAllInfoUser(ctx context.Context, uuid string) ([]models.DataRecordModel, error) {
	return nil, nil
}

// DeleteRecord - delete current record from database.
func (c *ClientPostgres) DeleteRecord(ctx context.Context, model models.IDModel) error {
	return nil
}

// DeleteAllRecords - delete all records specified in the list from database.
func (c *ClientPostgres) DeleteAllRecords(ctx context.Context, model models.ListRecordsModel) error {
	return nil
}
