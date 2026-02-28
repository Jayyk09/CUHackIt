package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		HTTP   HTTP
		Log    Log
		DB     DB
		Auth0  Auth0
		Gemini Gemini
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
	Auth0 struct {
		Domain        string `env:"AUTH0_DOMAIN,required"`
		ClientID      string `env:"AUTH0_CLIENT_ID,required"`
		ClientSecret  string `env:"AUTH0_CLIENT_SECRET,required"`
		CallbackURL   string `env:"AUTH0_CALLBACK_URL,required"`
		SessionSecret string `env:"SESSION_SECRET,required"`
	}
	Gemini struct {
		APIKey string `env:"GEMINI_API_KEY"`
		Model  string `env:"GEMINI_MODEL" envDefault:"gemini-1.5-flash"`
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
