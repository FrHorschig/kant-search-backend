package metadatamapping

import (
	"fmt"

	"github.com/frhorschig/kant-search-backend/common/errs"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/metadatamapping/metadatamapping/metadata"
)

func MapMetadata(volume *model.Volume, works []model.Work, metadata metadata.Metadata) errs.UploadError {
	md, mdErr := metadata.Read(volume.VolumeNumber)
	if mdErr != nil {
		return errs.New(nil, mdErr)
	}
	volume.Title = md.Title

	err := addWorkMetadata(volume.VolumeNumber, works, md)
	if err.HasError {
		return err
	}
	return errs.Nil()
}

func addWorkMetadata(volNr int32, works []model.Work, metadata metadata.VolumeMetadata) errs.UploadError {
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
	handleSpecialCases(volNr, works)
	return errs.Nil()
}

func handleSpecialCases(volNr int32, works []model.Work) {
	if volNr == 8 {
		// This work has a h1 heading with the text "Nachtrag" and a h2 sub heading. The h1 text is not the heading of a work of Kant, but a heading from the editor ("Nachtrag" = addendum), the h2 text is the work title. We fix this here.
		rezSilber := &works[27]
		rezSilber.Title = rezSilber.Sections[0].Heading.TocText
		rezSilber.Paragraphs = rezSilber.Sections[0].Paragraphs
		rezSilber.Sections = []model.Section{}

		// This work has a h1 heading with the text "Anhang" and a h2 sub heading. The h1 text is not the heading of a work of Kant, but a heading from the editor ("Anhang" = appendix), the h2 text is the work title. We fix this here.
		rezUlrich := &works[28]
		rezUlrich.Title = rezUlrich.Sections[0].Heading.TocText
		rezUlrich.Paragraphs = rezUlrich.Sections[0].Paragraphs
		rezUlrich.Sections = []model.Section{}
	}
}
