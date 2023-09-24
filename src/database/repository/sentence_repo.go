package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/sentence_repo_mock.go -package=mocks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/FrHorschig/kant-search-backend/database/model"
)

type SentenceRepo interface {
	Insert(ctx context.Context, sentences []model.Sentence) ([]int32, error)
}

type sentenceRepoImpl struct {
	db *sql.DB
}

func NewSentenceRepo(db *sql.DB) SentenceRepo {
	return &sentenceRepoImpl{
		db: db,
	}
}

func (repo *sentenceRepoImpl) Insert(ctx context.Context, sentences []model.Sentence) ([]int32, error) {
	var builder strings.Builder
	builder.WriteString(`INSERT INTO sentences (content, paragraph_id) VALUES `)
	values := make([]interface{}, 0)
	for i, sentence := range sentences {
		if i > 0 {
			builder.WriteString(`, `)
		}
		builder.WriteString(`($` + fmt.Sprint(i*2+1) + `, $` + fmt.Sprint(i*2+2) + `)`)
		values = append(values, sentence.Text)
		values = append(values, sentence.ParagraphId)
	}
	builder.WriteString(` RETURNING id`)

	rows, err := repo.db.QueryContext(ctx, builder.String(), values...)
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
