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
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/rs/zerolog/log"
)

type ContentRepo interface {
	Insert(ctx context.Context, data []esmodel.Content) error
	GetFootnotesByWork(ctx context.Context, workCode string, ordinals []int32) ([]esmodel.Content, error)
	GetHeadingsByWork(ctx context.Context, workCode string, ordinals []int32) ([]esmodel.Content, error)
	GetParagraphsByWork(ctx context.Context, workCode string, ordinals []int32) ([]esmodel.Content, error)
	GetSummariesByWork(ctx context.Context, workCode string, ordinals []int32) ([]esmodel.Content, error)
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
		Mappings: esmodel.ContentMapping,
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
				"german_stemming_analyzer": &types.CustomAnalyzer{
					Tokenizer: "standard",
					Filter: []string{
						"lowercase",            // only case-insensitive
						"german_normalization", // normalize Umlauts and ß
						"german_stemmer",       // see below
					},
				},
			},
			Filter: map[string]types.TokenFilter{
				"german_stemmer": &types.StemmerTokenFilter{
					Type:     "stemmer",
					Language: commonutil.StrPtr("german"),
				},
			},
		},
	}
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

func (rec *contentRepoImpl) GetFootnotesByWork(ctx context.Context, workCode string, ordinals []int32) ([]esmodel.Content, error) {
	contentQuery := util.CreateContentQuery(
		workCode,
		[]esmodel.Type{esmodel.Footnote},
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

func (rec *contentRepoImpl) GetHeadingsByWork(ctx context.Context, workCode string, ordinals []int32) ([]esmodel.Content, error) {
	contentQuery := util.CreateContentQuery(
		workCode,
		[]esmodel.Type{esmodel.Heading},
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

func (rec *contentRepoImpl) GetParagraphsByWork(ctx context.Context, workCode string, ordinals []int32) ([]esmodel.Content, error) {
	contentQuery := util.CreateContentQuery(
		workCode,
		[]esmodel.Type{esmodel.Paragraph},
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

func (rec *contentRepoImpl) GetSummariesByWork(ctx context.Context, workCode string, ordinals []int32) ([]esmodel.Content, error) {
	contentQuery := util.CreateContentQuery(
		workCode,
		[]esmodel.Type{esmodel.Summary},
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

	// TODO check if sorting make this so slow
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
			Size:      commonutil.IntPtr(10000), // TODO check: what to do when there are more results; and how to show them in the frontend
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
			Text:     c.SearchText,
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
