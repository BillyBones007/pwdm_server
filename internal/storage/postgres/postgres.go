package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/BillyBones007/pwdm_server/internal/customerror"
	"github.com/BillyBones007/pwdm_server/internal/datatypes"
	"github.com/BillyBones007/pwdm_server/internal/storage/models"
	"github.com/BillyBones007/pwdm_server/internal/tools/convertuuid"
	"github.com/BillyBones007/pwdm_server/internal/tools/encpass"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

// ClientPostgres - type for working with PostgreSQL.
type ClientPostgres struct {
	Pool     *pgxpool.Pool
	ConfigCP *pgxpool.Config
	Logger   *logrus.Logger
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
	err = cp.createTable()
	if err != nil {
		log.Fatal(err)
	}
	return &cp
}

// (c *ClientPostgres) createTable - creates a table when the server starts, if one does not already exist.
func (c *ClientPostgres) createTable() error {
	m, err := migrate.New("file://migrations/postgres", c.ConfigCP.ConnString())
	if err != nil {
		c.Logger.WithField("err", err).Fatal("Migration error")
	}
	err = m.Up()
	if err != migrate.ErrNoChange {
		c.Logger.WithField("err", err).Error("Migration error")
		return err
	}
	return nil
}

// (c *ClientPostgres) Close - close the pool connections.
func (c *ClientPostgres) Close() {
	c.Pool.Close()
	c.Logger.Info("Pool connections is closed")
	fmt.Println("Pool connections is closed")
}

// CreateUser - creating a new user in database.
func (c *ClientPostgres) CreateUser(ctx context.Context, model models.UserModel) (string, error) {
	var uuid string
	exUser, err := c.UserIsExists(ctx, model)
	if err != nil {
		return "", err
	}

	if exUser {
		return "", customerror.ErrUserIsExists
	}

	encPass, err := encpass.EncPassword(model.Password)
	if err != nil {
		return "", err
	}

	q := "INSERT INTO users (uuid, login, password) VALUES (uuid_generate_v4(), $1, $2) RETURNING uuid;"
	if err := c.Pool.QueryRow(ctx, q, model.Login, encPass).Scan(&uuid); err != nil {
		return "", err
	}

	return uuid, nil
}

// ValidUser - user validation. Checks the correctness of the login and password.
func (c *ClientPostgres) ValidUser(ctx context.Context, model models.UserModel) (bool, error) {
	var encPass string
	q := "SELECT password FROM users WHERE login = $1;"
	if err := c.Pool.QueryRow(ctx, q, model.Login).Scan(&encPass); err != nil {
		if errors.Is(err, customerror.ErrNoRows) {
			return false, customerror.ErrLoginOrPassIncorrect
		}
		return false, err
	}

	if !encpass.ComparePassword(model.Password, encPass) {
		return false, customerror.ErrLoginOrPassIncorrect
	}
	return true, nil
}

// UserIsExists - Checks if the user exists.
func (c *ClientPostgres) UserIsExists(ctx context.Context, model models.UserModel) (bool, error) {
	var flag bool
	q := "SELECT EXISTS(SELECT login FROM users WHERE login = $1);"
	if err := c.Pool.QueryRow(ctx, q, model.Login).Scan(&flag); err != nil {
		return flag, err
	}
	return flag, nil
}

// GetUUID - get uuid current user from database.
func (c *ClientPostgres) GetUUID(ctx context.Context, model models.UserModel) (string, error) {
	var uuid [16]byte
	q := "SELECT uuid FROM users WHERE login = $1;"
	if err := c.Pool.QueryRow(ctx, q, model.Login).Scan(&uuid); err != nil {
		return "", err
	}

	return convertuuid.UUID(uuid).String(), nil
}

// DeleteUser - delete user from database.
func (c *ClientPostgres) DeleteUser(ctx context.Context, uuid string) error {
	q := "DELETE FROM users WHERE uuid = $1;"
	_, err := c.Pool.Exec(ctx, q, uuid)
	if err != nil {
		return err
	}
	return nil
}

// INsertLogPwdPair - writes the login/password pair in database.
func (c *ClientPostgres) InsertLogPwdPair(ctx context.Context, model models.ReqLogPwdModel) (models.InsertRespModel, error) {
	res := models.InsertRespModel{}
	var id int32
	q := `INSERT INTO log_pwd_data(uuid, type, title, login, password, tag, comment) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
	if err := c.Pool.QueryRow(ctx, q, model.UUID, model.TechData.Type, model.TechData.Title, model.Data.Login,
		model.Data.Password, model.TechData.Tag, model.TechData.Comment).Scan(&id); err != nil {
		return res, err
	}
	res.ID = id
	res.Title = model.TechData.Title
	return res, nil
}

// INsertCardData - writes the card data in database.
func (c *ClientPostgres) InsertCardData(ctx context.Context, model models.ReqCardModel) (models.InsertRespModel, error) {
	res := models.InsertRespModel{}
	var id int32
	q := `INSERT INTO card_data(uuid, type, title, num, date, cvc, first_name, last_name, tag, comment) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;`
	if err := c.Pool.QueryRow(ctx, q, model.UUID, model.TechData.Type, model.TechData.Title, model.Data.Num, model.Data.Date,
		model.Data.CVC, model.Data.FirstName, model.Data.LastName, model.TechData.Tag, model.TechData.Comment).Scan(&id); err != nil {
		return res, err
	}
	res.ID = id
	res.Title = model.TechData.Title
	return res, nil
}

// InsertTextData - writes the some text data in database.
func (c *ClientPostgres) InsertTextData(ctx context.Context, model models.ReqTextModel) (models.InsertRespModel, error) {
	res := models.InsertRespModel{}
	var id int32
	q := `INSERT INTO text_data(uuid, type, title, data, tag, comment) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	if err := c.Pool.QueryRow(ctx, q, model.UUID, model.TechData.Type, model.TechData.Title, model.Data.Data,
		model.TechData.Tag, model.TechData.Comment).Scan(&id); err != nil {
		return res, err
	}
	res.ID = id
	res.Title = model.TechData.Title
	return res, nil
}

// InsertBinaryData - writes the some binary data in database.
func (c *ClientPostgres) InsertBinaryData(ctx context.Context, model models.ReqBinaryModel) (models.InsertRespModel, error) {
	res := models.InsertRespModel{}
	var id int32
	q := `INSERT INTO binary_data(uuid, type, title, data, tag, comment) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	if err := c.Pool.QueryRow(ctx, q, model.UUID, model.TechData.Type, model.TechData.Title, model.Data.Data,
		model.TechData.Tag, model.TechData.Comment).Scan(&id); err != nil {
		return res, err
	}
	res.ID = id
	res.Title = model.TechData.Title
	return res, nil
}

// SelectLogPwdPair - get a login/password pair from database.
func (c *ClientPostgres) SelectLogPwdPair(ctx context.Context, model models.IDModel) (models.RespLogPwdModel, error) {
	res := models.RespLogPwdModel{}

	q := `SELECT login, password, title, tag, comment, type FROM log_pwd_data WHERE id = $1 AND uuid = $2 AND deleted = false;`
	if err := c.Pool.QueryRow(ctx, q, model.ID, model.UUID).Scan(&res.Data.Login, &res.Data.Password,
		&res.TechData.Title, &res.TechData.Tag, &res.TechData.Comment, &res.TechData.Type); err != nil {
		return res, err
	}

	return res, nil
}

// SelectCardData - get a card data from database.
func (c *ClientPostgres) SelectCardData(ctx context.Context, model models.IDModel) (models.RespCardModel, error) {
	res := models.RespCardModel{}

	q := `SELECT num, date, cvc, first_name, last_name, title, tag, comment, type FROM card_data WHERE id = $1 AND uuid = $2 AND deleted = false;`
	if err := c.Pool.QueryRow(ctx, q, model.ID, model.UUID).Scan(&res.Data.Num, &res.Data.Date,
		&res.Data.CVC, &res.Data.FirstName, &res.Data.LastName, &res.TechData.Title, &res.TechData.Tag,
		&res.TechData.Comment, &res.TechData.Type); err != nil {
		return res, err
	}

	return res, nil
}

// SelectTextData - get some text data from database.
func (c *ClientPostgres) SelectTextData(ctx context.Context, model models.IDModel) (models.RespTextModel, error) {
	res := models.RespTextModel{}
	q := `SELECT data, title, tag, comment, type FROM text_data WHERE id = $1 AND uuid = $2 AND deleted = false;`
	if err := c.Pool.QueryRow(ctx, q, model.ID, model.UUID).Scan(&res.Data.Data, &res.TechData.Title,
		&res.TechData.Tag, &res.TechData.Comment, &res.TechData.Type); err != nil {
		return res, err
	}

	return res, nil
}

// SelectBinaryData - get some binary data from database.
func (c *ClientPostgres) SelectBinaryData(ctx context.Context, model models.IDModel) (models.RespBinaryModel, error) {
	res := models.RespBinaryModel{}

	q := `SELECT data, title, tag, comment, type FROM binary_data WHERE id = $1 AND uuid = $2 AND deleted = false;`
	if err := c.Pool.QueryRow(ctx, q, model.ID, model.UUID).Scan(&res.Data.Data, &res.TechData.Title,
		&res.TechData.Tag, &res.TechData.Comment, &res.TechData.Type); err != nil {
		return res, err
	}

	return res, nil
}

// SelectAllInfoUser - get all info by current user.
func (c *ClientPostgres) SelectAllInfoUser(ctx context.Context, uuid string) ([]models.DataRecordModel, error) {
	res := make([]models.DataRecordModel, 0)

	q := []string{
		`SELECT title, tag, comment, type, id FROM log_pwd_data WHERE uuid = $1 AND deleted = false;`,
		`SELECT title, tag, comment, type, id FROM card_data WHERE uuid = $1 AND deleted = false;`,
		`SELECT title, tag, comment, type, id FROM text_data WHERE uuid = $1 AND deleted = false;`,
		`SELECT title, tag, comment, type, id FROM binary_data WHERE uuid = $1 AND deleted = false;`,
	}

	for _, query := range q {
		rows, err := c.Pool.Query(ctx, query, uuid)
		if err != nil {
			return res, err
		}
		for rows.Next() {
			record := models.DataRecordModel{}
			err := rows.Scan(&record.Title, &record.Tag, &record.Comment, &record.Type, &record.ID)
			if err != nil {
				return res, err
			}
			res = append(res, record)
		}
	}
	return res, nil
}

// DeleteRecord - delete current record from database.
func (c *ClientPostgres) DeleteRecord(ctx context.Context, model models.IDModel) error {
	switch model.Type {
	case datatypes.LoginPasswordDataType:
		q := `UPDATE log_pwd_data SET deleted = true WHERE id = $1 AND uuid = $2;`
		_, err := c.Pool.Exec(ctx, q, model.ID, model.UUID)
		if err != nil {
			return err
		}
		return nil
	case datatypes.CardDataType:
		q := `UPDATE card_data SET deleted = true WHERE id = $1 AND uuid = $2;`
		_, err := c.Pool.Exec(ctx, q, model.ID, model.UUID)
		if err != nil {
			return err
		}
		return nil
	case datatypes.TextDataType:
		q := `UPDATE text_data SET deleted = true WHERE id = $1 AND uuid = $2;`
		_, err := c.Pool.Exec(ctx, q, model.ID, model.UUID)
		if err != nil {
			return err
		}
		return nil
	case datatypes.BinaryDataType:
		q := `UPDATE binary_data SET deleted = true WHERE id = $1 AND uuid = $2;`
		_, err := c.Pool.Exec(ctx, q, model.ID, model.UUID)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

// DeleteAllRecords - delete all records specified in the list from database.
func (c *ClientPostgres) DeleteAllRecords(ctx context.Context, model models.ListRecordsModel) error {
	return nil
}
