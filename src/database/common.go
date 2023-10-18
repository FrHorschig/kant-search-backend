package database

func buildParams() (snippetParams string, textParams string) {
	snippetParams = `FragmentDelimiter=" ...<br>... ",
		MaxFragments=10,
		MaxWords=16,
		MinWords=6`
	textParams = `MaxWords=99999, MinWords=99998`
	return
}
