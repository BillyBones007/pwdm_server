package customerror

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

const (
	ErrInternalServer string = "internal server error"
	ErrMissingToken   string = "missing token"
	ErrSignIn         string = "error sing in"
	ErrCreateUser     string = "create user error"
	ErrUserIsExist    string = "user is exists"
)

var (
	ErrNoRows               error = pgx.ErrNoRows
	ErrLoginOrPassIncorrect error = errors.New("login or password incorrect")
)
