package dataaccess

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_repo_mock.go -package=mocks

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/esmodel"
)

type VolumeRepo interface {
	Insert(ctx context.Context, volNum int32, data esmodel.Volume) error
	Delete(ctx context.Context, volNum int32) error
}

type volumeRepoImpl struct {
	dbClient *elasticsearch.TypedClient
}

func NewVolumeRepo(dbClient *elasticsearch.TypedClient) VolumeRepo {
	return &volumeRepoImpl{
		dbClient: dbClient,
	}
}

func (rec *volumeRepoImpl) Insert(ctx context.Context, volNum int32, data esmodel.Volume) error {
	// TODO implement me
	return nil
}

func (rec *volumeRepoImpl) Delete(ctx context.Context, volNum int32) error {
	// TODO implement me
	return nil
}
