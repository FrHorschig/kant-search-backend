package upload

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type WorkProcessor interface {
	Process(ctx context.Context, work model.Work) error
}

type WorkProcessorImpl struct {
	workRepo      repository.WorkRepo
	paragraphRepo repository.ParagraphRepo
	sentenceRepo  repository.SentenceRepo
}

func NewWorkProcessor(workRepo repository.WorkRepo, paragraphRepo repository.ParagraphRepo, sentenceRepo repository.SentenceRepo) WorkProcessor {
	processor := WorkProcessorImpl{
		workRepo:      workRepo,
		paragraphRepo: paragraphRepo,
		sentenceRepo:  sentenceRepo,
	}
	return &processor
}

func (proc *WorkProcessorImpl) Process(ctx context.Context, work model.Work) error {
	workId, err := proc.workRepo.Insert(ctx, work)
	if err != nil {
		return err
	}

	rawParas := strings.Split(work.Text, "{pr}")
	paras := make([]model.Paragraph, len(rawParas))
	for i, rawPara := range rawParas {
		para := strings.TrimSpace(rawPara)
		p, err := extractModelData(para, workId)
		if err != nil {
			return err
		}
		paras[i] = p
	}
	paras = mergeParasAroundFootnotes(paras)
	paras = removeEmptyParas(paras)

	// For now remove all line numbering
	r, _ := regexp.Compile(`\s*\{l\d+\}\s*`)
	for i := range paras {
		paras[i].Text = r.ReplaceAllString(paras[i].Text, " ")
	}

	// TODO write paragraphs to db
	for _, p := range paras {
		println(p.Text)
		println("-------------------")
	}
	return nil
}

func extractModelData(p string, workId int32) (model.Paragraph, error) {
	pages, err := findPages(p)
	if err != nil {
		return model.Paragraph{}, err
	}

	text, footnoteName, err := findTextOrFootnote(p)
	if err != nil {
		return model.Paragraph{}, err
	}

	para := model.Paragraph{Text: text, WorkId: workId, Pages: pages}
	if footnoteName != "" {
		para.FootnoteName = footnoteName
	}
	return para, nil
}

func findPages(p string) ([]int32, error) {
	pages := make([]int32, 0)
	r, _ := regexp.Compile(`\{p(\d+)\}`)
	matches := r.FindAllStringSubmatch(p, -1)
	for _, match := range matches {
		n, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
		pages = append(pages, int32(n))
		p = strings.ReplaceAll(p, match[0], "")
	}
	return pages, nil
}

func findTextOrFootnote(p string) (string, string, error) {
	r, _ := regexp.Compile(`\{fn(\d+\.\d+)\}\{(\w+)\}`)
	match := r.FindStringSubmatch(p)
	if len(match) > 0 {
		return match[2], match[1], nil
	}
	return p, "", nil
}

func mergeParasAroundFootnotes(paras []model.Paragraph) []model.Paragraph {
	merged := make([]model.Paragraph, 0)
	for i, p := range paras {
		if p.FootnoteName == "" && i > 0 && paras[i-1].FootnoteName != "" {
			lastNormal := findLastNonFootnote(merged)
			if lastNormal == -1 {
				merged = append(merged, model.Paragraph{Text: p.Text, WorkId: p.WorkId, Pages: p.Pages})
				continue
			}
			merged[lastNormal].Text += " " + p.Text
			merged[lastNormal].Pages = append(merged[lastNormal].Pages, p.Pages...)
		} else {
			merged = append(merged, model.Paragraph{Text: p.Text, WorkId: p.WorkId, Pages: p.Pages})
		}
	}
	return merged
}

func findLastNonFootnote(paras []model.Paragraph) int {
	for i := len(paras) - 1; i >= 0; i-- {
		if paras[i].FootnoteName == "" {
			return i
		}
	}
	return -1
}

func removeEmptyParas(paras []model.Paragraph) []model.Paragraph {
	filtered := make([]model.Paragraph, 0)
	for _, p := range paras {
		if p.Text != "" {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
