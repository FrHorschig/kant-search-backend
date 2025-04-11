package dataaccess

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/esmodel"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/work_repo_mock.go -package=mocks

type WorkRepo interface {
	Insert(ctx context.Context, data esmodel.Work) error
	Update(ctx context.Context, data esmodel.Work) error
	Delete(ctx context.Context, id int32) error
}

type workRepoImpl struct {
	dbClient *elasticsearch.TypedClient
}

func NewWorkRepo(dbClient *elasticsearch.TypedClient) WorkRepo {
	return &workRepoImpl{
		dbClient: dbClient,
	}
}

func (rec *workRepoImpl) Insert(ctx context.Context, data esmodel.Work) error {
	// TODO implement me
	return nil
}

func (rec *workRepoImpl) Update(ctx context.Context, data esmodel.Work) error {
	// TODO implement me
	return nil
}

func (rec *workRepoImpl) Delete(ctx context.Context, id int32) error {
	// TODO implement me
	return nil
}
