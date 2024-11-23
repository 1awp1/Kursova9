package usecase

import (
	"context"
	"dim_kurs/internal/domain/model"
	"dim_kurs/internal/domain/request"
)

type IUser interface {
	GetProfile(ctx context.Context, token string) (model.User, error)
	Update(ctx context.Context, req request.UpdateUser, token string) (model.User, error)
}

type User struct{}

func (u *User) GetProfile(ctx context.Context, token string) (model.User, error) {

}

func (u *User) Update(ctx context.Context, req request.UpdateUser, token string) (model.User, error) {

}
