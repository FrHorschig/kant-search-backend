package dataaccess

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/util"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/content_repo_mock.go -package=mocks

type ContentRepo interface {
	Insert(ctx context.Context, data esmodel.Content) error
	Delete(ctx context.Context, id int32) error
}

type contentRepoImpl struct {
	dbClient  *elasticsearch.TypedClient
	indexName string
}

func NewContentRepo(dbClient *elasticsearch.TypedClient) ContentRepo {
	repo := &contentRepoImpl{
		dbClient:  dbClient,
		indexName: "contents",
	}
	err := util.CreateIndex(repo.dbClient, repo.indexName, esmodel.ContentMapping)
	if err != nil {
		panic(err)
	}
	return repo
}

func (rec *contentRepoImpl) Insert(ctx context.Context, data esmodel.Content) error {
	// TODO implement me
	return nil
}

func (rec *contentRepoImpl) Delete(ctx context.Context, id int32) error {
	// TODO implement me
	return nil
}
