package database

//go:generate mockgen -source=$GOFILE -destination=mocks/work_repo_mock.go -package=mocks

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/frhorschig/kant-search-backend/common/model"
)

type WorkRepo interface {
	SelectAll(ctx context.Context) ([]model.Work, error)
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
	query := `SELECT * FROM works ORDER BY volume_id, ordinal`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Work{}, nil
		}
		return nil, err
	}

	works, err := scanWorkRows(rows)
	if err != nil {
		return nil, err
	}

	return works, nil
}

func scanWorkRows(rows *sql.Rows) ([]model.Work, error) {
	works := make([]model.Work, 0)
	for rows.Next() {
		var work model.Work
		err := rows.Scan(&work.Id, &work.Title, &work.Abbreviation, &work.Ordinal, &work.Year, &work.Volume)
		if err != nil {
			return nil, fmt.Errorf("query row scan failed: %v", err)
		}
		works = append(works, work)
	}
	return works, nil
}
