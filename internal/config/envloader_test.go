package config

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLoadEnv_Success(t *testing.T) {
	fileName := ".env.test"
	err := os.WriteFile(fileName, []byte("KEY=VALUE"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	defer os.Remove(fileName)

	err = LoadEnv(fileName)
	assert.NoError(t, err)

	assert.Equal(t, "VALUE", os.Getenv("KEY"))
}

func TestLoadEnv_FileNotFound(t *testing.T) {
	err := LoadEnv(".env.notfound")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Error loading .env file")
}

func TestLoadEnv_InvalidFormat(t *testing.T) {
	fileName := ".env.invalid"
	err := os.WriteFile(fileName, []byte("KEY=VALUE\nINVALID_ENTRY"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	defer os.Remove(fileName)

	err = LoadEnv(fileName)
	assert.NoError(t, err)
}
