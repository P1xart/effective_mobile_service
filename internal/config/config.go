package config

import (
	"log/slog"
	"os"

	"github.com/P1xart/effective_mobile_service/pkg/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

const defaultConfigFile = "config.yaml"

type Config struct {
	HTTP       HTTP       `yaml:"http"`
	Postgresql Postgresql `yaml:"postgresql"`
	API        ApiUrls    `yaml:"api"`
}

type HTTP struct {
	Port string `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	Host string `yaml:"host" env:"HTTP_HOST" env-default:"127.0.0.1"`
}

type Postgresql struct {
	User     string `yaml:"user" env:"PG_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PG_PASSWORD" env-default:"5432"`
	Host     string `yaml:"host" env:"PG_HOST" env-default:"127.0.0.1"`
	Port     string `yaml:"port" env:"PG_PORT" env-default:"5432"`
	Database string `yaml:"database" env:"PG_DATABASE" env-default:"postgres"`
	SSLMode  string `yaml:"ssl_mode" env:"PG_SSL" env-default:"disable"`
}

type ApiUrls struct {
	Age    string `yaml:"age_url" env:"AGE_URL" env-default:"https://api.agify.io"`
	Gender string `yaml:"gender_url" env:"GENDER_URL" env-default:"https://api.genderize.io"`
	Nation string `yaml:"nation_url" env:"NATION_URL" env-default:"https://api.nationalize.io"`
}

func New(log *slog.Logger) (*Config, error) {
	path := fetchConfigPath()

	if _, err := os.Stat(path); err != nil {
		log.Error("failed to open config file", logger.Error(err))
		return nil, err
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	log.Debug("app configuration", slog.Any("cfg", cfg))

	return &cfg, nil
}

func fetchConfigPath() string {
	var path string

	if path = os.Getenv("CONFIG_PATH"); path == "" {
		path = defaultConfigFile
	}

	return path
}
