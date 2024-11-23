package usecase

import (
	"context"
	"dim_kurs/internal/domain/request"
)

type IAuth interface {
	Login(ctx context.Context, req request.Login) (string, error)
	Register(ctx context.Context, req request.Register) (string, error)
}

type Auth struct{}

func (u *Auth) Login(ctx context.Context, req request.Login) (string, error) {

}

func (u *Auth) Register(ctx context.Context, req request.Register) (string, error) {

}
