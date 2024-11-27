package usecase

import (
	"context"
	"dim_kurs/internal/domain/model"
	"dim_kurs/internal/domain/request"
	"dim_kurs/internal/repository"
	"dim_kurs/pkg/token"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type IUser interface {
	GetProfile(ctx context.Context, claims token.AuthInfo) (model.User, error)
	Update(ctx context.Context, claims token.AuthInfo, req request.UpdateUser) error
	Logout(ctx context.Context, claims token.AuthInfo) error
}

type User struct {
	userRepo repository.IUser
	log      *logrus.Logger
}

func NewUser(userRepo repository.IUser, log *logrus.Logger) *User {
	return &User{
		userRepo: userRepo,
		log:      log,
	}
}

func (u *User) GetProfile(ctx context.Context, claims token.AuthInfo) (model.User, error) {
	log := u.log.WithFields(logrus.Fields{
		"op": "internal/usecase/user/GetProfile",
	})

	user, err := u.userRepo.Get(ctx, claims.Login)
	if err != nil {
		log.Error(err)

		return model.User{}, err
	}

	return user, nil
}

func (u *User) Update(ctx context.Context, claims token.AuthInfo, req request.UpdateUser) error {
	log := u.log.WithFields(logrus.Fields{
		"op": "internal/usecase/user/Update",
	})

	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		log.Error(err)
		return err
	}

	var passHash *string
	if req.Password != nil && *req.Password != "" {
		p, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(err)
			return err
		}

		kosiak := string(p)
		passHash = &kosiak
	}
	user := model.User{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Login:     req.Login,
		Phone:     req.Phone,
		Email:     req.Email,
		Password:  passHash,
	}

	err = u.userRepo.Update(ctx, user)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *User) Logout(ctx context.Context, claims token.AuthInfo) error {
	log := u.log.WithFields(logrus.Fields{
		"op": "internal/usecase/user/Logout",
	})
	isOnline := false

	err := u.userRepo.Update(ctx, model.User{
		ID:       uuid.MustParse(claims.UserID),
		IsOnline: &isOnline,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
