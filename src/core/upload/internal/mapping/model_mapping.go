package mapping

//go:generate mockgen -source=$GOFILE -destination=mocks/model_mapper.go -package=mocks

import "context"

type ModelMapper interface {
	MyFunc(ctx context.Context) error
}

type ModelMapperImpl struct {
}

func NewModelMapper() ModelMapper {
	impl := ModelMapperImpl{}
	return &impl
}

func (rec *ModelMapperImpl) MyFunc(ctx context.Context) error {
	// TODO implement me
	// don't forget to extract possible page numbers from headings and if necessary to add them to the start of the appropriate paragraphs
	return nil
}
