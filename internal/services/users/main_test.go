package users

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := os.Setenv("ENV", "test"); err != nil {
		log.Fatalf("failed to set env file: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}
