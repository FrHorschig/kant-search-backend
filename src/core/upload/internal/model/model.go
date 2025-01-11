package model

import "github.com/frhorschig/kant-search-backend/common/model"

type Section struct {
	Heading    Heading
	Paragraphs []string
	Sections   []Section
	Parent     *Section
}

type Heading struct {
	TocTitle  string
	TextTitle string
	Level     model.Level
}

type Randtext struct {
	// TODO frhorschig implement me
}

type Footnote struct {
	// TODO frhorschig implement me
}
