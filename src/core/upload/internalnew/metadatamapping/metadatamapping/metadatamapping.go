package metadatamapping

import (
	"fmt"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/metadatamapping/metadatamapping/metadata"
)

func MapMetadata(volume *model.Volume, works []model.Work, metadata metadata.Metadata) errs.UploadError {
	md, mdErr := metadata.Read(volume.VolumeNumber)
	if mdErr != nil {
		return errs.New(nil, mdErr)
	}
	volume.Title = md.Title

	err := addWorkMetadata(works, md)
	if err.HasError {
		return err
	}
	return errs.Nil()
}

func addWorkMetadata(works []model.Work, metadata metadata.VolumeMetadata) errs.UploadError {
	for i := range works {
		work := &works[i]
		workMd := metadata.Works[i]

		work.Code = workMd.Code
		work.Siglum = workMd.Siglum
		if workMd.Year != nil {
			work.Year = *workMd.Year
		}
		if work.Year == "" {
			return errs.New(fmt.Errorf("the year for the work '%s' can neither be found in the XML data nor in the volume metadata", work.Title), nil)
		}
	}
	return errs.Nil()
}
