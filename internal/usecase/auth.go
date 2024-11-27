package usecase

import (
	"context"
	"dim_kurs/internal/config"
	"dim_kurs/internal/custom_errors"
	"dim_kurs/internal/domain/model"
	"dim_kurs/internal/domain/request"
	"dim_kurs/internal/repository"
	"dim_kurs/pkg/token"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type IAuth interface {
	Login(ctx context.Context, req request.Login) (string, error)
	Register(ctx context.Context, req request.Register) (string, error)
	VerifyToken(ctx context.Context, tkn string) (*token.AuthInfo, error)
}

type Auth struct {
	userRepo     repository.IUser
	log          *logrus.Logger
	cfg          config.Auth
	tokenManager token.TokenManager
}

func NewAuth(userRepo repository.IUser, log *logrus.Logger, cfg config.Auth, tokenManager token.TokenManager) *Auth {
	return &Auth{
		userRepo:     userRepo,
		log:          log,
		cfg:          cfg,
		tokenManager: tokenManager,
	}
}

func (u *Auth) Login(ctx context.Context, req request.Login) (string, error) {
	log := u.log.WithFields(logrus.Fields{
		"op":       "internal/usecase/auth/SignIn",
		"login":    req.Login,
		"password": req.Password,
	})

	user, err := u.userRepo.Get(ctx, req.Login)
	log.Infof("user %v", user)
	if err != nil {
		err = custom_errors.UserNotExist
		log.Error(err.Error())

		return "", err
	}
	if user == (model.User{}) {
		err = custom_errors.UserNotExist
		log.Error(err.Error())

		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)); err != nil {
		log.Error(err.Error())
		return "", err
	}

	accessToken, err := u.tokenManager.NewJWT(token.AuthInfo{
		UserID: user.ID.String(),
		Login:  *user.Login,
		Role:   *user.Role,
	})
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		log.Error(err.Error())
		return "", err
	}

	return accessToken, nil
}

func (u *Auth) Register(ctx context.Context, req request.Register) (string, error) { //TODO
	log := u.log.WithFields(logrus.Fields{
		"op":           "internal/usecase/auth/SignUp",
		"first":        req.FirstName,
		"last":         req.LastName,
		"login":        req.Login,
		"password":     req.Password,
		"phone_number": req.Phone,
		"email":        req.Email,
	})

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return "", err
	}

	passStr := string(passHash)
	role := "user"
	isOnline := true
	status := true

	user := model.User{
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		Login:     &req.Login,
		Phone:     &req.Phone,
		Email:     &req.Email,
		Password:  &passStr,
		Role:      &role,
		IsOnline:  &isOnline,
		Status:    &status,
	}

	_, err = u.userRepo.Create(ctx, user)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	createdUser, err := u.userRepo.Get(ctx, req.Login)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	accessToken, err := u.tokenManager.NewJWT(token.AuthInfo{
		UserID: createdUser.ID.String(),
		Login:  req.Login,
		Role:   "user",
	})

	return accessToken, nil
}

func (u *Auth) VerifyToken(ctx context.Context, tkn string) (*token.AuthInfo, error) {
	log := u.log.WithFields(logrus.Fields{
		"op": "internal/usecase/auth/VerifyToken",
	})

	claims, err := u.tokenManager.Parse(tkn)
	if err != nil {
		log.Errorf("access token: %v", err)
		return nil, err
	}

	return &claims, nil
}
