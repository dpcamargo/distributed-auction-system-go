package user_entity

import (
	"context"

	"github.com/dpcamargo/distributed-auction-system-go/internal/internal_error"
)

type User struct {
	ID   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(
		ctx context.Context, id string) (*User, *internal_error.InternalError)
}
