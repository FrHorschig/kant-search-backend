package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/sentence_repo_mock.go -package=mocks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/FrHorschig/kant-search-backend/common/model"
	"github.com/lib/pq"
)

type SentenceRepo interface {
	Insert(ctx context.Context, sentences []model.Sentence) ([]int32, error)
	Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error)
	DeleteByWorkId(ctx context.Context, workId int32) error
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

func (repo *sentenceRepoImpl) Search(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error) {
	snippetParams, textParams := buildParams()
	query := `
		SELECT
			'... ' || ts_headline('german', s.content, to_tsquery('german', $2), $3) || ' ...',
			ts_headline('german', s.content, to_tsquery('german', $2), $4),
			p.pages,
			s.id, 
			p.id,
			p.work_id
		FROM sentences s
		LEFT JOIN paragraphs p ON s.paragraph_id = p.id
		WHERE p.work_id = ANY($1) AND s.search @@ to_tsquery('german', $2)
		ORDER BY p.work_id, p.id, s.id`

	rows, err := repo.db.QueryContext(ctx, query, pq.Array(criteria.WorkIds), criteria.SearchString, snippetParams, textParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.SearchResult{}, nil
		}
		return nil, err
	}

	return scanSentenceSearchMatchRow(rows)
}

func (repo *sentenceRepoImpl) DeleteByWorkId(ctx context.Context, workId int32) error {
	query := `DELETE FROM sentences s USING paragraphs p WHERE s.paragraph_id = p.id AND p.work_id = $1`
	_, err := repo.db.ExecContext(ctx, query, workId)
	if err != nil {
		return err
	}
	return nil
}

func scanSentenceSearchMatchRow(rows *sql.Rows) ([]model.SearchResult, error) {
	matches := make([]model.SearchResult, 0)
	for rows.Next() {
		var match model.SearchResult
		err := rows.Scan(&match.Snippet, &match.Text, pq.Array(&match.Pages), &match.SentenceId, &match.ParagraphId, &match.WorkId)
		if err != nil {
			return nil, fmt.Errorf("search match row scan failed: %v", err)
		}
		matches = append(matches, match)
	}
	return matches, nil
}
