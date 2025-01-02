package validation

import (
	"github.com/frhorschig/kant-search-backend/api/search/validation/internal"
	"github.com/frhorschig/kant-search-backend/common/errors"
)

func CheckSyntax(searchTerms string) (string, *errors.Error) {
	tokens, err := internal.Tokenize(searchTerms)
	if err != nil {
		return "", err
	}
	err = internal.Parse(tokens)
	if err != nil {
		return "", err
	}
	return internal.GetSearchString(tokens), nil
}
