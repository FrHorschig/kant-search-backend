package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/FrHorschig/kant-search-backend/core/model"
	"github.com/lib/pq"
)

type SearchRepo interface {
	SearchParagraphs(ctx context.Context, searchTerms model.SearchCriteria) ([]model.SearchMatch, error)
}

type SearchRepoImpl struct {
	db *sql.DB
}

func NewSearchRepo(db *sql.DB) SearchRepo {
	impl := SearchRepoImpl{
		db: db,
	}
	return &impl
}

func (repo *SearchRepoImpl) SearchParagraphs(ctx context.Context, searchCriteria model.SearchCriteria) ([]model.SearchMatch, error) {
	searchString := strings.Join(searchCriteria.SearchTerms, " & ")
	query := `SELECT w.volume, w.title, p.text, p.pages, p.id
		FROM paragraphs p 
		JOIN works w ON p.work_id = w.id 
		WHERE work_id = ANY($1) AND search @@ to_tsquery('german', $2)
		ORDER BY w.volume, w.ordinal, p.id`
	rows, err := repo.db.QueryContext(ctx, query, pq.Array(searchCriteria.WorkIds), searchString)
	if err != nil {
		return nil, err
	}

	matches, err := scanSearchMatchRow(rows)
	if err == sql.ErrNoRows {
		return []model.SearchMatch{}, nil
	}
	return matches, err
}

func scanSearchMatchRow(rows *sql.Rows) ([]model.SearchMatch, error) {
	matches := make([]model.SearchMatch, 0)
	for rows.Next() {
		var match model.SearchMatch
		var pages pq.Int64Array
		err := rows.Scan(&match.Volume, &match.WorkTitle, &match.Snippet, &pages, &match.MatchId)
		if err != nil {
			return nil, fmt.Errorf("search match row scan failed: %v", err)
		}
		matches = append(matches, match)
	}
	return matches, nil
}
