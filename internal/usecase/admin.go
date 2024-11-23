package usecase

import (
	"context"
	"dim_kurs/internal/domain/model"
	"dim_kurs/internal/domain/request"
	"dim_kurs/internal/domain/response"
	"dim_kurs/pkg/token"
)

type IAdmin interface {
	GetUsers(ctx context.Context, claims token.AuthInfo, req request.GetUsers) (response.GetUsers, error)
	GetUser(ctx context.Context, claims token.AuthInfo, id int64) (model.User, error)
	UpdateUser(ctx context.Context, claims token.AuthInfo, req request.UpdateUser) (model.User, error)
	DeleteUser(ctx context.Context, claims token.AuthInfo, id int64) error
	BanAkk(ctx context.Context, claims token.AuthInfo, id int64) (bool, error)
	UnBanAkk(ctx context.Context, claims token.AuthInfo, id int64) (bool, error)
}

type Admin struct {
}

func (u *Admin) GetUsers(ctx context.Context, claims token.AuthInfo, req request.GetUsers) (response.GetUsers, error) {

}

func (u *Admin) GetUser(ctx context.Context, claims token.AuthInfo, id int64) (model.User, error) {

}

func (u *Admin) UpdateUser(ctx context.Context, claims token.AuthInfo, req request.UpdateUser) (model.User, error) {

}

func (u *Admin) DeleteUser(ctx context.Context, claims token.AuthInfo, id int64) error {

}

func (u *Admin) BanAkk(ctx context.Context, claims token.AuthInfo, id int64) (bool, error) {

}

func (u *Admin) UnBanAkk(ctx context.Context, claims token.AuthInfo, id int64) (bool, error) {

}
