package dataaccess

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/frhorschig/kant-search-backend/dataaccess/internal/esmodel"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/content_repo_mock.go -package=mocks

type ContentRepo interface {
	InsertHeading(ctx context.Context, data esmodel.Heading) error
	DeleteHeading(ctx context.Context, id int32) error
	InsertParagraph(ctx context.Context, data esmodel.Paragraph) error
	DeleteParagraph(ctx context.Context, id int32) error
	InsertFootnote(ctx context.Context, data esmodel.Footnote) error
	DeleteFootnote(ctx context.Context, id int32) error
	InsertSummary(ctx context.Context, data esmodel.Summary) error
	DeleteSummary(ctx context.Context, id int32) error
}

type contentRepoImpl struct {
	dbClient *elasticsearch.TypedClient
}

func NewContentRepo(dbClient *elasticsearch.TypedClient) ContentRepo {
	return &contentRepoImpl{
		dbClient: dbClient,
	}
}

func (rec *contentRepoImpl) InsertHeading(ctx context.Context, data esmodel.Heading) error {
	// TODO implement me
	return nil
}

func (rec *contentRepoImpl) DeleteHeading(ctx context.Context, id int32) error {
	// TODO implement me
	return nil
}

func (rec *contentRepoImpl) InsertParagraph(ctx context.Context, data esmodel.Paragraph) error {
	// TODO implement me
	return nil
}

func (rec *contentRepoImpl) DeleteParagraph(ctx context.Context, id int32) error {
	// TODO implement me
	return nil
}

func (rec *contentRepoImpl) InsertFootnote(ctx context.Context, data esmodel.Footnote) error {
	// TODO implement me
	return nil
}

func (rec *contentRepoImpl) DeleteFootnote(ctx context.Context, id int32) error {
	// TODO implement me
	return nil
}

func (rec *contentRepoImpl) InsertSummary(ctx context.Context, data esmodel.Summary) error {
	// TODO implement me
	return nil
}

func (rec *contentRepoImpl) DeleteSummary(ctx context.Context, id int32) error {
	// TODO implement me
	return nil
}
