package upload

//go:generate mockgen -source=$GOFILE -destination=mocks/volume_mock.go -package=mocks

import (
	"context"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadatamapping/metadatamapping/metadata"
	"github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
)

type UploadProcessor interface {
	Process(ctx context.Context, volNum int32, xml string) errs.UploadError
}

type uploadProcessorImpl struct {
	volumeRepo  dataaccess.VolumeRepo
	contentRepo dataaccess.ContentRepo
	xmlMapper   internal.XmlMapper
}

func NewUploadProcessor(volumeRepo dataaccess.VolumeRepo, contentRepo dataaccess.ContentRepo, configPath string) UploadProcessor {
	processor := uploadProcessorImpl{
		volumeRepo:  volumeRepo,
		contentRepo: contentRepo,
		xmlMapper:   internal.NewXmlMapper(metadata.NewMetadata(configPath)),
	}
	return &processor
}

func (rec *uploadProcessorImpl) Process(ctx context.Context, volNr int32, xml string) errs.UploadError {
	volume, contents, err := rec.xmlMapper.MapXml(volNr, xml)
	if err.HasError {
		return err
	}
	errDelete := deleteExistingData(ctx, rec.volumeRepo, rec.contentRepo, volNr)
	if errDelete != nil {
		return errs.New(nil, errDelete)
	}

	err = insertNewData(ctx, rec.volumeRepo, rec.contentRepo, &volume, contents)
	if err.HasError {
		deleteExistingData(ctx, rec.volumeRepo, rec.contentRepo, volNr) // ignore the potential delete error, because here the insertion error is the more interesting one
		return err
	}
	return errs.Nil()
}

func deleteExistingData(ctx context.Context, volRepo dataaccess.VolumeRepo, contentRepo dataaccess.ContentRepo, volNr int32) error {
	vol, err := volRepo.GetByVolumeNumber(ctx, volNr)
	if err != nil {
		return err
	}
	if vol == nil {
		return nil
	}
	for _, workRef := range vol.Works {
		err = contentRepo.DeleteByWork(ctx, workRef.Code)
		if err != nil {
			return err
		}
	}
	err = volRepo.Delete(ctx, volNr)
	if err != nil {
		return err
	}
	return nil
}

func insertNewData(ctx context.Context, volRepo dataaccess.VolumeRepo, contentRepo dataaccess.ContentRepo, volume *model.Volume, contents []model.Content) errs.UploadError {
	err := contentRepo.Insert(ctx, contents)
	if err != nil {
		return errs.New(nil, err)
	}
	err = volRepo.Insert(ctx, volume)
	if err != nil {
		return errs.New(nil, err)
	}
	return errs.Nil()
}
