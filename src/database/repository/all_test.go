package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
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
	ctx := context.Background()
	envVarValue := "kantsearch"

	req := testcontainers.ContainerRequest{
		Env: map[string]string{
			"POSTGRES_USER":     envVarValue,
			"POSTGRES_PASSWORD": envVarValue,
			"POSTGRES_DB":       envVarValue,
		},
		ExposedPorts: []string{"5432/tcp"},
		Image:        "postgres:14.3",
		WaitingFor: wait.ForExec([]string{"pg_isready"}).
			WithPollInterval(2 * time.Second).
			WithExitCodeMatcher(func(exitCode int) bool {
				return exitCode == 0
			}),
	}
	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		panic(err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", envVarValue, envVarValue, host, mappedPort, envVarValue)
}
