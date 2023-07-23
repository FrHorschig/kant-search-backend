package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/lib/pq"
)

type ParagraphRepo interface {
	Select(ctx context.Context, id int32) (model.Paragraph, error)
	SelectRange(ctx context.Context, workId int32, start_id int32, end_id int32) ([]model.Paragraph, error)
	Insert(ctx context.Context, paragraph model.Paragraph) (int32, error)
}

type ParagraphRepoImpl struct {
	db *sql.DB
}

func NewParagraphRepo(db *sql.DB) ParagraphRepo {
	return &ParagraphRepoImpl{
		db: db,
	}
}

func (repo *ParagraphRepoImpl) Select(ctx context.Context, id int32) (model.Paragraph, error) {
	var paragraph model.Paragraph
	query := `SELECT * FROM paragraphs WHERE id=$1`
	row := repo.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&paragraph.Id, &paragraph.WorkId, &paragraph.Text)
	if err != nil {
		return paragraph, err
	}

	return paragraph, nil
}

func (repo *ParagraphRepoImpl) SelectRange(ctx context.Context, workId int32, start_id int32, end_id int32) ([]model.Paragraph, error) {
	var paragraphs []model.Paragraph
	query := `SELECT * FROM paragraphs WHERE work_id = $1 AND $2 <= ANY(pages) AND $3 >= ANY(pages)`
	rows, err := repo.db.QueryContext(ctx, query, workId, start_id, end_id)
	if err != nil {
		return nil, err
	}

	paragraphs, err = scanParagraphRows(rows)
	if err != nil {
		return nil, err
	}

	return paragraphs, nil
}

func (repo *ParagraphRepoImpl) Insert(ctx context.Context, paragraphs model.Paragraph) (int32, error) {
	var id int32
	query := `INSERT INTO paragraphs (text, pages, work_id) VALUES ($1, $2, $3) RETURNING id`
	row := repo.db.QueryRowContext(ctx, query, paragraphs.Text, pq.Array(paragraphs.Pages), paragraphs.WorkId)

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func scanParagraphRows(rows *sql.Rows) ([]model.Paragraph, error) {
	paragraphs := make([]model.Paragraph, 0)
	for rows.Next() {
		var work model.Paragraph
		var pages pq.Int64Array
		err := rows.Scan(&work.Id, &work.Text, &pages, &work.WorkId)
		if err != nil {
			return nil, fmt.Errorf("query row scan failed: %v", err)
		}
		work.Pages = make([]int32, len(pages))
		for i, page := range pages {
			work.Pages[i] = int32(page)
		}
		paragraphs = append(paragraphs, work)
	}
	return paragraphs, nil
}
