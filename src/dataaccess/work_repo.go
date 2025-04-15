package dataaccess

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/result"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/util"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/work_repo_mock.go -package=mocks

type WorkRepo interface {
	Insert(ctx context.Context, data *esmodel.Work) error
	Update(ctx context.Context, data *esmodel.Work) error
	Get(ctx context.Context, id string) (*esmodel.Work, error)
	Delete(ctx context.Context, id string) error
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

func (rec *workRepoImpl) Insert(ctx context.Context, data *esmodel.Work) error {
	res, err := rec.dbClient.Index(rec.indexName).Document(&data).Do(ctx)
	if err != nil {
		return err
	}
	if res.Result != result.Created {
		return fmt.Errorf("unable to create work with id %s", data.Id)
	}
	data.Id = res.Id_
	return nil
}

func (rec *workRepoImpl) Update(ctx context.Context, data *esmodel.Work) error {
	res, err := rec.dbClient.Update(rec.indexName, data.Id).Doc(&data).Do(ctx)
	if err != nil {
		return err
	}
	if res.Result != result.Updated {
		return fmt.Errorf("unable to update work with id %s", data.Id)
	}
	return nil
}

func (rec *workRepoImpl) Get(ctx context.Context, id string) (*esmodel.Work, error) {
	res, err := rec.dbClient.Get(rec.indexName, id).Do(context.TODO())
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, fmt.Errorf("work with ID %s not found", id)
	}

	var work esmodel.Work
	err = json.Unmarshal(res.Source_, &work)
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func (rec *workRepoImpl) Delete(ctx context.Context, id string) error {
	res, err := rec.dbClient.Delete(rec.indexName, id).Do(ctx)
	if err != nil {
		return err
	}
	if res.Result != result.Deleted {
		return fmt.Errorf("unable to delete work with id %s", id)
	}
	return nil
}
