package dataaccess

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_repo_mock.go -package=mocks

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type VolumeRepo interface {
}

type volumeRepoImpl struct {
	dbClient *elasticsearch.TypedClient
}

func NewVolumeRepo(dbClient *elasticsearch.TypedClient) VolumeRepo {
	return &volumeRepoImpl{
		dbClient: dbClient,
	}
}
