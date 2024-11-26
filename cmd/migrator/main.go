package main

import (
	"dim_kurs/internal/config"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	var migrationPath string
	var up, down bool

	godotenv.Load(".env.auth.example")

	flag.StringVar(&migrationPath, "m", "", "path to the migration dir")
	flag.BoolVar(&up, "up", false, "apply migrations")
	flag.BoolVar(&down, "down", false, "rollback migrations")
	flag.Parse()

	log := logrus.New()

	if !up && !down {
		log.Fatal("You must specify either --up or --down")
	}

	cfg, err := config.InitConfig("")
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PG.Username, cfg.PG.Password, cfg.PG.Host, cfg.PG.Port, cfg.PG.DBName)

	m, err := migrate.New("file://"+migrationPath, dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	if up {
		if err := m.Up(); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatal(err)
			}
		}
		fmt.Println("Migrations applied")
	} else if down {
		if err := m.Down(); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatal(err)
			}
		}
		fmt.Println("Migrations rolled back")
	}
}
