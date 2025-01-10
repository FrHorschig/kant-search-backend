package errors

// TODO frhorschig: rename
type ErrorNew struct {
	HasError       bool
	DomainError    error
	TechnicalError error
}

func NewError(domainErr error, techError error) ErrorNew {
	hasError := false
	if domainErr != nil || techError != nil {
		hasError = true
	}
	return ErrorNew{
		HasError:       hasError,
		DomainError:    domainErr,
		TechnicalError: techError,
	}
}

func NilError() ErrorNew {
	return ErrorNew{false, nil, nil}
}

type Error struct {
	Msg    ErrMsg
	Params []string
}

type ErrMsg string

const (
	GO_ERR ErrMsg = "GO_ERR" // for normal `error` errors

	// search term parsing
	UNEXPECTED_TOKEN            ErrMsg = "UNEXPECTED_TOKEN"
	WRONG_STARTING_CHAR         ErrMsg = "WRONG_STARTING_CHAR"
	WRONG_ENDING_CHAR           ErrMsg = "WRONG_ENDING_CHAR"
	UNEXPECTED_END_OF_INPUT     ErrMsg = "UNEXPECTED_END_OF_INPUT"
	MISSING_CLOSING_PARENTHESIS ErrMsg = "MISSING_CLOSING_PARENTHESIS"
	UNTERMINATED_DOUBLE_QUOTE   ErrMsg = "UNTERMINATED_DOUBLE_QUOTE"

	// kantf parsing
	UPLOAD_GO_ERR              ErrMsg = "UPLOAD_GO_ERR"
	UPLOAD_WRONG_STARTING_CHAR ErrMsg = "UPLOAD_WRONG_STARTING_CHAR"
	MISSING_EXPR_TYPE          ErrMsg = "MISSING_EXPR_TYPE"
	MISSING_CLOSING_BRACE      ErrMsg = "MISSING_CLOSING_BRACE"
	WRONG_START_EXPRESSION     ErrMsg = "WRONG_START_EXPRESSION"
	WRONG_END_EXPRESSION       ErrMsg = "WRONG_END_EXPRESSION"
	UNKNOWN_EXPRESSION_CLASS   ErrMsg = "UNKNOWN_EXPRESSION_CLASS"
)
