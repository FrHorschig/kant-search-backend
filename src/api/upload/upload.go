package upload

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/beevik/etree"
	"github.com/frhorschig/kant-search-backend/api/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type UploadHandler interface {
	PostVolume(ctx echo.Context) error
}

type uploadHandlerImpl struct {
	volumeProcessor upload.UploadProcessor
}

func NewUploadHandler(volumeProcessor upload.UploadProcessor) UploadHandler {
	return &uploadHandlerImpl{
		volumeProcessor: volumeProcessor,
	}
}

func (rec *uploadHandlerImpl) PostVolume(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		msg := fmt.Sprintf("error reading request body: %v", err.Error())
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	xml, err := readToString(body)
	if err != nil {
		msg := fmt.Sprintf("error unmarshaling request body: %v", err.Error())
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	xml = replaceHtml(xml)
	xml = replaceCustomEntities(xml)

	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		msg := fmt.Sprintf("error unmarshaling request body: %v", err.Error())
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}

	band := doc.FindElement("//band")
	if band == nil || band.SelectAttr("nr") == nil {
		msg := "missing element 'band' with attribute 'nr'"
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	nrStr := strings.TrimLeft(band.SelectAttr("nr").Value, "0")
	if nrStr == "" {
		msg := "the volume number is 0, but it must be between 1 and 9"
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	volNum, err := strconv.Atoi(nrStr)
	if err != nil {
		msg := fmt.Sprintf("attribute 'nr' of element 'band' can't be converted to a number: %v", err.Error())
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	}
	if volNum < 1 {
		msg := fmt.Sprintf("the volume number is %d, but it must be between 1 and 9", volNum)
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusBadRequest, msg)
	} else if volNum > 9 {
		msg := "uploading volumes greater than 9 is not yet implemented"
		log.Error().Msg(msg)
		return errors.JsonError(ctx, http.StatusNotImplemented, msg)
	}

	if err := rec.volumeProcessor.Process(ctx.Request().Context(), doc); err != nil {
		msg := fmt.Sprintf("error processing XML data for volume %d", volNum)
		log.Error().Err(err).Msg(msg)
		return errors.JsonError(ctx, http.StatusInternalServerError, msg)
	}

	return ctx.NoContent(http.StatusCreated)
}

func readToString(input []byte) (string, error) {
	xmlDecl := `<?xml version="1.0" encoding="`
	encodingStart := bytes.Index(input, []byte(xmlDecl))
	if encodingStart == -1 {
		return string(input), nil
	}
	encodingStart += len(xmlDecl)
	encodingEnd := bytes.IndexByte(input[encodingStart:], '"')
	if encodingEnd == -1 {
		return "", fmt.Errorf("invalid XML declaration, no closing quote for encoding")
	}

	encoding := string(input[encodingStart : encodingStart+encodingEnd])
	if encoding == "" || encoding == "UTF-8" {
		return string(input), nil
	} else if encoding == "ISO-8859-1" {
		decoder := charmap.ISO8859_1.NewDecoder()
		utf8Bytes, _, err := transform.Bytes(decoder, input)
		if err != nil {
			return "", err
		}
		return string(utf8Bytes), nil
	} else {
		return "", fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func replaceCustomEntities(xml string) string {
	customEntities := map[string]string{
		"&kreis;":   "○",
		"&quadrat;": "■",
	}
	for entity, replacement := range customEntities {
		xml = strings.ReplaceAll(xml, entity, replacement)
	}
	return xml
}

func replaceHtml(xml string) string {
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
