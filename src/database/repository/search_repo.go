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
	SearchParagraphs(ctx context.Context, searchTerms model.SearchCriteria) ([]model.SearchResult, error)
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

func (repo *searchRepoImpl) SearchParagraphs(ctx context.Context, searchCriteria model.SearchCriteria) ([]model.SearchResult, error) {
	searchString := strings.Join(searchCriteria.SearchTerms, " & ")
	query := `SELECT p.id, ts_headline('german', p.content, to_tsquery('german', $2), 'FragmentDelimiter="...<br>... ", MaxFragments=10, MaxWords=16, MinWords=6'), p.pages, p.work_id
		FROM paragraphs p 
		WHERE work_id = ANY($1) AND search @@ to_tsquery('german', $2)
		ORDER BY p.work_id, p.ordinal`
	rows, err := repo.db.QueryContext(ctx, query, pq.Array(searchCriteria.WorkIds), searchString)
	if err != nil {
		return nil, err
	}

	matches, err := scanSearchMatchRow(rows)
	if err == sql.ErrNoRows {
		return []model.SearchResult{}, nil
	}
	return matches, err
}

func scanSearchMatchRow(rows *sql.Rows) ([]model.SearchResult, error) {
	matches := make([]model.SearchResult, 0)
	for rows.Next() {
		var match model.SearchResult
		err := rows.Scan(&match.ElementId, &match.Snippet, pq.Array(&match.Pages), &match.WorkId)
		if err != nil {
			return nil, fmt.Errorf("search match row scan failed: %v", err)
		}
		matches = append(matches, match)
	}
	return matches, nil
}
