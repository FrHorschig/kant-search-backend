package validation

import (
	"github.com/frhorschig/kant-search-backend/api/search/internal/errors"
	"github.com/frhorschig/kant-search-backend/api/search/internal/validation/internal"
)

func CheckSyntax(searchTerms string) (string, *errors.ValidationError) {
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
