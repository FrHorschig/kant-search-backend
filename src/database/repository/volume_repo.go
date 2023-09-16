package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_repo_mock.go -package=mocks

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/database/model"
)

type VolumeRepo interface {
	SelectAll(ctx context.Context) ([]model.Volume, error)
}

type volumeRepoImpl struct {
	db *sql.DB
}

func NewVolumeRepo(db *sql.DB) VolumeRepo {
	return &volumeRepoImpl{
		db: db,
	}
}

func (repo *volumeRepoImpl) SelectAll(ctx context.Context) ([]model.Volume, error) {
	query := `SELECT * FROM volumes ORDER BY id`
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	volumes, err := scanVolumeRows(rows)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Volume{}, nil
		}
		return nil, err
	}

	return volumes, nil
}

func scanVolumeRows(rows *sql.Rows) ([]model.Volume, error) {
	volumes := make([]model.Volume, 0)
	for rows.Next() {
		var volume model.Volume
		err := rows.Scan(&volume.Id, &volume.Title, &volume.Section)
		if err != nil {
			return nil, fmt.Errorf("query row scan failed: %v", err)
		}
		volumes = append(volumes, volume)
	}
	return volumes, nil
}
