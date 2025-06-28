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
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/rs/zerolog/log"
)

const analyzerPrefix = "searchText."

type ContentRepo interface {
	Insert(ctx context.Context, data []model.Content) error
	GetFootnotesByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	GetHeadingsByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	GetParagraphsByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	GetSummariesByWork(ctx context.Context, workCode string, ordinals []int32) ([]model.Content, error)
	DeleteByWork(ctx context.Context, workCode string) error
	Search(ctx context.Context, ast *model.SearchTermNode, options model.SearchOptions) ([]model.SearchResult, error)
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
					Language: util.StrPtr("german"),
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
	contentQuery := createContentQuery(
		workCode,
		[]model.Type{model.Footnote},
	)
	if len(ordinals) > 0 {
		contentQuery.Bool.Filter = append(
			contentQuery.Bool.Filter,
			createOrdinalQuery(ordinals),
		)
	}
	res, err := rec.dbClient.Search().
		AllowPartialSearchResults(false).
		Request(&search.Request{
			Query: contentQuery,
			Size:  util.IntPtr(resultsSize),
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
	contentQuery := createContentQuery(
		workCode,
		[]model.Type{model.Heading},
	)
	if len(ordinals) > 0 {
		contentQuery.Bool.Filter = append(
			contentQuery.Bool.Filter,
			createOrdinalQuery(ordinals),
		)
	}
	res, err := rec.dbClient.Search().
		AllowPartialSearchResults(false).
		Request(&search.Request{
			Query: contentQuery,
			Size:  util.IntPtr(resultsSize),
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
	contentQuery := createContentQuery(
		workCode,
		[]model.Type{model.Paragraph},
	)
	if len(ordinals) > 0 {
		contentQuery.Bool.Filter = append(
			contentQuery.Bool.Filter,
			createOrdinalQuery(ordinals),
		)
	}
	res, err := rec.dbClient.Search().
		AllowPartialSearchResults(false).
		Request(&search.Request{
			Query: contentQuery,
			Size:  util.IntPtr(resultsSize),
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
	contentQuery := createContentQuery(
		workCode,
		[]model.Type{model.Summary},
	)
	if len(ordinals) > 0 {
		contentQuery.Bool.Filter = append(
			contentQuery.Bool.Filter,
			createOrdinalQuery(ordinals),
		)
	}
	res, err := rec.dbClient.Search().
		AllowPartialSearchResults(false).
		Request(&search.Request{
			Query: contentQuery,
			Size:  util.IntPtr(resultsSize),
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
	res, err := rec.dbClient.DeleteByQuery(rec.indexName).
		Request(&deletebyquery.Request{
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

func (rec *contentRepoImpl) Search(ctx context.Context, ast *model.SearchTermNode, options model.SearchOptions) ([]model.SearchResult, error) {
	analyzer := model.NoStemming
	if options.WithStemming {
		analyzer = model.GermanStemming
	}
	searchQuery, err := createSearchQuery(ast, analyzer)
	if err != nil {
		return nil, err
	}
	if searchQuery == nil {
		// empty search term (== nil searchQueries) is catched in api layer, so if this is the case, the error is technical, not a user error
		return nil, errors.New("search AST must not be nil")
	}
	optionQueries := createOptionQueries(options)

	res, err := rec.dbClient.Search().Index(rec.indexName).
		AllowPartialSearchResults(false).
		Request(
			&search.Request{
				Query: &types.Query{
					Bool: &types.BoolQuery{
						Must:   []types.Query{*searchQuery},
						Filter: optionQueries,
					},
				},
				Sort:      createSortOptions(),
				Highlight: createHighlightOptions(analyzer),
				Size:      util.IntPtr(10000),
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

func createWorkCodeQuery(workCode string) types.Query {
	return types.Query{
		Term: map[string]types.TermQuery{
			"workCode": {Value: workCode},
		},
	}
}

func createContentQuery(workCode string, cType []model.Type) *types.Query {
	return &types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				createWorkCodeQuery(workCode),
				createTypeQuery(cType),
			},
		},
	}
}

func createSearchQuery(node *model.SearchTermNode, analyzer model.Analyzer) (*types.Query, error) {
	if node == nil {
		return nil, nil
	}
	if node.Token.IsAnd {
		return createAndQuery(node, analyzer)
	}
	if node.Token.IsOr {
		return createOrQuery(node, analyzer)
	}
	if node.Token.IsNot {
		return createNotQuery(node, analyzer)
	}
	if node.Token.IsWord {
		return createTextMatchQuery(node.Token.Text, analyzer), nil
	}
	if node.Token.IsPhrase {
		return createPhraseQuery(node.Token.Text, analyzer), nil
	}
	return nil, errors.New("invalid token type")
}

func createSortOptions() []types.SortCombinations {
	return []types.SortCombinations{
		types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"ordinal": {Order: &sortorder.Asc},
			},
		},
	}
}

func createHighlightOptions(analyzer model.Analyzer) *types.Highlight {
	return &types.Highlight{
		Fields: map[string]types.HighlightField{
			analyzerPrefix + string(analyzer): {
				NumberOfFragments: util.IntPtr(0),
			},
		},
		PreTags:  []string{"<ks-meta-hit>"},
		PostTags: []string{"</ks-meta-hit>"},
	}
}

func createTypeQuery(cType []model.Type) types.Query {
	return types.Query{Terms: &types.TermsQuery{
		TermsQuery: map[string]types.TermsQueryField{
			"type": cType,
		},
	}}
}

func createOrdinalQuery(ordinals []int32) types.Query {
	values := make([]interface{}, len(ordinals))
	for i, v := range ordinals {
		values[i] = v
	}
	return types.Query{
		Terms: &types.TermsQuery{
			TermsQuery: map[string]types.TermsQueryField{
				"ordinal": values,
			},
		},
	}
}

func createAndQuery(node *model.SearchTermNode, analyzer model.Analyzer) (*types.Query, error) {
	q1, err := createSearchQuery(node.Left, analyzer)
	if err != nil {
		return nil, err
	}
	q2, err := createSearchQuery(node.Right, analyzer)
	if err != nil {
		return nil, err
	}
	if q1 == nil || q2 == nil {
		return nil, errors.New("AND nodes must have both a left and a right child")
	}
	return &types.Query{Bool: &types.BoolQuery{
		Must: []types.Query{*q1, *q2},
	}}, nil
}

func createOrQuery(node *model.SearchTermNode, analyzer model.Analyzer) (*types.Query, error) {
	q1, err := createSearchQuery(node.Left, analyzer)
	if err != nil {
		return nil, err
	}
	q2, err := createSearchQuery(node.Right, analyzer)
	if err != nil {
		return nil, err
	}
	if q1 == nil || q2 == nil {
		return nil, errors.New("OR nodes must have both a left and a right child")
	}
	return &types.Query{Bool: &types.BoolQuery{
		Should: []types.Query{*q1, *q2},
	}}, nil
}

func createNotQuery(node *model.SearchTermNode, analyzer model.Analyzer) (*types.Query, error) {
	q1, err := createSearchQuery(node.Left, analyzer)
	if err != nil {
		return nil, err
	}
	if q1 == nil {
		q2, err := createSearchQuery(node.Right, analyzer)
		if err != nil {
			return nil, err
		}
		if q2 == nil {
			return nil, errors.New("NOT nodes must have either a left and a right child")
		}
		return &types.Query{Bool: &types.BoolQuery{
			MustNot: []types.Query{*q2},
		}}, nil
	}
	return &types.Query{Bool: &types.BoolQuery{
		MustNot: []types.Query{*q1},
	}}, nil
}

func createPhraseQuery(phrase string, analyzer model.Analyzer) *types.Query {
	return &types.Query{
		MatchPhrase: map[string]types.MatchPhraseQuery{
			analyzerPrefix + string(analyzer): {Query: phrase},
		},
	}
}

func createTextMatchQuery(term string, analyzer model.Analyzer) *types.Query {
	return &types.Query{
		Match: map[string]types.MatchQuery{
			analyzerPrefix + string(analyzer): {Query: term},
		},
	}
}

func createOptionQueries(opts model.SearchOptions) []types.Query {
	tps := []model.Type{model.Paragraph}
	if opts.IncludeHeadings {
		tps = append(tps, model.Heading)
	}
	if opts.IncludeFootnotes {
		tps = append(tps, model.Footnote)
	}
	if opts.IncludeSummaries {
		tps = append(tps, model.Summary)
	}
	return []types.Query{
		createWorkCodesQuery(opts.WorkCodes),
		createTypeQuery(tps),
	}
}

func createWorkCodesQuery(workCodes []string) types.Query {
	return types.Query{Terms: &types.TermsQuery{
		TermsQuery: map[string]types.TermsQueryField{
			"workCode": workCodes,
		},
	}}
}
