package postgres

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	const dbConnEnv = "DATABASE_CONNECTION"
	ddlConnStr := os.Getenv(dbConnEnv)
	if ddlConnStr == "" {
		log.Printf("%s is not set, skipping", dbConnEnv)
		return
	}

	var teardown func()
	_testStorage, teardown = NewTestStorage(ddlConnStr, filepath.Join("..", "..", "migrations", "sql"))

	exitCode := m.Run()

	if teardown != nil {
		teardown()
	}

	os.Exit(exitCode)
}
