package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Jobs struct {
	ExpirePeriod   time.Duration `env:"JOBS_EXPIRE_PERIOD" env-default:"1m"`
	AutoJoinPeriod time.Duration `env:"JOBS_AUTO_JOIN_PERIOD" env-default:"30s"`
}

func NewJobs() (*Jobs, error) {
	var jobs Jobs
	if err := cleanenv.ReadEnv(&jobs); err != nil {
		return nil, err
	}

	return &jobs, nil
}
