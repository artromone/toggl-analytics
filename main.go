package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type UserCredentials struct {
	APIKey      string
	WorkspaceID string
	FileName    string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	user1 := UserCredentials{
		APIKey:      os.Getenv("U1_API_KEY"),
		WorkspaceID: os.Getenv("U1_WORKSPACE_ID"),
		FileName:    os.Getenv("U1_FILE_NAME"),
	}

	fmt.Printf("User 1: %+v\n", user1)
}
