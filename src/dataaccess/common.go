package dataaccess

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

func createTermQuery(fieldName string, value any) *types.Query {
	return &types.Query{
		Term: map[string]types.TermQuery{
			fieldName: {Value: value},
		},
	}
}
