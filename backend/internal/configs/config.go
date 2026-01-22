package configs

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Database databaseConfig
	Backend  backendConfig
}

type databaseConfig struct {
	Username string
	Password string
	Hostname string
	Port     int
	Name     string
}

type backendConfig struct {
	Port int
}

var instance *Config
var once sync.Once

func getRequiredValue(key string) string {
	value := getValue(key)
	if value == "" {
		panic(fmt.Sprintf("an required value \"%s\" is not in .env file", key))
	}

	return value
}

func getValueToInt(key string) int {
	value := getValue(key)
	if value == "" {
		return 0
	}

	parsedInt, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("an value \"%s\" should to be an int", key))
	}

	return parsedInt
}

func getValue(key string) string {
	return os.Getenv(key)
}

func GetConfig() *Config {
	once.Do(func() {
		godotenv.Load()

		instance = &Config{
			Database: databaseConfig{
				Username: getRequiredValue("DATABASE_USERNAME"),
				Password: getRequiredValue("DATABASE_PASSWORD"),
				Hostname: getRequiredValue("DATABASE_HOSTNAME"),
				Port:     getValueToInt("DATABASE_PORT"),
				Name:     getRequiredValue("DATABASE_NAME"),
			},
			Backend: backendConfig{
				Port: getValueToInt("BACKEND_PORT"),
			},
		}

	})

	return instance
}
