package config

import "github.com/ilyakaznacheev/cleanenv"

type Database struct {
	DriverName     string `env:"DB_DRIVER" env-default:"postgres"`
	DataSourceName string `env:"DB_DSN" env-default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	MigrationsPath string `env:"DB_MIGRATIONS_PATH" env-default:"file://./migrations/postgres"`
}

func NewDatabase() (*Database, error) {
	var cfg Database
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
