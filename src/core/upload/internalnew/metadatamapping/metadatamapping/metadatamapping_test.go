package metadatamapping

import (
	"errors"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/util"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/common/testutil"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/metadatamapping/metadatamapping/metadata"
	"github.com/frhorschig/kant-search-backend/core/upload/internalnew/metadatamapping/metadatamapping/metadata/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMetadataMapping(t *testing.T) {
	testCases := []struct {
		name        string
		metadata    metadata.VolumeMetadata
		metadataErr error
		volume      model.Volume
		works       []model.Work
		expVolume   model.Volume
		expWorks    []model.Work
		expError    string
	}{
		{
			name: "map volume and work data",
			metadata: metadata.VolumeMetadata{
				VolumeNumber: 2,
				Title:        "vol title",
				Works: []metadata.WorkMetadata{
					{Code: "code 1", Siglum: util.StrPtr("siglum 1"), Year: util.StrPtr("5678")},
					{Code: "code 1", Siglum: util.StrPtr("siglum 1")},
				},
			},
			volume:    model.Volume{VolumeNumber: 2},
			works:     []model.Work{{}, {Year: "1234"}},
			expVolume: model.Volume{VolumeNumber: 2, Title: "vol title"},
			expWorks: []model.Work{
				{Code: "code 1", Siglum: util.StrPtr("siglum 1"), Year: "5678"},
				{Code: "code 1", Siglum: util.StrPtr("siglum 1"), Year: "1234"},
			},
		},
		{
			name: "no year error",
			metadata: metadata.VolumeMetadata{
				VolumeNumber: 2,
				Works:        []metadata.WorkMetadata{{Code: "code 1"}},
			},
			volume:   model.Volume{VolumeNumber: 2},
			works:    []model.Work{{}},
			expError: "the year",
		},
		{
			name:        "error from metadata",
			metadataErr: errors.New("test err"),
			expError:    "test err",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			metadata := mocks.NewMockMetadata(ctrl)
			metadata.EXPECT().Read(tc.volume.VolumeNumber).Return(tc.metadata, tc.metadataErr)

			err := MapMetadata(&tc.volume, tc.works, metadata)
			if tc.expError != "" {
				assert.True(t, err.HasError)
				if err.DomainError != nil {
					assert.Contains(t, err.DomainError.Error(), tc.expError)
				}
				if err.TechnicalError != nil {
					assert.Contains(t, err.TechnicalError.Error(), tc.expError)
				}
			} else {
				assert.False(t, err.HasError)
				assert.Equal(t, tc.expVolume.VolumeNumber, tc.volume.VolumeNumber)
				assert.Equal(t, tc.expVolume.Title, tc.volume.Title)
				testutil.AssertWorks(t, tc.expWorks, tc.works)
			}
		})
	}
}
