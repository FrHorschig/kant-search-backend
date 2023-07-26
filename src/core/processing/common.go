package processing

import "regexp"

func RemoveFormatting(text string) string {
	re := regexp.MustCompile("<[^>]*>")
	text = re.ReplaceAllString(text, " ")
	re = regexp.MustCompile(`\{[^}]*\}`)
	text = re.ReplaceAllString(text, " ")
	re = regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(text, " ")
}
