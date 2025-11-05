package apperror

import (
	"errors"

	"github.com/diagnosis/luxsuv-api-v2/internal/store"
	"github.com/jackc/pgx/v5"
)

func MapDBError(err error) *AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, store.ErrNotFound) || errors.Is(err, pgx.ErrNoRows) {
		return NotFound("Resource not found")
	}

	if errors.Is(err, store.ErrDuplicateEmail) {
		return EmailAlreadyExists()
	}

	if errors.Is(err, store.ErrTokenInvalid) {
		return Unauthorized("Token is invalid or expired")
	}

	return Wrap(CodeDatabaseError, "Database operation failed", 500, err)
}
