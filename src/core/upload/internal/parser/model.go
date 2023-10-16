package parser

type Type int

const (
	OPEN Type = iota
	CLOSE
	SEPARATOR
	CLASS
	PARAM
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
func newParam(text string) Token {
	return Token{Type: PARAM, Text: text}
}
func newText(text string) Token {
	return Token{Type: TEXT, Text: text}
}

type Expression struct {
	Metadata Metadata
	Content  *string
}

func (e *Expression) String() string {
	s := e.Metadata.String()
	if e.Content != nil {
		s += "|" + *e.Content
	}
	return s
}

type Metadata struct {
	Class string
	Param *string
}

func (m *Metadata) String() string {
	s := m.Class
	if m.Param != nil {
		s += *m.Param
	}
	return s
}
