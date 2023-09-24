//go:build integration
// +build integration

package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDb *sql.DB

func (repo *paragraphRepoImpl) insertParagraphs(workId int32, text string) (int32, error) {
	ctx := context.Background()
	p := model.Paragraph{Text: text, Pages: []int32{1, 2}, WorkId: workId}
	return repo.Insert(ctx, p)
}

func (repo *sentenceRepoImpl) insertSentences(paragraphId int32, texts []string) ([]int32, error) {
	ctx := context.Background()
	var sentences []model.Sentence
	for _, text := range texts {
		sentences = append(sentences, model.Sentence{Text: text, ParagraphId: paragraphId})
	}
	return repo.Insert(ctx, sentences)
}

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
