package util

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func CreateIndex(es *elasticsearch.TypedClient, name string, mapping *types.TypeMapping) error {
	_, err := es.Indices.Create(name).Request(&create.Request{
		Mappings: mapping,
	}).Do(context.Background())
	// TODO check what happens if the index already exists
	return err
}
