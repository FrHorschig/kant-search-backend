package model

import "github.com/frhorschig/kant-search-backend/common/model"

type Work struct {
	Title    string
	Year     *string
	Sections []Section
}

type Section struct {
	Heading    Heading
	Paragraphs []string
	Sections   []Section
}

type Heading struct {
	Title string
	Level model.Level
}

type Randtext struct {
	// TODO frhorschig implement me
}

type Footnote struct {
	// TODO frhorschig implement me
}
