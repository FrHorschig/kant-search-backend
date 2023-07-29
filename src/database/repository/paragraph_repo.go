package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/lib/pq"
)

type ParagraphRepo interface {
	Insert(ctx context.Context, paragraph model.Paragraph) (int32, error)
	UpdateText(ctx context.Context, paragraph model.Paragraph, reindex bool) error
	Select(ctx context.Context, workId int32, paragraphId int32) (model.Paragraph, error)
	SelectOfPages(ctx context.Context, workId int32, page_start int32, page_end int32) ([]model.Paragraph, error)
}

type paragraphRepoImpl struct {
	db *sql.DB
}

func NewParagraphRepo(db *sql.DB) ParagraphRepo {
	return &paragraphRepoImpl{
		db: db,
	}
}

func (repo *paragraphRepoImpl) Insert(ctx context.Context, paragraphs model.Paragraph) (int32, error) {
	var id int32
	query := `INSERT INTO paragraphs (text, pages, work_id, reindex) VALUES ($1, $2, $3, $4) RETURNING id`
	row := repo.db.QueryRowContext(ctx, query, paragraphs.Text, pq.Array(paragraphs.Pages), paragraphs.WorkId, true)

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *paragraphRepoImpl) UpdateText(ctx context.Context, paragraph model.Paragraph, reindex bool) error {
	query := `UPDATE paragraphs SET text = $1, reindex = $2 where id = $3`
	repo.db.ExecContext(ctx, query, paragraph.Text, reindex, paragraph.Id)
	return nil
}

func (repo *paragraphRepoImpl) Select(ctx context.Context, workId int32, paragraphId int32) (model.Paragraph, error) {
	println(workId, paragraphId)
	var paragraph model.Paragraph
	query := `SELECT id, text, pages, work_id FROM paragraphs WHERE work_id = $1 AND id = $2`
	err := repo.db.QueryRowContext(ctx, query, workId, paragraphId).Scan(&paragraph.Id, &paragraph.Text, pq.Array(&paragraph.Pages), &paragraph.WorkId)
	return paragraph, err
}

func (repo *paragraphRepoImpl) SelectOfPages(ctx context.Context, workId int32, page_start int32, page_end int32) ([]model.Paragraph, error) {
	query := `SELECT id, text, pages, work_id FROM paragraphs WHERE work_id = $1 AND $2 <= ANY(pages) AND $3 >= ANY(pages) ORDER BY id`
	rows, err := repo.db.QueryContext(ctx, query, workId, page_start, page_end)
	if err != nil {
		return nil, err
	}

	paras, err := scanParagraphRows(rows)
	if err == sql.ErrNoRows {
		return []model.Paragraph{}, nil
	}
	return paras, err
}

func scanParagraphRows(rows *sql.Rows) ([]model.Paragraph, error) {
	paragraphs := make([]model.Paragraph, 0)
	for rows.Next() {
		var work model.Paragraph
		var pages pq.Int64Array
		err := rows.Scan(&work.Id, &work.Text, &pages, &work.WorkId)
		if err != nil {
			return nil, fmt.Errorf("paragraph row scan failed: %v", err)
		}
		work.Pages = make([]int32, len(pages))
		for i, page := range pages {
			work.Pages[i] = int32(page)
		}
		paragraphs = append(paragraphs, work)
	}
	return paragraphs, nil
}
