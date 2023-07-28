package repository

import (
	"context"
	"database/sql"
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
	// TODO implement new search design
	searchString := strings.Join(searchCriteria.SearchWords, " & ")
	query := `SELECT id, text, pages, work_id FROM paragraphs WHERE work_id = ANY($1) AND search @@ to_tsquery('german', $2)`
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
		// TODO implement me
		matches = append(matches, match)
	}
	return matches, nil
}
