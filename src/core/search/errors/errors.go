package errors

type SearchError struct {
	HasError       bool
	SyntaxError    *SyntaxError
	TechnicalError error
}

func New(syntaxErr *SyntaxError, technicalError error) SearchError {
	hasError := false
	if syntaxErr != nil || technicalError != nil {
		hasError = true
	}
	return SearchError{
		HasError:       hasError,
		SyntaxError:    syntaxErr,
		TechnicalError: technicalError,
	}
}

func Nil() SearchError {
	return SearchError{false, nil, nil}
}

type SyntaxError struct {
	Msg    ErrMsg
	Params []string
}

type ErrMsg string

const (
	UnexpectedToken         ErrMsg = "UNEXPECTED_TOKEN"
	WrongStartingChar       ErrMsg = "WRONG_STARTING_CHAR"
	WrongEndingChar         ErrMsg = "WRONG_ENDING_CHAR"
	UnexpectedEndOfInput    ErrMsg = "UNEXPECTED_END_OF_INPUT"
	MissingCloseParenthesis ErrMsg = "MISSING_CLOSING_PARENTHESIS"
	UnterminatedDoubleQuote ErrMsg = "UNTERMINATED_DOUBLE_QUOTE"
)
