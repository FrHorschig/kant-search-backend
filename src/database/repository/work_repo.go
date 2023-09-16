package repository

//go:generate mockgen -source=$GOFILE -destination=work_repo_mock.go -package=repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/database/model"
)

type WorkRepo interface {
	SelectAll(ctx context.Context) ([]model.Work, error)
	UpdateText(ctx context.Context, upload model.WorkUpload) error
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

func (repo *workRepoImpl) UpdateText(ctx context.Context, upload model.WorkUpload) error {
	query := `UPDATE works SET text = $1 WHERE id = $2;`
	_, err := repo.db.ExecContext(ctx, query, upload.Text, upload.WorkId)
	return err
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
