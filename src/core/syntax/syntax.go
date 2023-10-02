package syntax

func CheckSyntax(searchTerms string) (string, error) {
	tokens, err := tokenize(searchTerms)
	if err != nil {
		return "", err
	}
	return getSearchString(tokens), nil
}
