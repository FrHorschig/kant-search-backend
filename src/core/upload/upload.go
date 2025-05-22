package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/extract"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	"github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
)

type UploadProcessor interface {
	Process(ctx context.Context, volNum int32, xml string) errors.UploadError
}

type uploadProcessorImpl struct {
	volumeRepo  dataaccess.VolumeRepo
	workRepo    dataaccess.WorkRepo
	contentRepo dataaccess.ContentRepo
	xmlMapper   internal.XmlMapper
}

func NewUploadProcessor(volumeRepo dataaccess.VolumeRepo, workRepo dataaccess.WorkRepo, contentRepo dataaccess.ContentRepo) UploadProcessor {
	processor := uploadProcessorImpl{
		volumeRepo:  volumeRepo,
		workRepo:    workRepo,
		contentRepo: contentRepo,
		xmlMapper:   internal.NewXmlMapper(),
	}
	return &processor
}

func (rec *uploadProcessorImpl) Process(ctx context.Context, volNr int32, xml string) errors.UploadError {
	vol, err := rec.xmlMapper.MapVolume(volNr, xml)
	if err.HasError {
		return err
	}
	works, err := rec.xmlMapper.MapWorks(volNr, xml)
	if err.HasError {
		return err
	}

	errDelete := deleteExistingData(ctx, rec.volumeRepo, rec.workRepo, rec.contentRepo, volNr)
	if errDelete != nil {
		return errors.New(nil, errDelete)
	}

	errInsert := insertNewData(ctx, rec.volumeRepo, rec.workRepo, rec.contentRepo, vol, works)
	if errInsert != nil {
		deleteExistingData(ctx, rec.volumeRepo, rec.workRepo, rec.contentRepo, volNr) // ignore the error, because here the insertion error is the more interesting one
		return errors.New(nil, errInsert)
	}
	return errors.Nil()
}

func deleteExistingData(ctx context.Context, volRepo dataaccess.VolumeRepo, workRepo dataaccess.WorkRepo, contentRepo dataaccess.ContentRepo, volNr int32) error {
	vol, err := volRepo.GetByVolumeNumber(ctx, volNr)
	if err != nil {
		return err
	}
	if vol == nil {
		return nil
	}
	err = volRepo.Delete(ctx, volNr)
	if err != nil {
		return err
	}
	for _, workRef := range vol.Works {
		err = workRepo.Delete(ctx, workRef.Id)
		if err != nil {
			return err
		}
		err = contentRepo.DeleteByWorkId(ctx, workRef.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

type loopVariables struct {
	ctx         context.Context
	contentRepo dataaccess.ContentRepo
	fnByRef     map[string]model.Footnote
	summByRef   map[string]model.Summary
	workId      string
	ordinal     int32
}

func insertNewData(ctx context.Context, volRepo dataaccess.VolumeRepo, workRepo dataaccess.WorkRepo, contentRepo dataaccess.ContentRepo, v *model.Volume, works []model.Work) error {
	// TODO problem: using auto generated ids makes it so that the ids are different if the volume is uploaded again
	vol := esmodel.Volume{
		VolumeNumber: v.VolumeNumber,
		Title:        v.Title,
	}
	for i, w := range works {
		work := createWork(w, int32(i))
		err := workRepo.Insert(ctx, &work)
		if err != nil {
			return err
		}

		loopVars := &loopVariables{
			ctx:         ctx,
			contentRepo: contentRepo,
			fnByRef:     createFnByRef(w.Footnotes),
			summByRef:   createSummsByRef(w.Summaries),
			workId:      work.Id,
			ordinal:     int32(0),
		}

		paragraphs, err := insertParagraphs(loopVars, w.Paragraphs)
		if err != nil {
			return err
		}
		work.Paragraphs = append(work.Paragraphs, paragraphs...)
		sections, err := insertContents(loopVars, w.Sections)
		if err != nil {
			return err
		}
		work.Sections = append(work.Sections, sections...)
		err = workRepo.Update(ctx, &work)
		if err != nil {
			return err
		}
		vol.Works = append(vol.Works, createWorkRef(&work))
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
func createWork(w model.Work, ordinal int32) esmodel.Work {
	return esmodel.Work{
		Ordinal:      ordinal,
		Code:         w.Code,
		Abbreviation: w.Abbreviation,
		Title:        w.Title,
		Year:         w.Year,
		Paragraphs:   []string{},
		Sections:     []esmodel.Section{},
	}
}

func insertContents(lv *loopVariables, sections []model.Section) ([]esmodel.Section, error) {
	result := []esmodel.Section{}
	for _, s := range sections {
		headId, err := insertHeading(lv, s.Heading)
		if err != nil {
			return nil, err
		}
		parIds, err := insertParagraphs(lv, s.Paragraphs)
		if err != nil {
			return nil, err
		}
		secs, err := insertContents(lv, s.Sections)
		if err != nil {
			return nil, err
		}
		result = append(result, esmodel.Section{
			Heading:    headId,
			Paragraphs: parIds,
			Sections:   secs,
		})
	}
	return result, nil
}

func insertHeading(lv *loopVariables, h model.Heading) (string, error) {
	contents := []esmodel.Content{}
	contents = append(contents, createHeading(lv, h))
	for _, fnRef := range h.FnRefs {
		fn := lv.fnByRef[fnRef]
		contents = append(contents, createFootnote(lv, fn))
	}
	err := lv.contentRepo.Insert(lv.ctx, contents)
	if err != nil {
		return "", err
	}
	return contents[0].Id, nil
}

func insertParagraphs(lv *loopVariables, paragraphs []model.Paragraph) ([]string, error) {
	pars := []esmodel.Content{}
	fnsAndSumm := []esmodel.Content{}
	for _, p := range paragraphs {
		if p.SummaryRef != nil {
			summ := lv.summByRef[*p.SummaryRef]
			fnsAndSumm = append(fnsAndSumm, createSummary(lv, summ))
			for _, fnRef := range summ.FnRefs {
				fnsAndSumm = append(fnsAndSumm, createFootnote(lv, lv.fnByRef[fnRef]))
			}
		}
		pars = append(pars, createParagraph(lv, p))
		for _, fnRef := range p.FnRefs {
			fnsAndSumm = append(fnsAndSumm, createFootnote(lv, lv.fnByRef[fnRef]))
		}
	}
	if len(pars) == 0 {
		return []string{}, nil
	}

	if len(fnsAndSumm) > 0 {
		err := lv.contentRepo.Insert(lv.ctx, fnsAndSumm)
		if err != nil {
			return nil, err
		}
	}
	err := lv.contentRepo.Insert(lv.ctx, pars)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(pars))
	for i, p := range pars {
		ids[i] = p.Id
	}
	return ids, nil
}

func createHeading(lv *loopVariables, h model.Heading) esmodel.Content {
	heading := esmodel.Content{
		Type:       esmodel.Heading,
		Ordinal:    lv.ordinal,
		FmtText:    h.Text,
		TocText:    util.StrPtr(h.TocText),
		SearchText: extract.RemoveTags(h.Text),
		Pages:      h.Pages,
		FnRefs:     h.FnRefs,
		WorkId:     lv.workId,
	}
	lv.ordinal += 1
	return heading
}

func createParagraph(lv *loopVariables, p model.Paragraph) esmodel.Content {
	paragraph := esmodel.Content{
		Type:       esmodel.Paragraph,
		Ordinal:    lv.ordinal,
		FmtText:    p.Text,
		SearchText: extract.RemoveTags(p.Text),
		Pages:      p.Pages,
		FnRefs:     p.FnRefs,
		SummaryRef: p.SummaryRef,
		WorkId:     lv.workId,
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
		SearchText: extract.RemoveTags(f.Text),
		Pages:      f.Pages,
		WorkId:     lv.workId,
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
		SearchText: extract.RemoveTags(s.Text),
		Pages:      s.Pages,
		FnRefs:     s.FnRefs,
		WorkId:     lv.workId,
	}
	lv.ordinal += 1
	return summary
}

func createWorkRef(work *esmodel.Work) esmodel.WorkRef {
	return esmodel.WorkRef{
		Id:    work.Id,
		Code:  work.Code,
		Title: work.Title,
	}
}
