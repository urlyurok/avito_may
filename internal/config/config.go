package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env"  envDefault:"development"`
	LogLevel   string `yaml:"log_level" env-default:"debug"`
	HTTPServer `yaml:"http_server"`
	Postgres   `yaml:"postgres"`
}

type HTTPServer struct {
	Adress      string        `yaml:"port" envDefault:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Postgres struct {
	URL string `env-required:"true" env:"POSTGRES_CONN"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file was not found by this path %s: ", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	err := cleanenv.UpdateEnv(cfg)
	if err != nil {
		log.Fatalf("cannot update config: %s", err)
	}

	return cfg
}
