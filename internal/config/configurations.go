package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	Host                 string
	User                 string
	Password             string
	DBName               string
	Port                 string
	SslMode              string
	TimeZone             string
	SecretKey            string
	MinioHost            string
	MinioRootUser        string
	MinioRootPassword    string
	MinioDefaultBucket   string
	MinioAccessKey       string
	MinioSecretAccessKey string
	MinioIsSecure        bool
	AllowedDomains       string
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
		Host:                 GetValue("POSTGRES_HOST", ""),
		User:                 GetValue("POSTGRES_USER", ""),
		Password:             GetValue("POSTGRES_PASSWORD", ""),
		DBName:               GetValue("POSTGRES_NAME", ""),
		Port:                 GetValue("POSTGRES_PORT", ""),
		SslMode:              GetValue("SSLMODE", ""),
		TimeZone:             GetValue("TIMEZONE", ""),
		SecretKey:            GetValue("SECRET_KEY", ""),
		MinioHost:            GetValue("MINIO_HOST", ""),
		MinioRootUser:        GetValue("MINIO_ROOT_USER", ""),
		MinioRootPassword:    GetValue("MINIO_ROOT_PASSWORD", ""),
		MinioDefaultBucket:   GetValue("MINIO_DEFAULT_BUCKETS", ""),
		MinioAccessKey:       GetValue("MINIO_ACCESS_KEY", ""),
		MinioSecretAccessKey: GetValue("MINIO_SECRET_ACCESS_KEY", ""),
		AllowedDomains:       GetValue("ALLOWED_DOMAINS", ""),
		MinioIsSecure:        GetBool("MINIO_IS_SECURE", false),
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

func GetBool(key string, defaultValue bool) bool {
	value_from_env := os.Getenv(key)
	value, err := strconv.ParseBool(value_from_env)
	if err != nil {
		value = defaultValue
	}
	return value
}

func GetValue(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
