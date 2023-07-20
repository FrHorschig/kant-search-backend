package repository

import (
	"database/sql"
)

type SentenceRepo interface {
}

type SentenceRepoImpl struct {
	db *sql.DB
}

func NewSentenceRepo(db *sql.DB) SentenceRepo {
	return &SentenceRepoImpl{
		db: db,
	}
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
