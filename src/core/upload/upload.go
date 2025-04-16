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
		deleteExistingData(ctx, rec.volumeRepo, rec.workRepo, rec.contentRepo, volNr) // ignore the error, because here insertion error is more interesting
		return errors.New(nil, errInsert)
	}
	return errors.Nil()
}

func deleteExistingData(ctx context.Context, volRepo dataaccess.VolumeRepo, workRepo dataaccess.WorkRepo, contentRepo dataaccess.ContentRepo, volNr int32) error {
	vol, err := volRepo.GetByVolumeNumber(ctx, volNr)
	if err != nil {
		return err
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

func insertNewData(ctx context.Context, volRepo dataaccess.VolumeRepo, workRepo dataaccess.WorkRepo, contentRepo dataaccess.ContentRepo, v *model.Volume, works []model.Work) error {
	vol := esmodel.Volume{
		VolumeNumber: v.VolumeNumber,
		Title:        v.Title,
	}
	for _, w := range works {
		work, err := insertWork(ctx, workRepo, w)
		if err != nil {
			return err
		}
		workId := work.Id
		sections, err := insertSections(ctx, contentRepo, w.Sections, workId)
		if err != nil {
			return err
		}
		work.Sections = append(work.Sections, sections...)
		err = workRepo.Update(ctx, &work)
		if err != nil {
			return err
		}

		err = insertFootnotes(ctx, contentRepo, w.Footnotes, workId)
		if err != nil {
			return err
		}
		err = insertSummaries(ctx, contentRepo, w.Summaries, workId)
		if err != nil {
			return err
		}
		vol.Works = append(vol.Works, createWorkRef(work))
	}
	return volRepo.Insert(ctx, &vol)
}

func insertWork(ctx context.Context, workRepo dataaccess.WorkRepo, w model.Work) (esmodel.Work, error) {
	work := esmodel.Work{
		Code:         w.Code,
		Abbreviation: w.Abbreviation,
		Title:        w.Title,
		Year:         w.Year,
		Sections:     []esmodel.Section{},
	}
	err := workRepo.Insert(ctx, &work)
	return work, err
}

func insertSections(ctx context.Context, contentRepo dataaccess.ContentRepo, sections []model.Section, workId string) ([]esmodel.Section, error) {
	result := []esmodel.Section{}
	for _, s := range sections {
		headId, err := insertHeading(ctx, contentRepo, &s.Heading, workId)
		if err != nil {
			return nil, err
		}
		parIds, err := insertParagraphs(ctx, contentRepo, s.Paragraphs, workId)
		if err != nil {
			return nil, err
		}
		secs, err := insertSections(ctx, contentRepo, s.Sections, workId)
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

func insertHeading(ctx context.Context, contentRepo dataaccess.ContentRepo, h *model.Heading, workId string) (string, error) {
	headings := []esmodel.Content{{
		Type:       esmodel.Heading,
		FmtText:    h.Text,
		TocText:    util.ToStrPtr(h.TocText),
		SearchText: extract.RemoveTags(h.Text),
		Pages:      h.Pages,
		FnRefs:     h.FnRefs,
		WorkId:     workId,
	}}
	err := contentRepo.Insert(ctx, headings)
	if err != nil {
		return "", err
	}
	return headings[0].Id, nil
}

func insertParagraphs(ctx context.Context, contentRepo dataaccess.ContentRepo, paragraphs []model.Paragraph, workId string) ([]string, error) {
	pars := []esmodel.Content{}
	for _, p := range paragraphs {
		pars = append(pars, esmodel.Content{
			Type:       esmodel.Paragraph,
			FmtText:    p.Text,
			SearchText: extract.RemoveTags(p.Text),
			Pages:      p.Pages,
			FnRefs:     p.FnRefs,
			SummaryRef: p.SummaryRef,
			WorkId:     workId,
		})
	}
	if len(pars) == 0 {
		return []string{}, nil
	}
	err := contentRepo.Insert(ctx, pars)
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(pars))
	for i, p := range pars {
		ids[i] = p.Id
	}
	return ids, nil
}

func insertFootnotes(ctx context.Context, contentRepo dataaccess.ContentRepo, footnotes []model.Footnote, workId string) error {
	fns := []esmodel.Content{}
	for _, f := range footnotes {
		fns = append(fns, esmodel.Content{
			Type:       esmodel.Footnote,
			Ref:        &f.Ref,
			FmtText:    f.Text,
			SearchText: extract.RemoveTags(f.Text),
			Pages:      f.Pages,
			WorkId:     workId,
		})
	}
	if len(fns) == 0 {
		return nil
	}
	return contentRepo.Insert(ctx, fns)
}

func insertSummaries(ctx context.Context, contentRepo dataaccess.ContentRepo, summaries []model.Summary, workId string) error {
	summs := []esmodel.Content{}
	for _, s := range summaries {
		summs = append(summs, esmodel.Content{
			Type:       esmodel.Summary,
			Ref:        &s.Ref,
			FmtText:    s.Text,
			SearchText: extract.RemoveTags(s.Text),
			Pages:      s.Pages,
			FnRefs:     s.FnRefs,
			WorkId:     workId,
		})
	}
	if len(summs) == 0 {
		return nil
	}
	return contentRepo.Insert(ctx, summs)
}

func createWorkRef(work esmodel.Work) esmodel.WorkRef {
	return esmodel.WorkRef{
		Id:    work.Id,
		Code:  work.Code,
		Title: work.Title,
	}
}
