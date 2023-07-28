package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/core/model"
)

type WorkRepo interface {
	SelectAll(ctx context.Context) ([]model.Work, error)
	Insert(ctx context.Context, work model.Work) (int32, error)
}

type workRepoImpl struct {
	db *sql.DB
}

func NewWorkRepo(db *sql.DB) WorkRepo {
	return &workRepoImpl{
		db: db,
	}
}

func (repo *workRepoImpl) SelectAll(ctx context.Context) ([]model.Work, error) {
	query := `SELECT * FROM works ORDER BY ordinal`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	works, err := scanWorkRows(rows)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Work{}, nil
		}
		return nil, err
	}

	return works, nil
}

func (repo *workRepoImpl) Insert(ctx context.Context, work model.Work) (int32, error) {
	query := `INSERT INTO works (title, abbreviation, volume, ordinal, year) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	row := repo.db.QueryRowContext(ctx, query, work.Title, work.Abbreviation, work.Volume, work.Ordinal, work.Year)

	var id int32
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func scanWorkRows(rows *sql.Rows) ([]model.Work, error) {
	works := make([]model.Work, 0)
	for rows.Next() {
		var work model.Work
		err := rows.Scan(&work.Id, &work.Title, &work.Abbreviation, &work.Volume, &work.Ordinal, &work.Year)
		if err != nil {
			return nil, fmt.Errorf("query row scan failed: %v", err)
		}
		works = append(works, work)
	}
	return works, nil
}
