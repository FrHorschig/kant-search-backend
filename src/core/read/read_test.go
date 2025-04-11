//go:build unit
// +build unit

package read

import (
	"context"
	"testing"

	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestReadProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	volumeRepo := dbMocks.NewMockVolumeRepo(ctrl)
	workRepo := dbMocks.NewMockWorkRepo(ctrl)
	contentRepo := dbMocks.NewMockContentRepo(ctrl)
	sut := &readProcessorImpl{
		volumeRepo:  volumeRepo,
		workRepo:    workRepo,
		contentRepo: contentRepo,
	}

	ctx := context.Background()

	// TODO implement me
	testCases := []struct {
		name      string
		err       error
		mockCalls func()
		assert    func(t *testing.T)
	}{
		{
			name: "Processing is successful",
			mockCalls: func() {
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()
			err := sut.Process(ctx)
			assert.Equal(t, tc.err, err)
		})
	}
}
