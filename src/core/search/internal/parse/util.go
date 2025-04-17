package parse

import "github.com/frhorschig/kant-search-backend/dataaccess/model"

func newAnd() model.Token {
	return model.Token{IsAnd: true, Text: "&"}
}
func newOr() model.Token {
	return model.Token{IsOr: true, Text: "|"}
}
func newNot() model.Token {
	return model.Token{IsNot: true, Text: "!"}
}
func newOpen() model.Token {
	return model.Token{IsOpen: true, Text: "("}
}
func newClose() model.Token {
	return model.Token{IsClose: true, Text: ")"}
}
func newWord(text string) model.Token {
	return model.Token{IsWord: true, Text: text}
}
func newPhrase(text string) model.Token {
	return model.Token{IsPhrase: true, Text: text}
}
