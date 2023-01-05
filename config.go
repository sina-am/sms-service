package main

import "github.com/Netflix/go-env"

type Config struct {
	Hostname string `env:"HOSTNAME,default=:8080"`
	Database string `env:"DATABASE,default=./test.db"`
}

func NewConfig() (*Config, error) {
	config := &Config{}
	_, err := env.UnmarshalFromEnviron(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
