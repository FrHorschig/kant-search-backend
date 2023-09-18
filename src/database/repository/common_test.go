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
	dbUrl := createDbContainer()
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}
	testDb = db
	code := m.Run()
	os.Exit(code)
}

func createDbContainer() string {
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

	connStr, err := cont.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		panic(err)
	}
	return connStr
}
