package internalnew

import (
	"context"

	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/pyutil"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt1"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt2"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt31"
	"github.com/frhorschig/kant-search-backend/core/upload/model/abt32"
	"github.com/frhorschig/kant-search-backend/core/upload/model/vol14"
)

type XmlMapper interface {
	MapAbt1(ctx context.Context, volNum int32, vol abt1.Band) (model.Volume, error)
	MapAbt2(ctx context.Context, volNum int32, vol abt2.Band) (model.Volume, error)
	MapVol14(ctx context.Context, volNum int32, vol vol14.Band) (model.Volume, error)
	MapAbt31(ctx context.Context, volNum int32, vol abt31.Band) (model.Volume, error)
	MapAbt32(ctx context.Context, volNum int32, vol abt32.Band) (model.Volume, error)
}

type xmlMapperImpl struct {
	pyUtil pyutil.PythonUtil
}

func NewXmlMapper() XmlMapper {
	impl := xmlMapperImpl{
		pyUtil: pyutil.NewPythonUtil(),
	}
	return &impl
}

func (rec *xmlMapperImpl) MapAbt1(ctx context.Context, volNum int32, vol abt1.Band) (model.Volume, error) {
	return model.Volume{}, nil
}

func (rec *xmlMapperImpl) MapAbt2(ctx context.Context, volNum int32, vol abt2.Band) (model.Volume, error) {
	return model.Volume{}, nil
}

func (rec *xmlMapperImpl) MapVol14(ctx context.Context, volNum int32, vol vol14.Band) (model.Volume, error) {
	return model.Volume{}, nil
}

func (rec *xmlMapperImpl) MapAbt31(ctx context.Context, volNum int32, vol abt31.Band) (model.Volume, error) {
	return model.Volume{}, nil
}

func (rec *xmlMapperImpl) MapAbt32(ctx context.Context, volNum int32, vol abt32.Band) (model.Volume, error) {
	return model.Volume{}, nil
}
