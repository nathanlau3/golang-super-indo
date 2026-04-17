package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort           string `mapstructure:"APP_PORT"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	RedisAddr         string `mapstructure:"REDIS_ADDR"`
	RedisPassword     string `mapstructure:"REDIS_PASSWORD"`
	JWTPrivateKeyPath string `mapstructure:"JWT_PRIVATE_KEY_PATH"`
	JWTPublicKeyPath  string `mapstructure:"JWT_PUBLIC_KEY_PATH"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "superindo")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("JWT_PRIVATE_KEY_PATH", "pkg/credentials/jwtRS256.key")
	viper.SetDefault("JWT_PUBLIC_KEY_PATH", "pkg/credentials/jwtRS256.key.pub")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("file .env tidak ditemukan, pakai env variable / default: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("gagal unmarshal config: %v", err)
	}

	return &cfg
}
