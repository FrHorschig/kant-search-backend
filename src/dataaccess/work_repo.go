package dataaccess

import "github.com/elastic/go-elasticsearch/v8"

//go:generate mockgen -source=$GOFILE -destination=mocks/work_repo_mock.go -package=mocks

type WorkRepo interface {
}

type workRepoImpl struct {
	dbClient *elasticsearch.TypedClient
}

func NewWorkRepo(dbClient *elasticsearch.TypedClient) WorkRepo {
	return &workRepoImpl{
		dbClient: dbClient,
	}
}
