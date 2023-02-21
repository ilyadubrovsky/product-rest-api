package config

import (
	"github.com/ilyadubrovsky/product-rest-api/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"username"`
}

func GetConfig() (*Config, error) {
	var err error = nil

	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err = cleanenv.ReadConfig("configs/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Debug(help)
			return
		}
	})

	return instance, err
}
