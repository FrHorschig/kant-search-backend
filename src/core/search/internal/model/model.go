package model

type AstNode struct {
	Left  *AstNode
	Right *AstNode
	Token *Token
}

type Token struct {
	IsAnd    bool
	IsOr     bool
	IsNot    bool
	IsOpen   bool
	IsClose  bool
	IsWord   bool
	IsPhrase bool
	Text     string
}
