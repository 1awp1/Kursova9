package inject

import (
	"dim_kurs/internal/config"
	"dim_kurs/internal/usecase"
	"dim_kurs/pkg/token"

	"github.com/sirupsen/logrus"
)

type UseCases struct {
	usecase.IAuth
}

func NewUseCases(repos *Repositories, log *logrus.Logger, cfg config.Auth, tokenManager token.TokenManager) *UseCases {
	return &UseCases{
		IAuth: usecase.NewAuth(repos.IUser, log, cfg, tokenManager),
	}
}
