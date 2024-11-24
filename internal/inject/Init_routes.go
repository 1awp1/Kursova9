package inject

import (
	"dim_kurs/internal/config"
	"dim_kurs/internal/routing"
	"dim_kurs/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func NewRouteHandlers(router *gin.Engine, pool *pgxpool.Pool, log *logrus.Logger, cfg config.Auth, tokenManager token.TokenManager) routing.Handlers {
	repositories := NewRepos(pool)
	usecases := NewUseCases(repositories, log, cfg, tokenManager)
	handlers := NewHandlers(usecases)

	routeHandlers := routing.Handlers{
		IAuth: handlers.IAuth,
	}

	return routeHandlers
}
