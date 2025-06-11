package util

// TODO all functions are only called in contentRepo
import (
	"errors"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

func CreateWorkCodeQuery(workCode string) types.Query {
	return types.Query{
		Term: map[string]types.TermQuery{
			"workCode": {Value: workCode},
		},
	}
}

func CreateContentQuery(workCode string, cType []esmodel.Type) *types.Query {
	return &types.Query{
		Bool: &types.BoolQuery{
			Filter: []types.Query{
				CreateWorkCodeQuery(workCode),
				createTypeQuery(cType),
			},
		},
	}
}

func CreateSearchQuery(node *model.AstNode, analyzer esmodel.Analyzer) (*types.Query, error) {
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

func CreateSortOptions() []types.SortCombinations {
	return []types.SortCombinations{
		types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"ordinal": {Order: &sortorder.Asc},
			},
		},
	}
}

func CreateHighlightOptions(analyzer esmodel.Analyzer) *types.Highlight {
	return &types.Highlight{
		Fields: map[string]types.HighlightField{
			"searchText." + string(analyzer): {
				FragmentSize:      util.IntPtr(150),
				NumberOfFragments: util.IntPtr(5),
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

func CreateOrdinalQuery(ordinals []int32) types.Query {
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

func createAndQuery(node *model.AstNode, analyzer esmodel.Analyzer) (*types.Query, error) {
	q1, err := CreateSearchQuery(node.Left, analyzer)
	if err != nil {
		return nil, err
	}
	q2, err := CreateSearchQuery(node.Right, analyzer)
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

func createOrQuery(node *model.AstNode, analyzer esmodel.Analyzer) (*types.Query, error) {
	q1, err := CreateSearchQuery(node.Left, analyzer)
	if err != nil {
		return nil, err
	}
	q2, err := CreateSearchQuery(node.Right, analyzer)
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

func createNotQuery(node *model.AstNode, analyzer esmodel.Analyzer) (*types.Query, error) {
	q1, err := CreateSearchQuery(node.Left, analyzer)
	if err != nil {
		return nil, err
	}
	if q1 == nil {
		q2, err := CreateSearchQuery(node.Right, analyzer)
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

func createPhraseQuery(phrase string, analyzer esmodel.Analyzer) *types.Query {
	return &types.Query{
		MatchPhrase: map[string]types.MatchPhraseQuery{
			"searchText." + string(analyzer): {Query: phrase},
		},
	}
}

func createTextMatchQuery(term string, analyzer esmodel.Analyzer) *types.Query {
	return &types.Query{
		Match: map[string]types.MatchQuery{
			"searchText." + string(analyzer): {Query: term},
		},
	}
}

func CreateOptionQueries(opts model.SearchOptions) []types.Query {
	tps := []esmodel.Type{esmodel.Paragraph}
	if opts.IncludeHeadings {
		tps = append(tps, esmodel.Heading)
	}
	if opts.IncludeFootnotes {
		tps = append(tps, esmodel.Footnote)
	}
	if opts.IncludeSummaries {
		tps = append(tps, esmodel.Summary)
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
