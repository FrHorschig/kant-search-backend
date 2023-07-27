package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/lib/pq"
)

type ParagraphRepo interface {
	Insert(ctx context.Context, paragraph model.Paragraph) (int32, error)
	UpdateText(ctx context.Context, paragraph model.Paragraph, reindex bool) error
	SelectOfPages(ctx context.Context, workId int32, start_id int32, end_id int32) ([]model.Paragraph, error)
	Search(ctx context.Context, searchTerms model.SearchCriteria) ([]model.Paragraph, error)
}

type ParagraphRepoImpl struct {
	db *sql.DB
}

func NewParagraphRepo(db *sql.DB) ParagraphRepo {
	return &ParagraphRepoImpl{
		db: db,
	}
}

func (repo *ParagraphRepoImpl) Insert(ctx context.Context, paragraphs model.Paragraph) (int32, error) {
	var id int32
	query := `INSERT INTO paragraphs (text, pages, work_id, reindex) VALUES ($1, $2, $3, $4) RETURNING id`
	row := repo.db.QueryRowContext(ctx, query, paragraphs.Text, pq.Array(paragraphs.Pages), paragraphs.WorkId, true)

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *ParagraphRepoImpl) UpdateText(ctx context.Context, paragraph model.Paragraph, reindex bool) error {
	query := `UPDATE paragraphs SET text = $1, reindex = $2 where id = $3`
	repo.db.ExecContext(ctx, query, paragraph.Text, reindex, paragraph.Id)
	return nil
}

func (repo *ParagraphRepoImpl) SelectOfPages(ctx context.Context, workId int32, start_id int32, end_id int32) ([]model.Paragraph, error) {
	query := `SELECT id, text, pages, work_id FROM paragraphs WHERE work_id = $1 AND $2 <= ANY(pages) AND $3 >= ANY(pages) ORDER BY id ASC`
	rows, err := repo.db.QueryContext(ctx, query, workId, start_id, end_id)
	if err != nil {
		return nil, err
	}

	paras, err := scanParagraphRows(rows)
	if err == sql.ErrNoRows {
		return []model.Paragraph{}, nil
	}
	return paras, err
}

func (repo *ParagraphRepoImpl) Search(ctx context.Context, searchCriteria model.SearchCriteria) ([]model.Paragraph, error) {
	searchString := strings.Join(searchCriteria.SearchWords, " & ")
	query := `SELECT id, text, pages, work_id FROM paragraphs WHERE work_id = ANY($1) AND search @@ to_tsquery('german', $2)`
	rows, err := repo.db.QueryContext(ctx, query, pq.Array(searchCriteria.WorkIds), searchString)
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
