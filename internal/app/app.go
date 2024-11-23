package app

import (
	"context"
	"dim_kurs/internal/config"
	"dim_kurs/internal/delivery"
	"dim_kurs/internal/repository"
	"dim_kurs/internal/server"
	"dim_kurs/internal/usecase"
	"dim_kurs/pkg/token"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type App struct {
	server *server.Server
	pool   *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg *config.Config, log *logrus.Logger) *App {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.PG.Username, cfg.PG.Password, cfg.PG.Host, cfg.PG.Port, cfg.PG.DBName)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	tokenManager, err := token.NewManager(cfg.Auth.ATDuration, cfg.Auth.PrivateKeyPath, cfg.Auth.PublicKeyPath)
	if err != nil {
		panic(err)
	}

	repos := repository.NewRepositories(pool)
	usecases := usecase.NewUsecases(repos, log, cfg.Auth, tokenManager)
	ctrs := controller.NewControllers(usecases)
	router := gin.New()

	delivery.InitRoutes(router, ctrs)
	server := server.NewServer(&cfg.Server, router)

	return &App{
		server: server,
		pool:   pool,
	}
}

func (a *App) Start() {
	a.server.Run()
}

func (a *App) Stop(ctx context.Context) {
	a.server.Stop(ctx)
	a.pool.Close()
}
