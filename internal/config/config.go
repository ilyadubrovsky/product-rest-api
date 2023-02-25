package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	Host     string `env:"MONGO_HOST" env-default:"localhost"`
	Port     string `env:"MONGO_PORT" env-default:"27017"`
	Database string `env:"MONGO_DATABASE" env-default:"restapi"`
	Username string `env:"MONGO_USERNAME" env-default:""`
	Password string `env:"MONGO_PASSWORD" env-default:""`
	AuthDB   string `env:"MONGO_AUTHDB" env-default:""`
}

func GetConfig() (*Config, error) {
	var err error = nil

	once.Do(func() {
		instance = &Config{}
		if err = cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			return
		}
	})

	return instance, err
}
