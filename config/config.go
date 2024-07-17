package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string
	MsgTime int
}

func LoadConfig() *Config {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("File .env does not exist")
	} else {
		log.Println("File .env exists")
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file:", err)
	} else {
		log.Println(".env file loaded successfully")
	}

	return &Config{
		Port:    getEnv("PORT", "8080"),
		MsgTime: getEnvInt("MSG_TIME", 2),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
