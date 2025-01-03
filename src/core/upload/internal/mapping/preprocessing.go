package mapping

import "regexp"

func Preprocess(xml []byte) (string, error) {
	// xml, err := pyutil.NewPythonUtil().PreprocessXml(xml)
	// if err != nil {
	//     return "", err
	// }

	str := string(xml)
	reZeile := regexp.MustCompile(`<zeile\s+nr="(\d+)"\s*/>`)
	str = reZeile.ReplaceAllString(str, `{l$1}`)
	reSeite := regexp.MustCompile(`<seite\s*[^>]\s*nr="(\d+)"\s*[^>]*\s*/>`)
	str = reSeite.ReplaceAllString(str, `{p$1}`)
	reTrenn := regexp.MustCompile(`<trenn\s*/>`)
	str = reTrenn.ReplaceAllString(str, "")

	return str, nil
}
