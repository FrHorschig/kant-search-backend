package dataaccess

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_repo_mock.go -package=mocks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/deletebyquery"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/result"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/rs/zerolog/log"
)

type VolumeRepo interface {
	Insert(ctx context.Context, data *esmodel.Volume) error
	GetAll(ctx context.Context) ([]esmodel.Volume, error)
	GetByVolumeNumber(ctx context.Context, volNum int32) (*esmodel.Volume, error)
	Delete(ctx context.Context, volNr int32) error
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
	err := createVolumeIndex(repo.dbClient, repo.indexName)
	if err != nil {
		panic(err)
	}
	return repo
}

func createVolumeIndex(es *elasticsearch.TypedClient, name string) error {
	ctx := context.Background()
	ok, err := es.Indices.Exists(name).Do(ctx)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	res, err := es.Indices.Create(name).Request(&create.Request{
		Mappings: esmodel.VolumeMapping,
	}).Do(ctx)
	if err != nil {
		return err
	}
	if !res.Acknowledged {
		return fmt.Errorf("creation of index '%s' not acknowledged", name)
	}
	return err
}

// TODO disallow partial results everywhere

func (rec *volumeRepoImpl) Insert(ctx context.Context, data *esmodel.Volume) error {
	existing, err := rec.GetByVolumeNumber(ctx, data.VolumeNumber)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("volume with volume number %d already exists", data.VolumeNumber)
	}

	createRes, err := rec.dbClient.Index(rec.indexName).Document(&data).Do(ctx)
	if err != nil {
		return err
	}
	if createRes.Result != result.Created {
		return fmt.Errorf("unable to create volume with title \"%s\"", data.Title)
	}
	return err
}

func (rec *volumeRepoImpl) GetAll(ctx context.Context) ([]esmodel.Volume, error) {
	res, err := rec.dbClient.Search().Index(rec.indexName).
		Request(&search.Request{
			Query: &types.Query{MatchAll: &types.MatchAllQuery{}},
			Sort: []types.SortCombinations{
				types.SortOptions{
					SortOptions: map[string]types.FieldSort{
						"volumeNumber": {Order: &sortorder.Asc},
					},
				},
			},
		}).Do(ctx)
	if err != nil {
		return nil, err
	}

	volumes := []esmodel.Volume{}
	for _, hit := range res.Hits.Hits {
		var vol esmodel.Volume
		err = json.Unmarshal(hit.Source_, &vol)
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, vol)
	}
	return volumes, nil
}

func (rec *volumeRepoImpl) GetByVolumeNumber(ctx context.Context, volNum int32) (*esmodel.Volume, error) {
	res, err := rec.dbClient.Search().Index(rec.indexName).
		Request(&search.Request{
			Query: createTermQuery("volumeNumber", volNum),
		}).Do(ctx)
	if err != nil {
		return nil, err
	}
	numOfHits := len(res.Hits.Hits)
	if numOfHits == 0 {
		return nil, nil
	}
	if numOfHits > 1 {
		return nil, fmt.Errorf("more than one volume with volume number %d found", volNum)
	}

	var vol esmodel.Volume
	err = json.Unmarshal(res.Hits.Hits[0].Source_, &vol)
	if err != nil {
		return nil, err
	}
	return &vol, nil
}

func (rec *volumeRepoImpl) Delete(ctx context.Context, volNr int32) error {
	res, err := rec.dbClient.DeleteByQuery(rec.indexName).Request(&deletebyquery.Request{
		Query: createTermQuery("volumeNumber", volNr),
	}).Do(ctx)
	if err != nil {
		return err
	}

	if len(res.Failures) > 0 {
		e := res.Failures[0].Cause.Reason
		if e != nil {
			log.Error().Msgf("Failed to delete content: %s", *e)
		}
		return fmt.Errorf("unable to delete volume %d", volNr)
	}

	_, err = rec.dbClient.Indices.Refresh().Index(rec.indexName).Do(ctx)
	return err
}
