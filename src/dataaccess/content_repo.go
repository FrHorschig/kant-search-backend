package dataaccess

//go:generate mockgen -source=$GOFILE -destination=mocks/content_repo_mock.go -package=mocks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/deletebyquery"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/operationtype"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/util"
	"github.com/rs/zerolog/log"
)

type ContentRepo interface {
	Insert(ctx context.Context, data []esmodel.Content) error
	GetFootnotesByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error)
	GetHeadingsByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error)
	GetParagraphsByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error)
	GetSummariesByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error)
	DeleteByWorkId(ctx context.Context, workId string) error
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

func (rec *contentRepoImpl) Insert(ctx context.Context, data []esmodel.Content) error {
	bulk := rec.dbClient.Bulk().Index(rec.indexName)
	for _, c := range data {
		bulk.CreateOp(*types.NewCreateOperation(), c)
	}
	res, err := bulk.Do(context.TODO())
	if err != nil {
		return err
	}

	for i, item := range res.Items {
		e := item[operationtype.Create].Error
		if e != nil && e.Reason != nil {
			return errors.New(*e.Reason)
		}
		if *item[operationtype.Create].Result != "created" {
			return errors.New("unable to create new document")
		}
		data[i].Id = *item[operationtype.Create].Id_
	}
	return nil
}

func (rec *contentRepoImpl) GetFootnotesByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: createContentQuery(workId, string(esmodel.Footnote)),
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	contents := []esmodel.Content{}
	for _, hit := range res.Hits.Hits {
		var c esmodel.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (rec *contentRepoImpl) GetHeadingsByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: createContentQuery(workId, string(esmodel.Heading)),
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	contents := []esmodel.Content{}
	for _, hit := range res.Hits.Hits {
		var c esmodel.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (rec *contentRepoImpl) GetParagraphsByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: createContentQuery(workId, string(esmodel.Paragraph)),
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	contents := []esmodel.Content{}
	for _, hit := range res.Hits.Hits {
		var c esmodel.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (rec *contentRepoImpl) GetSummariesByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: createContentQuery(workId, string(esmodel.Summary)),
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	contents := []esmodel.Content{}
	for _, hit := range res.Hits.Hits {
		var c esmodel.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (rec *contentRepoImpl) DeleteByWorkId(ctx context.Context, workId string) error {
	res, err := rec.dbClient.DeleteByQuery(rec.indexName).Request(&deletebyquery.Request{
		Query: createTermQuery("workId", workId),
	}).Do(ctx)
	if err != nil {
		return err
	}

	if len(res.Failures) > 0 {
		for _, fail := range res.Failures {
			e := fail.Cause.Reason
			if e != nil {
				log.Error().Msgf("Failed to delete content: %s", *e)
			}
		}
		return fmt.Errorf("unable to delete work with id %s", workId)
	}
	return nil
}

func createContentQuery(workId string, cType string) *types.Query {
	return &types.Query{
		Bool: &types.BoolQuery{
			Must: []types.Query{
				{
					Term: map[string]types.TermQuery{
						"workId": {Value: workId},
					},
				},
				{
					Term: map[string]types.TermQuery{
						"type": {Value: cType},
					},
				},
			},
		},
	}
}
