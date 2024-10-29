package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvs() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func GetEnvOrDie(key string) string {
	value := os.Getenv(key)

	if value == "" {
		err := fmt.Errorf("missing environment variable %s", key)
		panic(err)
	}

	return value
}
