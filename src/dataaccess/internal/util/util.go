package util

import (
	"context"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func CreateIndex(es *elasticsearch.TypedClient, name string, mapping *types.TypeMapping) error {
	ctx := context.Background()
	ok, err := es.Indices.Exists(name).Do(ctx)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	res, err := es.Indices.Create(name).Request(&create.Request{
		Mappings: mapping,
	}).Do(ctx)
	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return fmt.Errorf("creation of index '%s' not acknowledged", name)
	}
	return err
}

func createWorkIdQuery(workId string) types.Query {
	return types.Query{
		Term: map[string]types.TermQuery{
			"workId": {Value: workId},
		},
	}
}

func CreateContentQuery(workId string, cType []esmodel.Type) *types.Query {
	return &types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				createWorkIdQuery(workId),
				createTypeQuery(cType),
			},
		},
	}
}

func CreateSearchQuery(node *model.AstNode) (*types.Query, error) {
	if node == nil {
		return nil, nil
	}
	if node.Token.IsAnd {
		return createAndQuery(node)
	}
	if node.Token.IsOr {
		return createOrQuery(node)
	}
	if node.Token.IsNot {
		return createNotQuery(node)
	}
	if node.Token.IsWord {
		return createTextMatchQuery(node.Token.Text), nil
	}
	if node.Token.IsPhrase {
		return createPhraseQuery(node.Token.Text), nil
	}
	return nil, errors.New("invalid token type")
}

func CreateSortOptions() []types.SortCombinations {
	return []types.SortCombinations{
		types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"ordinal": {Order: &sortorder.Asc},
			},
		},
	}
}

func CreateHighlightOptions() *types.Highlight {
	return &types.Highlight{
		Fields: map[string]types.HighlightField{
			"searchText": {
				// TODO adjust this after testing with the frontend
				FragmentSize:      util.IntPtr(150),
				NumberOfFragments: util.IntPtr(3),
			},
		},
		PreTags:  []string{"<ks-meta-hit>"},
		PostTags: []string{"</ks-meta-hit>"},
	}
}

func createTypeQuery(cType []esmodel.Type) types.Query {
	return types.Query{Terms: &types.TermsQuery{
		TermsQuery: map[string]types.TermsQueryField{
			"type": cType,
		},
	}}
}

func createAndQuery(node *model.AstNode) (*types.Query, error) {
	q1, err := CreateSearchQuery(node.Left)
	if err != nil {
		return nil, err
	}
	q2, err := CreateSearchQuery(node.Right)
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

func createOrQuery(node *model.AstNode) (*types.Query, error) {
	q1, err := CreateSearchQuery(node.Left)
	if err != nil {
		return nil, err
	}
	q2, err := CreateSearchQuery(node.Right)
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

func createNotQuery(node *model.AstNode) (*types.Query, error) {
	q1, err := CreateSearchQuery(node.Left)
	if err != nil {
		return nil, err
	}
	if q1 == nil {
		q2, err := CreateSearchQuery(node.Right)
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

func createPhraseQuery(phrase string) *types.Query {
	return &types.Query{
		MatchPhrase: map[string]types.MatchPhraseQuery{
			"searchText": {Query: phrase},
		},
	}
}

func createTextMatchQuery(term string) *types.Query {
	return &types.Query{
		Match: map[string]types.MatchQuery{
			"searchText": {Query: term},
		},
	}
}

func CreateOptionQueries(opts model.SearchOptions) []types.Query {
	qs := []types.Query{}
	for _, wId := range opts.WorkIds {
		qs = append(qs, createWorkIdQuery(wId))
	}
	types := []esmodel.Type{esmodel.Paragraph}
	if opts.IncludeHeadings {
		types = append(types, esmodel.Heading)
	}
	if opts.IncludeFootnotes {
		types = append(types, esmodel.Footnote)
	}
	if opts.IncludeSummaries {
		types = append(types, esmodel.Summary)
	}
	qs = append(qs, createTypeQuery(types))
	return qs
}
