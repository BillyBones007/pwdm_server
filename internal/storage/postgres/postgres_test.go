package postgres

import (
	"context"
	"testing"

	"github.com/BillyBones007/pwdm_server/internal/datatypes"
	"github.com/BillyBones007/pwdm_server/internal/storage/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

const (
	dsn             string = "postgres://admin:root@localhost:5432/test?sslmode=disable"
	extension       string = `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	createUserTable string = `CREATE TABLE IF NOT EXISTS users(uuid UUID UNIQUE NOT NULL PRIMARY KEY, 
		login VARCHAR(255) UNIQUE NOT NULL, password VARCHAR(255) NOT NULL, deleted BOOLEAN DEFAULT false);`
	createLPTable string = `CREATE TABLE IF NOT EXISTS log_pwd_data(uuid UUID NOT NULL, id SERIAL UNIQUE NOT
		 NULL PRIMARY KEY, type INTEGER NOT NULL, title VARCHAR(255), login VARCHAR(255), 
		 password VARCHAR(255), tag VARCHAR(255), comment TEXT, deleted BOOLEAN DEFAULT false);`
	createCardTable string = `CREATE TABLE IF NOT EXISTS card_data(uuid UUID NOT NULL, id SERIAL UNIQUE NOT 
		NULL PRIMARY KEY, type INTEGER NOT NULL, title VARCHAR(255), num VARCHAR(255), date VARCHAR(255),
		 cvc VARCHAR(255), first_name VARCHAR(255), last_name VARCHAR(255), tag VARCHAR(255), comment TEXT, deleted BOOLEAN DEFAULT false);`
	createTextTable string = `CREATE TABLE IF NOT EXISTS text_data(uuid UUID NOT NULL, id SERIAL UNIQUE NOT
		 NULL PRIMARY KEY, type INTEGER NOT NULL, title VARCHAR(255), data TEXT, tag VARCHAR(255),
		  comment TEXT, deleted BOOLEAN DEFAULT false);`
	createBinaryTable string = `CREATE TABLE IF NOT EXISTS binary_data(uuid UUID NOT NULL, id SERIAL UNIQUE NOT
		 NULL PRIMARY KEY, type INTEGER NOT NULL, title VARCHAR(255), data TEXT, tag VARCHAR(255), 
		 comment TEXT, deleted BOOLEAN DEFAULT false);`
	dropUserTable   string = "DROP TABLE IF EXISTS users;"
	dropLPTable     string = "DROP TABLE IF EXISTS log_pwd_data;"
	dropCardTable   string = "DROP TABLE IF EXISTS card_data;"
	dropTextTable   string = "DROP TABLE IF EXISTS text_data;"
	dropBinaryTable string = "DROP TABLE IF EXISTS binary_data;"
)

func NewTestClient(dsn string) (*ClientPostgres, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	cp := ClientPostgres{Pool: pool, ConfigCP: config}
	err = createTestTables(cp.Pool)
	if err != nil {
		return nil, err
	}
	return &cp, nil
}

func createTestTables(pool *pgxpool.Pool) error {
	tables := []string{extension, createUserTable, createLPTable, createCardTable, createTextTable, createBinaryTable}
	for _, table := range tables {
		conn, err := pool.Acquire(context.TODO())
		if err != nil {
			return err
		}
		_, err = conn.Exec(context.TODO(), table)
		if err != nil {
			return err
		}
		conn.Release()
	}
	return nil
}

func dropTestTables(pool *pgxpool.Pool) error {
	tables := []string{dropUserTable, dropLPTable, dropCardTable, dropTextTable, dropBinaryTable}
	for _, table := range tables {
		conn, err := pool.Acquire(context.TODO())
		if err != nil {
			return err
		}
		_, err = conn.Exec(context.TODO(), table)
		if err != nil {
			return err
		}
		conn.Release()
	}
	return nil
}

func TestStorage(t *testing.T) {
	t.Run("Create user", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "TestLogin", Password: "TestPassword"}
		_, err = client.CreateUser(ctx, args)
		assert.NoError(t, err)
	})

	t.Run("Bad create user", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "TestLogin", Password: "TestPassword"}
		_, err = client.CreateUser(ctx, args)
		assert.NoError(t, err)

		_, err = client.CreateUser(ctx, args)
		assert.Error(t, err)
	})

	t.Run("Valid user", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "TestLogin", Password: "TestPassword"}
		_, err = client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		ok, _ := client.ValidUser(ctx, args)
		assert.True(t, ok)
	})

	t.Run("Not valid user", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "TestLogin", Password: "TestPassword"}
		_, err = client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		newArgs := models.UserModel{Login: "Baduser", Password: "1234"}
		ok, _ := client.ValidUser(ctx, newArgs)
		assert.False(t, ok)
	})

	t.Run("User is exists", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "TestLogin", Password: "TestPassword"}
		_, err = client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		ok, _ := client.UserIsExists(ctx, args)
		assert.True(t, ok)
	})

	t.Run("Get UUID", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "TestLogin", Password: "TestPassword"}
		_, err = client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		uuid, _ := client.GetUUID(ctx, args)
		assert.NotEmpty(t, uuid)
	})

	t.Run("Delete user", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "TestLogin", Password: "TestPassword"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}

		err = client.DeleteUser(ctx, uuid)
		assert.NoError(t, err)

		ok, _ := client.UserIsExists(ctx, args)
		assert.False(t, ok)
	})

	t.Run("Insert login/password", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqLogPwdModel{UUID: uuid, Data: models.LogPwdModel{
			Login:    "Test",
			Password: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.LoginPasswordDataType,
		},
		}

		resp, _ := client.InsertLogPwdPair(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)
	})

	t.Run("Update login/password", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqLogPwdModel{UUID: uuid, Data: models.LogPwdModel{
			Login:    "Test",
			Password: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.LoginPasswordDataType,
		},
		}

		resp, _ := client.InsertLogPwdPair(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)

		newData := models.ReqLogPwdModel{UUID: uuid, Data: models.LogPwdModel{
			Login:    "Test",
			Password: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "New title",
			Type:  datatypes.LoginPasswordDataType,
		},
		}

		newResp, _ := client.UpdateLogPwdPair(ctx, newData)
		assert.NotEqual(t, data.TechData.Title, newResp.Title)
	})

	t.Run("Insert card data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqCardModel{UUID: uuid, Data: models.CardModel{
			Num:       "1234 1234 1234 1234",
			Date:      "04/23",
			CVC:       "123",
			FirstName: "Ivan",
			LastName:  "Ivanov",
		}, TechData: models.ReqTechDataModel{
			Title: "Ivan card",
			Type:  datatypes.CardDataType,
		},
		}

		resp, _ := client.InsertCardData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)
	})

	t.Run("Update card data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqCardModel{UUID: uuid, Data: models.CardModel{
			Num:       "1234 1234 1234 1234",
			Date:      "04/23",
			CVC:       "123",
			FirstName: "Ivan",
			LastName:  "Ivanov",
		}, TechData: models.ReqTechDataModel{
			Title: "Ivan card",
			Type:  datatypes.CardDataType,
		},
		}

		resp, _ := client.InsertCardData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)

		newData := models.ReqCardModel{UUID: uuid, Data: models.CardModel{
			Num:       "1234 1234 1234 1234",
			Date:      "04/23",
			CVC:       "123",
			FirstName: "Hren morjovi",
			LastName:  "Ivanov",
		}, TechData: models.ReqTechDataModel{
			Title: "Hren",
			Type:  datatypes.CardDataType,
		},
		}

		newResp, _ := client.UpdateCardData(ctx, newData)
		assert.NotEqual(t, data.TechData.Title, newResp.Title)
	})

	t.Run("Insert text data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqTextModel{UUID: uuid, Data: models.TextDataModel{
			Data: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.TextDataType,
		},
		}

		resp, _ := client.InsertTextData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)
	})

	t.Run("Update text data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqTextModel{UUID: uuid, Data: models.TextDataModel{
			Data: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.TextDataType,
		},
		}

		resp, _ := client.UpdateTextData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)

		newData := models.ReqTextModel{UUID: uuid, Data: models.TextDataModel{
			Data: "new test",
		}, TechData: models.ReqTechDataModel{
			Title: "New title",
			Type:  datatypes.TextDataType,
		},
		}

		newResp, _ := client.UpdateTextData(ctx, newData)
		assert.NotEqual(t, data.TechData.Title, newResp.Title)
	})

	t.Run("Insert binary data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqBinaryModel{UUID: uuid, Data: models.BinaryDataModel{
			Data: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.BinaryDataType,
		},
		}

		resp, _ := client.InsertBinaryData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)
	})

	t.Run("Update binary data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqBinaryModel{UUID: uuid, Data: models.BinaryDataModel{
			Data: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.BinaryDataType,
		},
		}

		resp, _ := client.UpdateBinaryData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)

		newData := models.ReqBinaryModel{UUID: uuid, Data: models.BinaryDataModel{
			Data: "new test",
		}, TechData: models.ReqTechDataModel{
			Title: "New title",
			Type:  datatypes.BinaryDataType,
		},
		}

		newResp, _ := client.UpdateBinaryData(ctx, newData)
		assert.NotEqual(t, data.TechData.Title, newResp.Title)
	})

	t.Run("Select login/password", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqLogPwdModel{UUID: uuid, Data: models.LogPwdModel{
			Login:    "Test",
			Password: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.LoginPasswordDataType,
		},
		}

		resp, _ := client.InsertLogPwdPair(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)

		selResp, _ := client.SelectLogPwdPair(ctx, models.IDModel{UUID: uuid, ID: resp.ID, Type: datatypes.LoginPasswordDataType})
		assert.Equal(t, selResp.Data.Login, data.Data.Login)
	})

	t.Run("Select card data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqCardModel{UUID: uuid, Data: models.CardModel{
			Num:       "1234 1234 1234 1234",
			Date:      "04/23",
			CVC:       "123",
			FirstName: "Ivan",
			LastName:  "Ivanov",
		}, TechData: models.ReqTechDataModel{
			Title: "Ivan card",
			Type:  datatypes.CardDataType,
		},
		}

		resp, _ := client.InsertCardData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)

		selResp, _ := client.SelectCardData(ctx, models.IDModel{UUID: uuid, ID: resp.ID, Type: datatypes.CardDataType})
		assert.Equal(t, selResp.Data.Num, data.Data.Num)
	})

	t.Run("Select text data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqTextModel{UUID: uuid, Data: models.TextDataModel{
			Data: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.TextDataType,
		},
		}

		resp, _ := client.InsertTextData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)

		selResp, _ := client.SelectTextData(ctx, models.IDModel{UUID: uuid, ID: resp.ID, Type: datatypes.TextDataType})
		assert.Equal(t, selResp.Data.Data, data.Data.Data)
	})

	t.Run("Select binary data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqBinaryModel{UUID: uuid, Data: models.BinaryDataModel{
			Data: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.BinaryDataType,
		},
		}

		resp, _ := client.InsertBinaryData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)

		selResp, _ := client.SelectBinaryData(ctx, models.IDModel{UUID: uuid, ID: resp.ID, Type: datatypes.BinaryDataType})
		assert.Equal(t, selResp.Data.Data, data.Data.Data)
	})

	t.Run("Select all data", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}

		lpData := models.ReqLogPwdModel{UUID: uuid, Data: models.LogPwdModel{
			Login:    "Test",
			Password: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Login/Password",
			Type:  datatypes.LoginPasswordDataType,
		},
		}

		lpResp, _ := client.InsertLogPwdPair(ctx, lpData)
		assert.Equal(t, lpData.TechData.Title, lpResp.Title)

		cData := models.ReqCardModel{UUID: uuid, Data: models.CardModel{
			Num:       "1234 1234 1234 1234",
			Date:      "04/23",
			CVC:       "123",
			FirstName: "Ivan",
			LastName:  "Ivanov",
		}, TechData: models.ReqTechDataModel{
			Title: "Ivan card",
			Type:  datatypes.CardDataType,
		},
		}

		cResp, _ := client.InsertCardData(ctx, cData)
		assert.Equal(t, cData.TechData.Title, cResp.Title)

		tData := models.ReqTextModel{UUID: uuid, Data: models.TextDataModel{
			Data: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Text",
			Type:  datatypes.TextDataType,
		},
		}

		tResp, _ := client.InsertTextData(ctx, tData)
		assert.Equal(t, tData.TechData.Title, tResp.Title)

		bData := models.ReqBinaryModel{UUID: uuid, Data: models.BinaryDataModel{
			Data: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Binary",
			Type:  datatypes.BinaryDataType,
		},
		}

		bResp, _ := client.InsertBinaryData(ctx, bData)
		assert.Equal(t, bData.TechData.Title, bResp.Title)

		allResp, err := client.SelectAllInfoUser(ctx, uuid)
		if err != nil {
			t.Fail()
		}
		for _, v := range allResp {
			switch v.Type {
			case datatypes.LoginPasswordDataType:
				assert.Equal(t, v.Title, lpData.TechData.Title)
			case datatypes.CardDataType:
				assert.Equal(t, v.Title, cData.TechData.Title)
			case datatypes.TextDataType:
				assert.Equal(t, v.Title, tData.TechData.Title)
			case datatypes.BinaryDataType:
				assert.Equal(t, v.Title, bData.TechData.Title)
			}
		}

	})

	t.Run("Delete record", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}
		data := models.ReqBinaryModel{UUID: uuid, Data: models.BinaryDataModel{
			Data: "test",
		}, TechData: models.ReqTechDataModel{
			Title: "Title",
			Type:  datatypes.BinaryDataType,
		},
		}

		resp, _ := client.InsertBinaryData(ctx, data)
		assert.Equal(t, data.TechData.Title, resp.Title)

		err = client.DeleteRecord(ctx, models.IDModel{ID: resp.ID, UUID: uuid, Type: datatypes.BinaryDataType})
		assert.NoError(t, err)

		_, err = client.SelectBinaryData(ctx, models.IDModel{UUID: uuid, ID: resp.ID, Type: datatypes.BinaryDataType})
		assert.Error(t, err)
	})

	t.Run("Delete not exist record", func(t *testing.T) {
		client, err := NewTestClient(dsn)
		ctx := context.TODO()
		defer dropTestTables(client.Pool)
		if err != nil {
			t.Fatalf("Failed create client: %v", err)
		}
		args := models.UserModel{Login: "User", Password: "1234"}
		uuid, err := client.CreateUser(ctx, args)
		if err != nil {
			t.Fail()
		}

		err = client.DeleteRecord(ctx, models.IDModel{ID: 134, UUID: uuid, Type: datatypes.BinaryDataType})
		assert.NoError(t, err)
	})

	t.Run("NewClientPostgres empty dsn", func(t *testing.T) {
		_, err := NewClientPostgres("")
		assert.Error(t, err)
	})

	t.Run("NewClientPostgres bad dsn", func(t *testing.T) {
		_, err := NewClientPostgres("1234")
		assert.Error(t, err)
	})
}
