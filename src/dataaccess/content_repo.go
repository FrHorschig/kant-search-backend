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
	GetFootnotesByWorkCode(ctx context.Context, workCode string) ([]esmodel.Content, error)
	GetHeadingsByWorkCode(ctx context.Context, workCode string) ([]esmodel.Content, error)
	GetParagraphsByWorkCode(ctx context.Context, workCode string) ([]esmodel.Content, error)
	GetSummariesByWorkCode(ctx context.Context, workCode string) ([]esmodel.Content, error)
	DeleteByWorkCode(ctx context.Context, workCode string) error
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

	for _, item := range res.Items {
		e := item[operationtype.Create].Error
		if e != nil && e.Reason != nil {
			return errors.New(*e.Reason)
		}
		if *item[operationtype.Create].Result != "created" {
			return errors.New("unable to create new document")
		}
	}
	return nil
}

func (rec *contentRepoImpl) GetFootnotesByWorkCode(ctx context.Context, workCode string) ([]esmodel.Content, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: util.CreateContentQuery(workCode, []esmodel.Type{esmodel.Footnote}),
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

func (rec *contentRepoImpl) GetHeadingsByWorkCode(ctx context.Context, workCode string) ([]esmodel.Content, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: util.CreateContentQuery(workCode, []esmodel.Type{esmodel.Heading}),
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

func (rec *contentRepoImpl) GetParagraphsByWorkCode(ctx context.Context, workCode string) ([]esmodel.Content, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: util.CreateContentQuery(workCode, []esmodel.Type{esmodel.Paragraph}),
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

func (rec *contentRepoImpl) GetSummariesByWorkCode(ctx context.Context, workCode string) ([]esmodel.Content, error) {
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: util.CreateContentQuery(workCode, []esmodel.Type{esmodel.Summary}),
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

func (rec *contentRepoImpl) DeleteByWorkCode(ctx context.Context, workCode string) error {
	res, err := rec.dbClient.DeleteByQuery(rec.indexName).Request(&deletebyquery.Request{
		Query: createTermQuery("workCode", workCode),
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
		return fmt.Errorf("unable to delete work with code %s", workCode)
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
	optionQueries := util.CreateOptionQueries(options) // BUG HERE

	res, err := rec.dbClient.Search().Index(rec.indexName).Request(
		&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Must:   []types.Query{*searchQuery},
					Filter: optionQueries, // BUG HERE
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
			Snippets: createSnippets(hit.Highlight["searchText"]),
			Pages:    c.Pages,
			Ordinal:  c.Ordinal,
			WorkCode: c.WorkCode,
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
