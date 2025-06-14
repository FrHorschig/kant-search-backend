package metadata

//go:generate mockgen -source=$GOFILE -destination=mocks/metadata_mock.go -package=mocks

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// TODO keep in sync with <frontend-dir>/src/app/search/util/works-select-util.ts
type Metadata interface {
	Read(volNr int32) (VolumeMetadata, error)
}

type metadataImpl struct {
	metadataPath string
}

func NewMetadata(configPath string) Metadata {
	impl := metadataImpl{
		metadataPath: configPath + "/volume-metadata.json",
	}
	return &impl
}

type VolumeMetadata struct {
	VolumeNumber int32          `json:"volumeNumber"`
	Title        string         `json:"title"`
	Works        []WorkMetadata `json:"works"`
}

type WorkMetadata struct {
	Code   string  `json:"code"`
	Siglum *string `json:"siglum,omitempty"`
	Year   *string `json:"year,omitempty"`
}

func (rec *metadataImpl) Read(volNr int32) (VolumeMetadata, error) {
	file, err := os.Open(rec.metadataPath)
	if err != nil {
		return VolumeMetadata{}, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return VolumeMetadata{}, fmt.Errorf("error reading file: %w", err)
	}

	var volumes []VolumeMetadata
	err = json.Unmarshal(bytes, &volumes)
	if err != nil {
		return VolumeMetadata{}, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	// we can just use volNr-1 as index here, because from api/upload/upload.go we know that the number is between 1 and 9
	return volumes[volNr-1], nil
}
