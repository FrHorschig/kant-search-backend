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

type modelMapperImpl struct {
}

func NewModelMapper() ModelMapper {
	impl := modelMapperImpl{}
	return &impl
}

func (rec *modelMapperImpl) Map([]model.Section, []model.Summary, []model.Footnote) ([]dbmodel.Work, errors.ErrorNew) {
	// TODO implement me
	// don't forget handling of images and tables
	return nil, errors.NilError()
}

// TODO: handle levels H7 and H8 with some kind of error
// func mapLevel(lvl model.Level) dbmodel.Level {
// 	switch lvl {
// 	case model.H1:
// 		return dbmodel.HWork
// 	case model.H2:
// 		return dbmodel.H1
// 	case model.H3:
// 		return dbmodel.H2
// 	case model.H4:
// 		return dbmodel.H3
// 	case model.H5:
// 		return dbmodel.H4
// 	case model.H6:
// 		return dbmodel.H5
// 	}
// 	return dbmodel.H6 // model.H7
// }
