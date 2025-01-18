package mapping

//go:generate mockgen -source=$GOFILE -destination=mocks/model_mapper.go -package=mocks

import (
	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/model"
	dbmodel "github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type ModelMapper interface {
	Map([]model.Section, []model.Summary, []model.Footnote) ([]dbmodel.Work, errors.ErrorNew)
}

type ModelMapperImpl struct {
}

func NewModelMapper() ModelMapper {
	impl := ModelMapperImpl{}
	return &impl
}

func (rec *ModelMapperImpl) Map([]model.Section, []model.Summary, []model.Footnote) ([]dbmodel.Work, errors.ErrorNew) {
	// TODO implement me
	// don't forget handling of images and tables
	return nil, errors.NilError()
}
