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
	commonutil "github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/rs/zerolog/log"
)

type ContentRepo interface {
	Insert(ctx context.Context, data []esmodel.Content) error
	GetFootnotesByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error)
	GetHeadingsByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error)
	GetParagraphsByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error)
	GetSummariesByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error)
	DeleteByWorkId(ctx context.Context, workId string) error
	Search(ctx context.Context, ast *model.AstNode, options model.SearchOptions) ([]model.SearchResult, error)
}

const resultsSize = 10000

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
	insert := rec.dbClient.Bulk().Index(rec.indexName)
	for _, c := range data {
		insert.CreateOp(*types.NewCreateOperation(), c)
	}
	res, err := insert.Do(ctx)
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

	update := rec.dbClient.Bulk().Index(rec.indexName)
	for _, c := range data {
		update.UpdateOp(types.UpdateOperation{Id_: &c.Id}, c, types.NewUpdateAction())
	}
	_, err = update.Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (rec *contentRepoImpl) GetFootnotesByWorkId(ctx context.Context, workId string) ([]esmodel.Content, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: util.CreateContentQuery(workId, []esmodel.Type{esmodel.Footnote}),
		Size:  commonutil.IntPtr(resultsSize),
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
		Query: util.CreateContentQuery(workId, []esmodel.Type{esmodel.Heading}),
		Size:  commonutil.IntPtr(resultsSize),
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
		Query: util.CreateContentQuery(workId, []esmodel.Type{esmodel.Paragraph}),
		Size:  commonutil.IntPtr(resultsSize),
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
		Query: util.CreateContentQuery(workId, []esmodel.Type{esmodel.Summary}),
		Size:  commonutil.IntPtr(resultsSize),
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

func (rec *contentRepoImpl) Search(ctx context.Context, ast *model.AstNode, options model.SearchOptions) ([]model.SearchResult, error) {
	searchQuery, err := util.CreateSearchQuery(ast)
	if err != nil {
		return nil, err
	}
	if searchQuery == nil {
		// empty search term (== nil searchQueries) is catched in api layer, so if this is the case, the error is technical, not a user error
		return nil, errors.New("search AST must not be nil")
	}
	optionQueries := util.CreateOptionQueries(options)

	res, err := rec.dbClient.Search().Index(rec.indexName).Request(
		&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Must:   []types.Query{*searchQuery},
					Filter: optionQueries,
				},
			},
			Sort:      util.CreateSortOptions(),
			Highlight: util.CreateHighlightOptions(),
		}).Do(ctx)
	if err != nil {
		return nil, err
	}

	results := []model.SearchResult{}
	for _, hit := range res.Hits.Hits {
		var c esmodel.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		// TODO later: set numbers of fragments to zero, then extract the indices of the pre and post highlight tag, then map them to the fmtText field, use it with the indices for a) showing the full paragraph with highlighting and b) finding the exact page and line number of a match
		results = append(results, model.SearchResult{
			Snippets:  createSnippets(hit.Highlight["searchText"]),
			Pages:     c.Pages,
			ContentId: c.Id,
			WorkId:    c.WorkId,
		})
	}
	return results, nil
}

func createSnippets(snips []string) []string {
	result := []string{}
	for _, s := range snips {
		result = append(result, fmt.Sprintf("...%s...", s))
	}
	return result
}
