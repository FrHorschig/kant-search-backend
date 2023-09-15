package read

import (
	"context"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository"
)

type VolumeReader interface {
	FindAll(ctx context.Context) ([]model.Volume, error)
}

type volumeReaderImpl struct {
	volumeRepo repository.VolumeRepo
}

func NewVolumeReader(volumeRepo repository.VolumeRepo) VolumeReader {
	impl := volumeReaderImpl{volumeRepo: volumeRepo}
	return &impl
}

func (rec *volumeReaderImpl) FindAll(ctx context.Context) ([]model.Volume, error) {
	return rec.volumeRepo.SelectAll(ctx)
}
