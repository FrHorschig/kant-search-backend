package preprocessing

import (
	"html"
	"regexp"
)

func ReplaceHtml(xml string) string {
	xml = html.UnescapeString(xml)
	replacements := map[string]string{
		"&alpha;":   "α",
		"&Alpha;":   "Α",
		"&beta;":    "β",
		"&Beta;":    "Β",
		"&gamma;":   "γ",
		"&Gamma;":   "Γ",
		"&delta;":   "δ",
		"&Delta;":   "Δ",
		"&epsilon;": "ε",
		"&Epsilon;": "Ε",
		"&zeta;":    "ζ",
		"&Zeta;":    "Ζ",
		"&eta;":     "η",
		"&Eta;":     "Η",
		"&theta;":   "θ",
		"&theata;":  "θ",
		"&Theta;":   "Θ",
		"&iota;":    "ι",
		"&Iota;":    "Ι",
		"&kappa;":   "κ",
		"&Kappa;":   "Κ",
		"&lambda;":  "λ",
		"&Lambda;":  "Λ",
		"&my;":      "μ",
		"&My;":      "Μ",
		"&ny;":      "ν",
		"&Ny;":      "Ν",
		"&xi;":      "ξ",
		"&Xi;":      "Ξ",
		"&omikron;": "ο",
		"&Omikron;": "Ο",
		"&pi;":      "π",
		"&Pi;":      "Π",
		"&rho;":     "ρ",
		"&Rho;":     "Ρ",
		"&sigma;":   "σ",
		"&sigma2;":  "ς",
		"&Sigma;":   "Σ",
		"&tau;":     "τ",
		"&Tau;":     "Τ",
		"&ypsilon;": "υ",
		"&Ypsilon;": "Υ",
		"&phi;":     "φ",
		"&Phi;":     "Φ",
		"&chi;":     "χ",
		"&Chi;":     "Χ",
		"&psi;":     "ψ",
		"&Psi;":     "Ψ",
		"&omega;":   "ω",
		"&Omega;":   "Ω",
	}
	for pattern, replacement := range replacements {
		re := regexp.MustCompile(pattern)
		xml = re.ReplaceAllString(xml, replacement)
	}
	return xml
}

func Simplify(xml string) string {
	reZeile := regexp.MustCompile(`<zeile\s+nr="(\d+)"\s*/>`)
	xml = reZeile.ReplaceAllString(xml, `{l$1}`)
	reSeite := regexp.MustCompile(`<seite\s*[^>]\s*nr="(\d+)"\s*[^>]*\s*/>`)
	xml = reSeite.ReplaceAllString(xml, `{p$1}`)
	reTrenn := regexp.MustCompile(`<trenn\s*/>`)
	xml = reTrenn.ReplaceAllString(xml, "")
	return xml
}
