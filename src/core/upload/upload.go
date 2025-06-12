package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/frhorschig/kant-search-backend/common/errs"
	commonutil "github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internal"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadata"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/util"
	"github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
)

type UploadProcessor interface {
	Process(ctx context.Context, volNum int32, xml string) errs.UploadError
}

type uploadProcessorImpl struct {
	volumeRepo  dataaccess.VolumeRepo
	contentRepo dataaccess.ContentRepo
	xmlMapper   internal.XmlMapper
}

func NewUploadProcessor(volumeRepo dataaccess.VolumeRepo, contentRepo dataaccess.ContentRepo, configPath string) UploadProcessor {
	processor := uploadProcessorImpl{
		volumeRepo:  volumeRepo,
		contentRepo: contentRepo,
		xmlMapper:   internal.NewXmlMapper(metadata.NewMetadata(configPath)),
	}
	return &processor
}

func (rec *uploadProcessorImpl) Process(ctx context.Context, volNr int32, xml string) errs.UploadError {
	vol, err := rec.xmlMapper.MapVolume(volNr, xml)
	if err.HasError {
		return err
	}
	works, err := rec.xmlMapper.MapWorks(volNr, xml)
	if err.HasError {
		return err
	}

	errDelete := deleteExistingData(ctx, rec.volumeRepo, rec.contentRepo, volNr)
	if errDelete != nil {
		return errs.New(nil, errDelete)
	}

	// TODO extract esmodel building code to separate interface-impl

	err = insertNewData(ctx, rec.volumeRepo, rec.contentRepo, vol, works)
	if err.HasError {
		deleteExistingData(ctx, rec.volumeRepo, rec.contentRepo, volNr) // ignore the potential delete error, because here the insertion error is the more interesting one
		return err
	}
	return errs.Nil()
}

func deleteExistingData(ctx context.Context, volRepo dataaccess.VolumeRepo, contentRepo dataaccess.ContentRepo, volNr int32) error {
	vol, err := volRepo.GetByVolumeNumber(ctx, volNr)
	if err != nil {
		return err
	}
	if vol == nil {
		return nil
	}
	for _, workRef := range vol.Works {
		err = contentRepo.DeleteByWork(ctx, workRef.Code)
		if err != nil {
			return err
		}
	}
	err = volRepo.Delete(ctx, volNr)
	if err != nil {
		return err
	}
	return nil
}

type loopVariables struct {
	fnByRef   map[string]model.Footnote
	summByRef map[string]model.Summary
	workCode  string
	ordinal   int32
	contents  []esmodel.Content
}

func insertNewData(ctx context.Context, volRepo dataaccess.VolumeRepo, contentRepo dataaccess.ContentRepo, v *model.Volume, works []model.Work) errs.UploadError {
	vol := esmodel.Volume{
		VolumeNumber: v.VolumeNumber,
		Title:        v.Title,
	}

	loopVars := &loopVariables{
		contents: []esmodel.Content{},
	}
	for i, w := range works {
		loopVars.fnByRef = createFnByRef(w.Footnotes)
		loopVars.summByRef = createSummsByRef(w.Summaries)
		loopVars.workCode = w.Code
		loopVars.ordinal = int32(1)
		paragraphs, err := buildParagraphs(loopVars, w.Paragraphs)
		if err != nil {
			return errs.New(nil, err)
		}
		sections, err := buildContents(loopVars, w.Sections)
		if err != nil {
			return errs.New(nil, err)
		}
		work := createWork(w, int32(i+1), paragraphs, sections)
		vol.Works = append(vol.Works, work)
	}

	err := contentRepo.Insert(ctx, loopVars.contents)
	if err != nil {
		return errs.New(nil, err)
	}
	err = volRepo.Insert(ctx, &vol)
	if err != nil {
		return errs.New(nil, err)
	}
	return errs.Nil()
}

func createFnByRef(footnotes []model.Footnote) map[string]model.Footnote {
	result := make(map[string]model.Footnote, len(footnotes))
	for _, fn := range footnotes {
		result[fn.Ref] = fn
	}
	return result
}

func createSummsByRef(summaries []model.Summary) map[string]model.Summary {
	result := make(map[string]model.Summary, len(summaries))
	for _, summ := range summaries {
		result[summ.Ref] = summ
	}
	return result
}
func createWork(w model.Work, ordinal int32, pars []int32, secs []esmodel.Section) esmodel.Work {
	return esmodel.Work{
		Ordinal:      ordinal,
		Code:         w.Code,
		Abbreviation: w.Abbreviation,
		Title:        w.Title,
		Year:         w.Year,
		Paragraphs:   pars,
		Sections:     secs,
	}
}

func buildContents(lv *loopVariables, sections []model.Section) ([]esmodel.Section, error) {
	result := []esmodel.Section{}
	for _, s := range sections {
		headOrdinal, err := buildHeading(lv, s.Heading)
		if err != nil {
			return nil, err
		}
		parOrdinals, err := buildParagraphs(lv, s.Paragraphs)
		if err != nil {
			return nil, err
		}
		secs, err := buildContents(lv, s.Sections)
		if err != nil {
			return nil, err
		}
		resultSection := esmodel.Section{
			Heading:    headOrdinal,
			Paragraphs: parOrdinals,
		}
		if len(secs) > 0 {
			resultSection.Sections = secs
		}
		result = append(result, resultSection)
	}
	return result, nil
}

func buildHeading(lv *loopVariables, heading model.Heading) (int32, error) {
	contents := make([]esmodel.Content, len(heading.FnRefs)+1)
	h, err := createHeading(lv, heading)
	if err != nil {
		return 0, err
	}
	contents[0] = h
	for i, fnRef := range heading.FnRefs {
		f, err := createFootnote(lv, lv.fnByRef[fnRef])
		if err != nil {
			return 0, err
		}
		contents[i+1] = f
	}
	lv.contents = append(lv.contents, contents...)
	return contents[0].Ordinal, nil
}

func buildParagraphs(lv *loopVariables, paragraphs []model.Paragraph) ([]int32, error) {
	if len(paragraphs) == 0 {
		return []int32{}, nil
	}
	parOrds := make([]int32, len(paragraphs))
	toInsert := []esmodel.Content{}
	for i, par := range paragraphs {
		if par.SummaryRef != nil {
			summ := lv.summByRef[*par.SummaryRef]
			s, err := createSummary(lv, summ)
			if err != nil {
				return nil, err
			}
			toInsert = append(toInsert, s)
			for _, fnRef := range summ.FnRefs {
				f, err := createFootnote(lv, lv.fnByRef[fnRef])
				if err != nil {
					return nil, err
				}
				toInsert = append(toInsert, f)
			}
		}
		p, err := createParagraph(lv, par)
		if err != nil {
			return nil, err
		}
		parOrds[i] = p.Ordinal
		toInsert = append(toInsert, p)
		for _, fnRef := range p.FnRefs {
			f, err := createFootnote(lv, lv.fnByRef[fnRef])
			if err != nil {
				return nil, err
			}
			toInsert = append(toInsert, f)
		}
	}
	lv.contents = append(lv.contents, toInsert...)
	return parOrds, nil
}

func createHeading(lv *loopVariables, h model.Heading) (esmodel.Content, error) {
	searchText := util.RemoveTags(h.Text)
	wordIndexMap, err := findWordIndexMap(searchText, h.Text)
	if err != nil {
		return esmodel.Content{}, err
	}
	pageByIndex, err := findNumberByIndex(h.Text, util.PageMatch)
	if err != nil {
		return esmodel.Content{}, err
	}
	lineByIndex, err := findNumberByIndex(h.Text, util.LineMatch)
	if err != nil {
		return esmodel.Content{}, err
	}
	heading := esmodel.Content{
		Type:         esmodel.Heading,
		Ordinal:      lv.ordinal,
		FmtText:      h.Text,
		TocText:      commonutil.StrPtr(h.TocText),
		SearchText:   searchText,
		WordIndexMap: wordIndexMap,
		PageByIndex:  pageByIndex,
		LineByIndex:  lineByIndex,
		Pages:        h.Pages,
		FnRefs:       h.FnRefs,
		WorkCode:     lv.workCode,
	}
	lv.ordinal += 1
	return heading, nil
}

func createParagraph(lv *loopVariables, p model.Paragraph) (esmodel.Content, error) {
	searchText := util.RemoveTags(p.Text)
	wordIndexMap, err := findWordIndexMap(searchText, p.Text)
	if err != nil {
		return esmodel.Content{}, err
	}
	pageByIndex, err := findNumberByIndex(p.Text, util.PageMatch)
	if err != nil {
		return esmodel.Content{}, err
	}
	lineByIndex, err := findNumberByIndex(p.Text, util.LineMatch)
	if err != nil {
		return esmodel.Content{}, err
	}
	paragraph := esmodel.Content{
		Type:         esmodel.Paragraph,
		Ordinal:      lv.ordinal,
		FmtText:      p.Text,
		SearchText:   searchText,
		WordIndexMap: wordIndexMap,
		Pages:        p.Pages,
		PageByIndex:  pageByIndex,
		LineByIndex:  lineByIndex,
		FnRefs:       p.FnRefs,
		SummaryRef:   p.SummaryRef,
		WorkCode:     lv.workCode,
	}
	lv.ordinal += 1
	return paragraph, nil
}

func createFootnote(lv *loopVariables, f model.Footnote) (esmodel.Content, error) {
	searchText := util.RemoveTags(f.Text)
	wordIndexMap, err := findWordIndexMap(searchText, f.Text)
	if err != nil {
		return esmodel.Content{}, err
	}
	pageByIndex, err := findNumberByIndex(f.Text, util.PageMatch)
	if err != nil {
		return esmodel.Content{}, err
	}
	lineByIndex, err := findNumberByIndex(f.Text, util.LineMatch)
	if err != nil {
		return esmodel.Content{}, err
	}
	footnote := esmodel.Content{
		Type:         esmodel.Footnote,
		Ordinal:      lv.ordinal,
		Ref:          &f.Ref,
		FmtText:      f.Text,
		SearchText:   searchText,
		WordIndexMap: wordIndexMap,
		Pages:        f.Pages,
		PageByIndex:  pageByIndex,
		LineByIndex:  lineByIndex,
		WorkCode:     lv.workCode,
	}
	lv.ordinal += 1
	return footnote, nil
}

func createSummary(lv *loopVariables, s model.Summary) (esmodel.Content, error) {
	searchText := util.RemoveTags(s.Text)
	wordIndexMap, err := findWordIndexMap(searchText, s.Text)
	if err != nil {
		return esmodel.Content{}, err
	}
	pageByIndex, err := findNumberByIndex(s.Text, util.PageMatch)
	if err != nil {
		return esmodel.Content{}, err
	}
	lineByIndex, err := findNumberByIndex(s.Text, util.LineMatch)
	if err != nil {
		return esmodel.Content{}, err
	}
	summary := esmodel.Content{
		Type:         esmodel.Summary,
		Ordinal:      lv.ordinal,
		Ref:          &s.Ref,
		FmtText:      s.Text,
		SearchText:   searchText,
		WordIndexMap: wordIndexMap,
		Pages:        s.Pages,
		PageByIndex:  pageByIndex,
		LineByIndex:  lineByIndex,
		FnRefs:       s.FnRefs,
		WorkCode:     lv.workCode,
	}
	lv.ordinal += 1
	return summary, nil
}

type WordData struct {
	Text  string
	Index int32
}

func findWordIndexMap(rawText string, fmtText string) (map[int32]int32, error) {
	rawWords := extractWordData(rawText)
	fmtWords := extractWordData(util.MaskTags(fmtText)) // we mask the tags here so that no word from a tag in fmtText may be accidentally matched to a normal word in rawText; this could happen if we later introduce tags with metadata (like image or table tags)
	if len(rawWords) != len(fmtWords) {
		// and because we mask all known tags, we (should) get the exact same number of words in both cases...
		return nil, fmt.Errorf("unequal number of words in searchText and fmtText: {%s} vs {%s}", fmtText, rawText)
	}

	result := make(map[int32]int32)
	for i, rawWord := range rawWords {
		fmtWord := fmtWords[i]
		if rawWord.Text != fmtWord.Text { // just one more sanity check
			return nil, fmt.Errorf("unequal matched words '%s' at index %d in searchText and '%s' at %d in fmtText", rawWord.Text, rawWord.Index, fmtWord.Text, fmtWord.Index)
		}
		// ...which then makes it super simple to map the indices
		result[rawWord.Index] = fmtWord.Index
	}
	return result, nil
}

func extractWordData(text string) []WordData {
	words := []WordData{}
	var currentWord strings.Builder
	startIndex := 0
	for i, r := range text {
		if unicode.IsLetter(r) {
			if currentWord.Len() == 0 {
				startIndex = i
			}
			currentWord.WriteRune(r)
		} else if currentWord.Len() > 0 {
			words = append(words, WordData{
				Text:  currentWord.String(),
				Index: int32(startIndex),
			})
			currentWord.Reset()
		}
	}
	if currentWord.Len() > 0 {
		words = append(words, WordData{
			Text:  currentWord.String(),
			Index: int32(startIndex),
		})
	}
	return words
}

func findNumberByIndex(text string, matcher string) ([]esmodel.IndexNumberPair, error) {
	result := []esmodel.IndexNumberPair{}
	re := regexp.MustCompile(matcher)
	for _, match := range re.FindAllStringSubmatchIndex(text, -1) {
		num, err := strconv.ParseInt(text[match[2]:match[3]], 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, esmodel.IndexNumberPair{
			I:   int32(match[0]),
			Num: int32(num),
		})
	}
	return result, nil
}
