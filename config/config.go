package config

import (
	env "github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Database struct {
		Driver string `env:"DATABASE_DRIVER"`
		Path   string `env:"DATABASE_PATH"`
	}

	OperatorToken string `env:"OPERATOR_TOKEN"`
}

// ReadConfig from env
func ReadConfig() *Config {
	var config Config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Database.Driver == "" {
		log.Fatal("Need DATABASE_DRIVER env var")
	}

	if config.Database.Driver != "sqlite3" && config.Database.Driver != "mysql" {
		log.Fatal("Only supported driver are sqlite3, mysql")
	}

	if config.Database.Path == "" {
		log.Fatal("Need DATABASE_PATH env var")
	}

	return &config
}
