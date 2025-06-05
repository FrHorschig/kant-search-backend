package dataaccess

//go:generate mockgen -source=$GOFILE -destination=mocks/work_repo_mock.go -package=mocks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/result"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/util"
)

type WorkRepo interface {
	Insert(ctx context.Context, data *esmodel.Work) error
	Get(ctx context.Context, code string) (*esmodel.Work, error)
	Delete(ctx context.Context, code string) error
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
		return fmt.Errorf("unable to create work with code %s", data.Code)
	}
	return nil
}

func (rec *workRepoImpl) Get(ctx context.Context, code string) (*esmodel.Work, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: createCodeQuery(code),
	}).Do(ctx)
	if err != nil {
		return nil, err
	}
	if len(res.Hits.Hits) == 0 {
		return nil, fmt.Errorf("no work with code %s found", code)
	}
	if len(res.Hits.Hits) > 1 {
		return nil, fmt.Errorf("multiple works with code %s found", code)
	}

	var work esmodel.Work
	err = json.Unmarshal(res.Hits.Hits[0].Source_, &work)
	if err != nil {
		return nil, err
	}
	return &work, nil
}

func (rec *workRepoImpl) Delete(ctx context.Context, code string) error {
	res, err := rec.dbClient.DeleteByQuery(rec.indexName).Query(createCodeQuery(code)).Do(ctx)
	if err != nil {
		return err
	}
	if len(res.Failures) > 0 {
		return fmt.Errorf("unable to delete work with code %s", code)
	}
	return nil
}

func createCodeQuery(code string) *types.Query {
	return &types.Query{
		Term: map[string]types.TermQuery{"code": {Value: code}},
	}
}
