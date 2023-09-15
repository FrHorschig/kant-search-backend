package repository

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

var testDb *sql.DB

func TestMain(m *testing.M) {
	// pool, resource := setupDb()
	// testDb = util.InitDbConnection()
	// code := m.Run()
	// cleanupDb(pool, resource)
	// os.Exit(code)
}

func setupDb() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "kant-search-database",
		Env: []string{
			"DB_USER=postgres",
			"DB_PASSWORD=postgres",
			"DB_NAME=testdb",
			"DB_HOST=localhost",
			"DB_PORT=5432",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	databaseUrl := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/testdb?sslmode=disable", resource.GetPort("5432/tcp"))
	log.Println("Connecting to database on url: ", databaseUrl)
	resource.Expire(20) // kill container after timeout in case cleanup fails
	pool.MaxWait = 20 * time.Second
	if err = pool.Retry(func() error {
		db, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		testDb = db
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return pool, resource
}

func cleanupDb(pool *dockertest.Pool, resource *dockertest.Resource) {
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
