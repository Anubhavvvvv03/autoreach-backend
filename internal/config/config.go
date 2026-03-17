package config

import (
    "os"
)

type Config struct {
    Port       string
    OpenAIKey  string
    HunterAPI  string
}

func LoadConfig() *Config {
    port := os.Getenv("PORT")
    if port == "" {
        port = "4040"
    }
    return &Config{
        Port:      port,
        OpenAIKey: os.Getenv("OPENAI_API_KEY"),
        HunterAPI: os.Getenv("HUNTER_API_KEY"),
    }
}
