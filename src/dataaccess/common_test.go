//go:build integration
// +build integration

package dataaccess

import (
	"context"
	"os"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	estest "github.com/testcontainers/testcontainers-go/modules/elasticsearch"
)

var dbClient *elasticsearch.TypedClient

func TestMain(m *testing.M) {
	container := createEsContainer()
	code := m.Run()
	cleanupDbContainer(container)
	os.Exit(code)
}

func createEsContainer() *estest.ElasticsearchContainer {
	ctx := context.Background()
	cont, err := estest.Run(ctx, "docker.elastic.co/elasticsearch/elasticsearch:8.17.4")
	if err != nil {
		panic(err)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			cont.Settings.Address,
		},
		Username: "elastic",
		Password: cont.Settings.Password,
		CACert:   cont.Settings.CACert,
	}

	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		panic(err)
	}
	dbClient = client
	return cont
}

func cleanupDbContainer(container *estest.ElasticsearchContainer) {
	if err := container.Terminate(context.Background()); err != nil {
		panic(err)
	}
}
