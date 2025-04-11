package dataaccess

import "github.com/elastic/go-elasticsearch/v8"

//go:generate mockgen -source=$GOFILE -destination=mocks/content_repo_mock.go -package=mocks

type ContentRepo interface {
}

type contentRepoImpl struct {
	dbClient *elasticsearch.TypedClient
}

func NewContentRepo(dbClient *elasticsearch.TypedClient) ContentRepo {
	return &contentRepoImpl{
		dbClient: dbClient,
	}
}
