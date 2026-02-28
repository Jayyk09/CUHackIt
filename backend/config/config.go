package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP HTTP
		Log  Log
		DB   DB
	}
	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}
	DB struct {
		URL string `env:"DATABASE_URL,required"`
	}
)

var (
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	//this ensure that this config is only created and loaded once otherwise it just returns
	//a pointer to the currently made config
	once.Do(func() {
		_ = godotenv.Load()
		c := &Config{}
		if err := env.Parse(c); err != nil {
			panic(err)
		}
		cfg = c
	})
	return cfg
}
