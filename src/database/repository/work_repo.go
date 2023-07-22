package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/database/model"
)

type WorkRepo interface {
	Select(ctx context.Context, id int32) (model.Work, error)
	SelectAll(ctx context.Context) ([]model.Work, error)
	Insert(ctx context.Context, work model.Work) (int32, error)
}

type WorkRepoImpl struct {
	db *sql.DB
}

func NewWorkRepo(db *sql.DB) WorkRepo {
	return &WorkRepoImpl{
		db: db,
	}
}

func (repo *WorkRepoImpl) Select(ctx context.Context, id int32) (model.Work, error) {
	var work model.Work
	query := `SELECT * FROM works WHERE id=$1`
	row := repo.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&work.Id, &work.Title, &work.Abbrev, &work.Volume, &work.Year)
	if err != nil {
		return work, err
	}

	return work, nil
}

func (repo *WorkRepoImpl) SelectAll(ctx context.Context) ([]model.Work, error) {
	query := `SELECT * FROM works`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	works, err := scanWorkRows(rows)
	if err != nil {
		return nil, err
	}

	return works, nil
}

func (repo *WorkRepoImpl) Insert(ctx context.Context, work model.Work) (int32, error) {
	query := `INSERT INTO works (title, abbrev, aa_volume, year) VALUES ($1, $2, $3, $4) RETURNING id`
	row := repo.db.QueryRowContext(ctx, query, work.Title, work.Abbrev, work.Volume, work.Year)

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
		err := rows.Scan(&work.Id, &work.Title, &work.Abbrev, &work.Volume, &work.Year)
		if err != nil {
			return nil, fmt.Errorf("query row scan failed: %v", err)
		}
		works = append(works, work)
	}
	return works, nil
}
