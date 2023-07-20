package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/database/model"
)

type SentenceRepo interface {
	Insert(ctx context.Context, sentences []model.Sentence) ([]int32, error)
}

type SentenceRepoImpl struct {
	db *sql.DB
}

func NewSentenceRepo(db *sql.DB) SentenceRepo {
	return &SentenceRepoImpl{
		db: db,
	}
}

func (repo *SentenceRepoImpl) Insert(ctx context.Context, sentences []model.Sentence) ([]int32, error) {
	query := `INSERT INTO sentences (text, paragraph_id, work_id) VALUES `
	values := make([]interface{}, 0)
	for i, sentence := range sentences {
		if i > 0 {
			query += `, `
		}
		query += `($` + fmt.Sprint(i*3+1) + `, $` + fmt.Sprint(i*3+2) + `, $` + fmt.Sprint(i*3+3) + `)`

		values = append(values, sentence.Text)
		values = append(values, sentence.ParagraphId)
		values = append(values, sentence.WorkId)
	}
	query += ` RETURNING id`

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	var ids []int32
	for rows.Next() {
		var id int32
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

/*
func scanSentenceRows(rows *sql.Rows) ([]model.Sentence, error) {
	paragraphs := make([]model.Sentence, 0)
	for rows.Next() {
		var work model.Sentence
		err := rows.Scan(&work.Id, &work.Text, &work.WorkId)
		if err != nil {
			return nil, fmt.Errorf("query row scan failed: %v", err)
		}
		paragraphs = append(paragraphs, work)
	}
	return paragraphs, nil
}
*/
