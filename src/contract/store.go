package contract

import (
	"biatosh/entity"
	"context"
)

type UserStore interface {
	LoginUser(ctx context.Context, user *entity.User) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUser(ctx context.Context, id int) (*entity.User, error)
	GetUsers(ctx context.Context) ([]*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type Store interface {
	UserStore
}
