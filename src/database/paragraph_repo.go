package database

//go:generate mockgen -source=$GOFILE -destination=mocks/paragraph_repo_mock.go -package=mocks

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/lib/pq"
)

type ParagraphRepo interface {
	Insert(ctx context.Context, paragraph model.Paragraph) (int32, error)
	SelectAll(ctx context.Context, workId int32) ([]model.Paragraph, error)
	Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error)
	DeleteByWorkId(ctx context.Context, workId int32) error
}

type paragraphRepoImpl struct {
	db *sql.DB
}

func NewParagraphRepo(db *sql.DB) ParagraphRepo {
	return &paragraphRepoImpl{
		db: db,
	}
}

func (repo *paragraphRepoImpl) Insert(ctx context.Context, paragraph model.Paragraph) (int32, error) {
	var id int32
	query := `INSERT INTO paragraphs (content, pages, work_id, heading_level, footnote_name) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := repo.db.QueryRowContext(ctx, query, paragraph.Text, pq.Array(paragraph.Pages), paragraph.WorkId, paragraph.HeadingLevel, paragraph.FootnoteName).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *paragraphRepoImpl) SelectAll(ctx context.Context, workId int32) ([]model.Paragraph, error) {
	query := `SELECT id, content, pages, work_id, heading_level, footnote_name FROM paragraphs WHERE work_id = $1 ORDER BY id`
	rows, err := repo.db.QueryContext(ctx, query, workId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Paragraph{}, nil
		}
		return nil, err
	}

	paras, err := scanParagraphRows(rows)
	if err != nil {
		return nil, err
	}
	return paras, err
}

func (repo *paragraphRepoImpl) Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error) {
	snippetParams, textParams := buildParams()
	query := `
		SELECT
			'... ' || ts_headline('german', p.content, to_tsquery('german', $2), $3) || ' ...',
			ts_headline('german', p.content, to_tsquery('german', $2), $4),
			p.pages,
			p.id, 
			p.work_id
		FROM paragraphs p
		WHERE p.work_id = ANY($1) AND p.search @@ to_tsquery('german', $2)
		ORDER BY p.work_id, p.id`

	rows, err := repo.db.QueryContext(ctx, query, pq.Array(criteria.WorkIds), criteria.SearchString, snippetParams, textParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.SearchResult{}, nil
		}
		return nil, err
	}

	return scanParagraphSearchMatchRow(rows)
}

func (repo *paragraphRepoImpl) DeleteByWorkId(ctx context.Context, workId int32) error {
	query := `DELETE FROM paragraphs WHERE work_id = $1`
	_, err := repo.db.ExecContext(ctx, query, workId)
	if err != nil {
		return err
	}
	return nil
}

func scanParagraphRows(rows *sql.Rows) ([]model.Paragraph, error) {
	paragraphs := make([]model.Paragraph, 0)
	for rows.Next() {
		var paragraph model.Paragraph
		var pages pq.Int64Array
		err := rows.Scan(&paragraph.Id, &paragraph.Text, &pages, &paragraph.WorkId, &paragraph.HeadingLevel, &paragraph.FootnoteName)
		if err != nil {
			return nil, fmt.Errorf("paragraph row scan failed: %v", err)
		}
		paragraph.Pages = make([]int32, len(pages))
		for i, page := range pages {
			paragraph.Pages[i] = int32(page)
		}
		paragraphs = append(paragraphs, paragraph)
	}
	return paragraphs, nil
}

func scanParagraphSearchMatchRow(rows *sql.Rows) ([]model.SearchResult, error) {
	matches := make([]model.SearchResult, 0)
	for rows.Next() {
		var match model.SearchResult
		err := rows.Scan(&match.Snippet, &match.Text, pq.Array(&match.Pages), &match.ParagraphId, &match.WorkId)
		if err != nil {
			return nil, fmt.Errorf("search match row scan failed: %v", err)
		}
		matches = append(matches, match)
	}
	return matches, nil
}
