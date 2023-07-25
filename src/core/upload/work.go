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
	paras, err := parseParas(work.Text, workId)
	if err != nil {
		return err
	}

	// For now remove all line numbering
	r, _ := regexp.Compile(`\s*\{l\d+\}\s*`)
	for i := range paras {
		paras[i].Text = r.ReplaceAllString(paras[i].Text, " ")
	}

	for _, p := range paras {
		_, err := proc.paragraphRepo.Insert(ctx, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseParas(text string, workId int32) ([]model.Paragraph, error) {
	paras := make([]model.Paragraph, 0)
	lastTextPara := int32(0)
	lastPage := int32(0)
	lastIsFn := false
	for _, rawPara := range strings.Split(text, "{pr}") {
		para := strings.TrimSpace(rawPara)
		p, err := extractModelData(para, workId)
		if err != nil {
			return nil, err
		}

		if len(p.Pages) == 0 {
			p.Pages = append(p.Pages, lastPage)
		} else {
			lastPage = p.Pages[len(p.Pages)-1]
		}

		if isFn(p) {
			paras = append(paras, p)
			lastIsFn = true
		} else {
			if lastIsFn {
				paras[lastTextPara].Text += " " + p.Text
				paras[lastTextPara].Pages = append(paras[lastTextPara].Pages, p.Pages...)
			} else {
				paras = append(paras, p)
				lastTextPara = int32(len(paras) - 1)
			}
			lastIsFn = false
		}
	}
	return removeEmptyParas(paras), nil
}

func isFn(p model.Paragraph) bool {
	return p.FootnoteName != ""
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

	para := model.Paragraph{Text: text, WorkId: workId, Pages: pages, FootnoteName: footnoteName}
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
	r, _ := regexp.Compile(`\{fn(\d+\.\d+)\}\{([^}]+)\}`)
	match := r.FindStringSubmatch(p)
	if len(match) > 0 {
		return match[2], match[1], nil
	}
	return p, "", nil
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
