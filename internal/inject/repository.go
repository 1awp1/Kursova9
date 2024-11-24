package inject

import (
	"dim_kurs/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	repository.IUser
}

func NewRepos(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		IUser: repository.NewUser(pool),
	}
}
