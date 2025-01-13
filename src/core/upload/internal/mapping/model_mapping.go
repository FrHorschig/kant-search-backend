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
	return nil
}
