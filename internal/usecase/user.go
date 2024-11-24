package usecase

import (
	"context"
	"dim_kurs/internal/domain/model"
	"dim_kurs/internal/domain/request"
	"dim_kurs/internal/repository"
	"dim_kurs/pkg/token"
)

type IUser interface {
	GetProfile(ctx context.Context, claims token.AuthInfo) (model.User, error)
	Update(ctx context.Context, claims token.AuthInfo, req request.UpdateUser) (model.User, error)
}

type User struct {
	userRepo repository.IUser
}

func NewUser() *User {
	return &User{}
}

func (u *User) GetProfile(ctx context.Context, claims token.AuthInfo) (model.User, error) {
	user, err := u.userRepo.Get(ctx, claims.UserID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *User) Update(ctx context.Context, req request.UpdateUser, claims token.AuthInfo) error {

	// err := u.userRepo.Update(ctx, req)
	// if err != nil {
	// 	return err
	// }

	return nil
}
