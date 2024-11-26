package usecase

import (
	"context"
	"dim_kurs/internal/domain/model"
	"dim_kurs/internal/domain/request"
	"dim_kurs/internal/domain/response"
	"dim_kurs/pkg/token"

	"github.com/google/uuid"
)

type IAdmin interface {
	GetUsers(ctx context.Context, claims token.AuthInfo, req request.GetUsers) (response.GetUsers, error)
	GetUser(ctx context.Context, claims token.AuthInfo, id uuid.UUID) (model.User, error)
	UpdateUser(ctx context.Context, claims token.AuthInfo, req request.UpdateUser) (model.User, error)
	DeleteUser(ctx context.Context, claims token.AuthInfo, id uuid.UUID) error
	BanAkk(ctx context.Context, claims token.AuthInfo, id uuid.UUID) error
	UnBanAkk(ctx context.Context, claims token.AuthInfo, id uuid.UUID) error
	LogoutUser(ctx context.Context, claims token.AuthInfo, id uuid.UUID) error
}

type Admin struct {
}

func (u *Admin) GetUsers(ctx context.Context, claims token.AuthInfo, req request.GetUsers) (response.GetUsers, error) {
	// Заглушка: возвращаем пустой ответ и nil ошибку
	return response.GetUsers{}, nil
}

func (u *Admin) GetUser(ctx context.Context, claims token.AuthInfo, id uuid.UUID) (model.User, error) {
	// Заглушка: возвращаем пустого пользователя и nil ошибку
	return model.User{}, nil
}

func (u *Admin) UpdateUser(ctx context.Context, claims token.AuthInfo, req request.UpdateUser) (model.User, error) {
	// Заглушка: возвращаем пустого пользователя и nil ошибку
	return model.User{}, nil
}

func (u *Admin) DeleteUser(ctx context.Context, claims token.AuthInfo, id uuid.UUID) error {
	// Заглушка: возвращаем nil ошибку
	return nil
}

func (u *Admin) BanAkk(ctx context.Context, claims token.AuthInfo, id uuid.UUID) error {
	// Заглушка: возвращаем false и nil ошибку
	return nil
}

func (u *Admin) UnBanAkk(ctx context.Context, claims token.AuthInfo, id uuid.UUID) error {
	// Заглушка: возвращаем false и nil ошибку
	return nil
}

func (u *Admin) LogoutUser(ctx context.Context, claims token.AuthInfo, id uuid.UUID) error {
	return nil
}
