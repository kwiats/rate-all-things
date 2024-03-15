package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	Host      string
	User      string
	Password  string
	DBName    string
	Port      string
	SslMode   string
	TimeZone  string
	SecretKey string
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
		Host:      os.Getenv("POSTGRES_HOST"),
		User:      os.Getenv("POSTGRES_USER"),
		Password:  os.Getenv("POSTGRES_PASSWORD"),
		DBName:    os.Getenv("POSTGRES_NAME"),
		Port:      os.Getenv("POSTGRES_PORT"),
		SslMode:   os.Getenv("SSLMODE"),
		TimeZone:  os.Getenv("TIMEZONE"),
		SecretKey: os.Getenv("SECRET_KEY"),
	}

	return &Config{Envs: envs}, nil
}

func (c *Config) GetDBConnectionUri() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Envs.Host, c.Envs.User, c.Envs.Password, c.Envs.DBName, c.Envs.Port, c.Envs.SslMode)
}

func (c *Config) GetSecretKey() string {
	return c.Envs.SecretKey
}




