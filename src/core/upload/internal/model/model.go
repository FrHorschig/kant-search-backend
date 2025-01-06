package model

type Section struct {
	Heading   Heading
	Paragraph []string
	Sections  []Section
}

type Heading struct {
	Name  string
	Level Level
}

type Level int

const (
	Work Level = iota + 1
	H1
	H2
	H3
	H4
	H5
	H6
	H7
)

type Footnote struct {
	// TODO frhorschig implement me
}
