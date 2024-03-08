package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Environment struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SslMode  string
	TimeZone string
}

type Config struct {
	Envs Environment
}

type Configuration interface {
	GetDBConnectionUri() string
}

func NewConfiguration() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	envs := Environment{
		Host:     os.Getenv("HOST"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		DBName:   os.Getenv("DBNAME"),
		Port:     os.Getenv("PORT"),
		SslMode:  os.Getenv("SSLMODE"),
		TimeZone: os.Getenv("TIMEZONE"),
	}

	return &Config{Envs: envs}, nil
}

func (c *Config) GetDBConnectionUri() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Envs.Host, c.Envs.User, c.Envs.Password, c.Envs.DBName, c.Envs.Port, c.Envs.SslMode)
}
