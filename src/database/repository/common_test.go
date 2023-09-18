package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDb *sql.DB

func TestMain(m *testing.M) {
	container := createDbContainer()
	openDb(container)
	code := m.Run()
	cleanupDbContainer(container)
	os.Exit(code)
}

func createDbContainer() *postgres.PostgresContainer {
	envVarValue := "kantsearch"
	cont, err := postgres.RunContainer(context.Background(),
		testcontainers.WithImage("kant-search-database"),
		postgres.WithUsername(envVarValue),
		postgres.WithPassword(envVarValue),
		postgres.WithDatabase(envVarValue),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(20*time.Second)),
	)
	if err != nil {
		panic(err)
	}
	return cont
}

func openDb(container *postgres.PostgresContainer) {
	connStr, err := container.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	testDb = db
}

func cleanupDbContainer(container *postgres.PostgresContainer) {
	testDb.Close()
	if err := container.Terminate(context.Background()); err != nil {
		panic(err)
	}
}
