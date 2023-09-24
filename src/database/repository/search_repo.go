package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/search_repo_mock.go -package=mocks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	snippetParams, textParams := buildParams()
	query := `SELECT
			p.id, 
			ts_headline('german', p.content, to_tsquery('german', $2), $3),
			ts_headline('german', p.content, to_tsquery('german', $2), $4),
			p.pages,
			p.work_id
		FROM paragraphs p
		WHERE p.work_id = ANY($1) AND p.search @@ plainto_tsquery('german', $2)
		ORDER BY p.work_id, p.id`

	rows, err := repo.db.QueryContext(ctx, query, pq.Array(criteria.WorkIds), buildTerms(criteria), snippetParams, textParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.SearchResult{}, nil
		}
		return nil, err
	}

	return scanSearchMatchRow(rows)
}

func (repo *searchRepoImpl) SearchSentences(ctx context.Context, criteria model.SearchCriteria) ([]model.SearchResult, error) {
	snippetParams, textParams := buildParams()
	query := `SELECT
			s.id, 
			ts_headline('german', s.content, to_tsquery('german', $2), $3),
			ts_headline('german', s.content, to_tsquery('german', $2), $4),
			p.pages,
			p.work_id
		FROM sentences s
		LEFT JOIN paragraphs p ON s.paragraph_id = p.id
		WHERE p.work_id = ANY($1) AND s.search @@ plainto_tsquery('german', $2)
		ORDER BY p.work_id, s.id`

	rows, err := repo.db.QueryContext(ctx, query, pq.Array(criteria.WorkIds), buildTerms(criteria), snippetParams, textParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.SearchResult{}, nil
		}
		return nil, err
	}

	return scanSearchMatchRow(rows)
}

func buildParams() (snippetParams string, textParams string) {
	snippetParams = `FragmentDelimiter="...<br>... ",
		MaxFragments=10,
		MaxWords=16,
		MinWords=6`
	textParams = `MaxWords=99999, MinWords=99998`
	return
}

func buildTerms(c model.SearchCriteria) string {
	var builder strings.Builder
	builder.WriteString(strings.Join(c.SearchTerms, " & "))
	if len(c.ExcludedTerms) > 0 {
		builder.WriteString(" & !")
		builder.WriteString(strings.Join(c.ExcludedTerms, " & !"))
	}
	if len(c.OptionalTerms) > 0 {
		builder.WriteString(" | ")
		builder.WriteString(strings.Join(c.OptionalTerms, " | "))
	}
	return builder.String()
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
