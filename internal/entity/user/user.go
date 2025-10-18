package user

import (
	"context"

	"github.com/ivandersr/go-auction/internal/internal_errors"
)

type User struct {
	Id   string
	Name string
}

type UserRepostitoryInterface interface {
	FindUserById(ctx context.Context, id string) (*User, *internal_errors.InternalError)
}
