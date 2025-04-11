package dataaccess

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_repo_mock.go -package=mocks

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/util"
)

type VolumeRepo interface {
	Insert(ctx context.Context, volNum int32, data esmodel.Volume) error
	Delete(ctx context.Context, volNum int32) error
}

type volumeRepoImpl struct {
	dbClient  *elasticsearch.TypedClient
	indexName string
}

func NewVolumeRepo(dbClient *elasticsearch.TypedClient) VolumeRepo {
	repo := &volumeRepoImpl{
		dbClient:  dbClient,
		indexName: "volumes",
	}
	err := util.CreateIndex(repo.dbClient, repo.indexName, esmodel.VolumeMapping)
	if err != nil {
		panic(err)
	}
	return repo
}

func (rec *volumeRepoImpl) Insert(ctx context.Context, volNum int32, data esmodel.Volume) error {
	// TODO implement me
	return nil
}

func (rec *volumeRepoImpl) Delete(ctx context.Context, volNum int32) error {
	// TODO implement me
	return nil
}
