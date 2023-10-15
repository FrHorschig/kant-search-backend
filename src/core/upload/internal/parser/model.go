package parser

type Type int

const (
	OPEN Type = iota
	CLOSE
	SEPARATOR
	CLASS
	LOCATION
	TEXT
)

type Token struct {
	Type Type
	Text string
}

func newOpen() Token {
	return Token{Type: OPEN, Text: "{"}
}
func newClose() Token {
	return Token{Type: CLOSE, Text: "}"}
}
func newSeparator() Token {
	return Token{Type: SEPARATOR, Text: "|"}
}
func newClass(text string) Token {
	return Token{Type: CLASS, Text: text}
}
func newLocation(text string) Token {
	return Token{Type: LOCATION, Text: text}
}
func newText(text string) Token {
	return Token{Type: TEXT, Text: text}
}

type Expression struct {
	Metadata Metadata
	Content  *Content
}

type Metadata struct {
	Class    string
	Location *string
}

type Content struct {
	Expressions []*Expression
	Texts       []string
}
