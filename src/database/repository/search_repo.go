package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/search_repo_mock.go -package=mocks

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/lib/pq"
)

type SearchRepo interface {
	SearchParagraphs(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error)
	SearchSentences(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error)
}

type searchRepoImpl struct {
	db *sql.DB
}

func NewSearchRepo(db *sql.DB) SearchRepo {
	impl := searchRepoImpl{
		db: db,
	}
	return &impl
}

func (repo *searchRepoImpl) SearchParagraphs(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error) {
	snippetParams := `FragmentDelimiter="...<br>... ",
		MaxFragments=10,
		MaxWords=16,
		MinWords=6`
	textParams := `MaxWords=99999, MinWords=99998`
	query := `SELECT
			p.id, 
			ts_headline('german', p.content, plainto_tsquery('german', $2), $3),
			ts_headline('german', p.content, plainto_tsquery('german', $2), $4),
			p.pages,
			p.work_id
		FROM paragraphs p
		WHERE work_id = ANY($1) AND search @@ plainto_tsquery('german', $2)
		ORDER BY p.work_id, p.id`

	rows, err := repo.db.QueryContext(ctx, query, pq.Array(criteria.WorkIds), criteria.SearchTerms[0], snippetParams, textParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.SearchResult{}, nil
		}
		return nil, err
	}

	matches, err := scanSearchMatchRow(rows)
	if err != nil {
		return nil, err
	}
	return matches, err
}

func (repo *searchRepoImpl) SearchSentences(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error) {
	// TODO implement me
	return []model.SearchResult{}, nil
}

func scanSearchMatchRow(rows *sql.Rows) ([]model.SearchResult, error) {
	matches := make([]model.SearchResult, 0)
	for rows.Next() {
		var match model.SearchResult
		err := rows.Scan(&match.ElementId, &match.Snippet, &match.Text, pq.Array(&match.Pages), &match.WorkId)
		if err != nil {
			return nil, fmt.Errorf("search match row scan failed: %v", err)
		}
		matches = append(matches, match)
	}
	return matches, nil
}
