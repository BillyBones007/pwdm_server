package customerror

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

// Custom errors.
var (
	ErrNoRows               error = pgx.ErrNoRows
	ErrLoginOrPassIncorrect error = errors.New("login or password incorrect")
	ErrInternalServer       error = errors.New("internal server error")
	ErrUserIsExists         error = errors.New("user is exists")
	ErrCreateUser           error = errors.New("create user error")
	ErrLogIn                error = errors.New("error log in")
	ErrMissingToken         error = errors.New("missing token")
	ErrMissingMD            error = errors.New("missing metadata")
	ErrDSNEmpty             error = errors.New("dsn is empty")
	ErrMigrations           error = errors.New("migrations error")
)
