package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"

	commonutil "github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/errs"
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

	errInsert := insertNewData(ctx, rec.volumeRepo, rec.contentRepo, vol, works)
	if errInsert != nil {
		deleteExistingData(ctx, rec.volumeRepo, rec.contentRepo, volNr) // ignore the error, because here the insertion error is the more interesting one
		return errs.New(nil, errInsert)
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
	ctx         context.Context
	contentRepo dataaccess.ContentRepo
	fnByRef     map[string]model.Footnote
	summByRef   map[string]model.Summary
	workCode    string
	ordinal     int32
}

func insertNewData(ctx context.Context, volRepo dataaccess.VolumeRepo, contentRepo dataaccess.ContentRepo, v *model.Volume, works []model.Work) error {
	vol := esmodel.Volume{
		VolumeNumber: v.VolumeNumber,
		// TODO section
		Title: v.Title,
	}
	for i, w := range works {
		loopVars := &loopVariables{
			ctx:         ctx,
			contentRepo: contentRepo,
			fnByRef:     createFnByRef(w.Footnotes),
			summByRef:   createSummsByRef(w.Summaries),
			workCode:    w.Code,
			ordinal:     int32(1),
		}
		paragraphs, err := insertParagraphs(loopVars, w.Paragraphs)
		if err != nil {
			return err
		}
		sections, err := insertContents(loopVars, w.Sections)
		if err != nil {
			return err
		}
		work := createWork(w, int32(i+1), paragraphs, sections)
		vol.Works = append(vol.Works, work)
	}
	return volRepo.Insert(ctx, &vol)
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

func insertContents(lv *loopVariables, sections []model.Section) ([]esmodel.Section, error) {
	result := []esmodel.Section{}
	for _, s := range sections {
		headOrdinal, err := insertHeading(lv, s.Heading)
		if err != nil {
			return nil, err
		}
		parOrdinals, err := insertParagraphs(lv, s.Paragraphs)
		if err != nil {
			return nil, err
		}
		secs, err := insertContents(lv, s.Sections)
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

func insertHeading(lv *loopVariables, h model.Heading) (int32, error) {
	contents := make([]esmodel.Content, len(h.FnRefs)+1)
	contents[0] = createHeading(lv, h)
	for i, fnRef := range h.FnRefs {
		fn := lv.fnByRef[fnRef]
		contents[i+1] = createFootnote(lv, fn)
	}
	err := lv.contentRepo.Insert(lv.ctx, contents)
	if err != nil {
		return 0, err
	}
	return contents[0].Ordinal, nil
}

func insertParagraphs(lv *loopVariables, paragraphs []model.Paragraph) ([]int32, error) {
	if len(paragraphs) == 0 {
		return []int32{}, nil
	}
	parOrds := make([]int32, len(paragraphs))
	toInsert := []esmodel.Content{}
	for i, p := range paragraphs {
		if p.SummaryRef != nil {
			summ := lv.summByRef[*p.SummaryRef]
			toInsert = append(toInsert, createSummary(lv, summ))
			for _, fnRef := range summ.FnRefs {
				toInsert = append(toInsert, createFootnote(lv, lv.fnByRef[fnRef]))
			}
		}
		par := createParagraph(lv, p)
		parOrds[i] = par.Ordinal
		toInsert = append(toInsert, par)
		for _, fnRef := range p.FnRefs {
			toInsert = append(toInsert, createFootnote(lv, lv.fnByRef[fnRef]))
		}
	}
	err := lv.contentRepo.Insert(lv.ctx, toInsert)
	if err != nil {
		return nil, err
	}
	return parOrds, nil
}

func createHeading(lv *loopVariables, h model.Heading) esmodel.Content {
	heading := esmodel.Content{
		Type:       esmodel.Heading,
		Ordinal:    lv.ordinal,
		FmtText:    h.Text,
		TocText:    commonutil.StrPtr(h.TocText),
		SearchText: util.RemoveTags(h.Text),
		Pages:      h.Pages,
		FnRefs:     h.FnRefs,
		WorkCode:   lv.workCode,
	}
	lv.ordinal += 1
	return heading
}

func createParagraph(lv *loopVariables, p model.Paragraph) esmodel.Content {
	paragraph := esmodel.Content{
		Type:       esmodel.Paragraph,
		Ordinal:    lv.ordinal,
		FmtText:    p.Text,
		SearchText: util.RemoveTags(p.Text),
		Pages:      p.Pages,
		FnRefs:     p.FnRefs,
		SummaryRef: p.SummaryRef,
		WorkCode:   lv.workCode,
	}
	lv.ordinal += 1
	return paragraph
}

func createFootnote(lv *loopVariables, f model.Footnote) esmodel.Content {
	footnote := esmodel.Content{
		Type:       esmodel.Footnote,
		Ordinal:    lv.ordinal,
		Ref:        &f.Ref,
		FmtText:    f.Text,
		SearchText: util.RemoveTags(f.Text),
		Pages:      f.Pages,
		WorkCode:   lv.workCode,
	}
	lv.ordinal += 1
	return footnote
}

func createSummary(lv *loopVariables, s model.Summary) esmodel.Content {
	summary := esmodel.Content{
		Type:       esmodel.Summary,
		Ordinal:    lv.ordinal,
		Ref:        &s.Ref,
		FmtText:    s.Text,
		SearchText: util.RemoveTags(s.Text),
		Pages:      s.Pages,
		FnRefs:     s.FnRefs,
		WorkCode:   lv.workCode,
	}
	lv.ordinal += 1
	return summary
}
