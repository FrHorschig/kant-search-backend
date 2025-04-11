package dataaccess

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/util"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/work_repo_mock.go -package=mocks

type WorkRepo interface {
	Insert(ctx context.Context, data esmodel.Work) error
	Update(ctx context.Context, data esmodel.Work) error
	Delete(ctx context.Context, id int32) error
}

type workRepoImpl struct {
	dbClient  *elasticsearch.TypedClient
	indexName string
}

func NewWorkRepo(dbClient *elasticsearch.TypedClient) WorkRepo {
	repo := &workRepoImpl{
		dbClient:  dbClient,
		indexName: "works",
	}
	err := util.CreateIndex(repo.dbClient, repo.indexName, esmodel.WorkMapping)
	if err != nil {
		panic(err)
	}
	return repo
}

func (rec *workRepoImpl) Insert(ctx context.Context, data esmodel.Work) error {
	// TODO implement me
	return nil
}

func (rec *workRepoImpl) Update(ctx context.Context, data esmodel.Work) error {
	// TODO implement me
	return nil
}

func (rec *workRepoImpl) Delete(ctx context.Context, id int32) error {
	// TODO implement me
	return nil
}
