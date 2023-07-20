package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/database/model"
)

type ParagraphRepo interface {
	Select(ctx context.Context, id int32) (model.Paragraph, error)
	SelectRange(ctx context.Context, start_id int32, end_id int32) ([]model.Paragraph, error)
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

func (repo *ParagraphRepoImpl) SelectRange(ctx context.Context, start_id int32, end_id int32) ([]model.Paragraph, error) {
	var paragraphs []model.Paragraph
	query := `SELECT * FROM paragraphs WHERE id BETWEEN $1 AND $2`
	rows, err := repo.db.QueryContext(ctx, query, start_id, end_id)
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
	row := repo.db.QueryRowContext(ctx, query, paragraphs.Text, paragraphs.Pages, paragraphs.WorkId)

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
		err := rows.Scan(&work.Id, &work.Text, &work.Pages, &work.WorkId)
		if err != nil {
			return nil, fmt.Errorf("query row scan failed: %v", err)
		}
		paragraphs = append(paragraphs, work)
	}
	return paragraphs, nil
}
