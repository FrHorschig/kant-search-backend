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
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/operationtype"
	commonutil "github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/rs/zerolog/log"
)

type ContentRepo interface {
	Insert(ctx context.Context, data []model.Content) error
	GetFootnotesByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	GetHeadingsByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	GetParagraphsByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	GetSummariesByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	DeleteByWork(ctx context.Context, workCode string) error
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
	err := createContentIndex(repo.dbClient, repo.indexName)
	if err != nil {
		panic(err)
	}
	return repo
}

func createContentIndex(es *elasticsearch.TypedClient, name string) error {
	ctx := context.Background()
	ok, err := es.Indices.Exists(name).Do(ctx)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	res, err := es.Indices.Create(name).Request(&create.Request{
		Mappings: model.ContentMapping,
		Settings: buildSettings(),
	}).Do(ctx)
	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return fmt.Errorf("creation of index '%s' not acknowledged", name)
	}
	return err
}

func buildSettings() *types.IndexSettings {
	return &types.IndexSettings{
		Analysis: &types.IndexSettingsAnalysis{
			Analyzer: map[string]types.Analyzer{
				string(model.NoStemming): &types.CustomAnalyzer{
					Tokenizer: "standard",
					Filter:    []string{"lowercase"},
				},
				string(model.GermanStemming): &types.CustomAnalyzer{
					Tokenizer: "standard",
					Filter:    []string{"lowercase", string(model.GermanStemming)},
				},
			},
			Filter: map[string]types.TokenFilter{
				string(model.GermanStemming): &types.StemmerTokenFilter{
					Type:     "stemmer",
					Language: commonutil.StrPtr("german"),
				},
			},
		},
	}
}

func (rec *contentRepoImpl) Insert(ctx context.Context, data []model.Content) error {
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

func (rec *contentRepoImpl) GetFootnotesByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error) {
	contentQuery := util.CreateContentQuery(
		workCode,
		[]model.Type{model.Footnote},
	)
	if len(ordinals) > 0 {
		contentQuery.Bool.Filter = append(
			contentQuery.Bool.Filter,
			util.CreateOrdinalQuery(ordinals),
		)
	}
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: contentQuery,
		Size:  commonutil.IntPtr(resultsSize),
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	contents := []model.Content{}
	for _, hit := range res.Hits.Hits {
		var c model.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (rec *contentRepoImpl) GetHeadingsByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error) {
	contentQuery := util.CreateContentQuery(
		workCode,
		[]model.Type{model.Heading},
	)
	if len(ordinals) > 0 {
		contentQuery.Bool.Filter = append(
			contentQuery.Bool.Filter,
			util.CreateOrdinalQuery(ordinals),
		)
	}
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: contentQuery,
		Size:  commonutil.IntPtr(resultsSize),
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	contents := []model.Content{}
	for _, hit := range res.Hits.Hits {
		var c model.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (rec *contentRepoImpl) GetParagraphsByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error) {
	contentQuery := util.CreateContentQuery(
		workCode,
		[]model.Type{model.Paragraph},
	)
	if len(ordinals) > 0 {
		contentQuery.Bool.Filter = append(
			contentQuery.Bool.Filter,
			util.CreateOrdinalQuery(ordinals),
		)
	}
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: contentQuery,
		Size:  commonutil.IntPtr(resultsSize),
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	contents := []model.Content{}
	for _, hit := range res.Hits.Hits {
		var c model.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (rec *contentRepoImpl) GetSummariesByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error) {
	contentQuery := util.CreateContentQuery(
		workCode,
		[]model.Type{model.Summary},
	)
	if len(ordinals) > 0 {
		contentQuery.Bool.Filter = append(
			contentQuery.Bool.Filter,
			util.CreateOrdinalQuery(ordinals),
		)
	}
	res, err := rec.dbClient.Search().Request(&search.Request{
		Query: contentQuery,
		Size:  commonutil.IntPtr(resultsSize),
	}).Do(ctx)
	if err != nil {
		return nil, err
	}

	contents := []model.Content{}
	for _, hit := range res.Hits.Hits {
		var c model.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, nil
}

func (rec *contentRepoImpl) DeleteByWork(ctx context.Context, workCode string) error {
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

	_, err = rec.dbClient.Indices.Refresh().Index(rec.indexName).Do(ctx)
	return err
}

func (rec *contentRepoImpl) Search(ctx context.Context, ast *model.AstNode, options model.SearchOptions) ([]model.SearchResult, error) {
	analyzer := model.NoStemming
	if options.WithStemming {
		analyzer = model.GermanStemming
	}
	searchQuery, err := util.CreateSearchQuery(ast, analyzer)
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
			Highlight: util.CreateHighlightOptions(analyzer),
			Size:      commonutil.IntPtr(10000),
		}).Do(ctx)
	if err != nil {
		return nil, err
	}

	results := []model.SearchResult{}
	for _, hit := range res.Hits.Hits {
		var c model.Content
		err = json.Unmarshal(hit.Source_, &c)
		if err != nil {
			return nil, err
		}
		results = append(results, model.SearchResult{
			HighlightText: hit.Highlight["searchText."+string(analyzer)][0],
			FmtText:       c.FmtText,
			Pages:         c.Pages,
			PageByIndex:   c.PageByIndex,
			LineByIndex:   c.LineByIndex,
			Ordinal:       c.Ordinal,
			WorkCode:      c.WorkCode,
			WordIndexMap:  c.WordIndexMap,
		})
	}
	return results, nil
}
