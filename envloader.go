package main

import (
	"fmt"
	"github.com/joho/godotenv"
)

func LoadEnv(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}
