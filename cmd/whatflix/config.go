package main

import (
	"github.com/DTSL/golang-libraries/envutils"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type configs struct {
	Env         string `envconfig:"ENV" default:"development"`
	HTTPPort    string `envconfig:"HTTP_PORT" default:":8080"`
	DBConString string `envconfig:"MongoDB Con String" default:"mongodb://localhost:27017/"`
}

func getConfigs() (*configs, error) {
	var cfg configs
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "get configs")
	}

	err = checkConfigEnv(&cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "check configs")
	}
	return &cfg, nil
}

func checkConfigEnv(config *configs) error {
	err := envutils.Check(config.Env)
	if err != nil {
		return errors.WithMessage(err, "environment")
	}
	return nil
}
