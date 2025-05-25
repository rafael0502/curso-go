package user_entity

import (
	"context"

	"github.com/rafael0502/curso-go/Leilao/internal/internal_error"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserRepositoryInterface interface {
	FindUserById(ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
