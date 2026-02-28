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
		Host   string `env:"DB_HOST,required"`
		Port   string `env:"DB_PORT,required"`
		DBUser string `env:"DB_USER,required"`
		DBPass string `env:"DB_PASSWORD,required"`
		DBName string `env:"DB_NAME,required"`
		URL    string `env:"DB_URL,required"`
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
