package main

import (
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file: %v", err)
	}

	users := GetAllUserCredentials()

	for userPrefix, credentials := range users {
		fmt.Printf("Credentials for %s: %+v\n", userPrefix, credentials)
	}
}
