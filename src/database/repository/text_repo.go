package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/FrHorschig/kant-search-backend/database/models"
)

type TextRepo interface {
	Select(ctx context.Context, id int32) (models.Text, error)
}

type TextRepoImpl struct {
	db *sql.DB
}

func NewTextRepo(db *sql.DB) TextRepo {
	return &TextRepoImpl{
		db: db,
	}
}

func (repo *TextRepoImpl) Select(ctx context.Context, id int32) (models.Text, error) {
	var text models.Text
	query := `SELECT * FROM text_table WHERE id=$1`
	row := repo.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&text.Id, &text.Text)
	if err != nil {
		if err == sql.ErrNoRows {
			return text, fmt.Errorf("no rows found for id %d", id)
		}
		return text, fmt.Errorf("query row scan failed: %v", err)
	}

	return text, nil
}
