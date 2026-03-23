package config

import (
    "os"
    "github.com/joho/godotenv"
    "github.com/yourusername/autoreach-backend/pkg/logger"
)

type Config struct {
    Port           string
    OpenAIKey      string
    HunterAPI      string
    JWTSecret      string
    DBHost         string
    DBPort         string
    DBUser         string
    DBPass         string
    DBName         string
    AWSRegion      string
    AWSBucket      string
    AWSAccessKey   string
    AWSSecretKey   string
}

var AppConfig *Config

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        logger.Info("No .env file found, using system environment variables")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "4040"
    }
    return &Config{
        Port:         port,
        OpenAIKey:    os.Getenv("OPENAI_API_KEY"),
        HunterAPI:    os.Getenv("HUNTER_API_KEY"),
        JWTSecret:    os.Getenv("JWT_SECRET"),
        DBHost:       os.Getenv("DB_HOST"),
        DBPort:       os.Getenv("DB_PORT"),
        DBUser:       os.Getenv("DB_USER"),
        DBPass:       os.Getenv("DB_PASSWORD"),
        DBName:       os.Getenv("DB_NAME"),
        AWSRegion:    os.Getenv("AWS_REGION"),
        AWSBucket:    os.Getenv("AWS_S3_BUCKET"),
        AWSAccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
        AWSSecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
    }
}
