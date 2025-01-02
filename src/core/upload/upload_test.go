//go:build unit
// +build unit

package upload

import (
	"context"
	"fmt"
	"testing"

	"github.com/frhorschig/kant-search-backend/common/errors"
	"github.com/frhorschig/kant-search-backend/common/model"
	"github.com/frhorschig/kant-search-backend/core/upload/internal/mocks"
	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWorkUploadProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockTextMapper := mocks.NewMockTextMapper(ctrl)
	mockParagraphRepo := dbMocks.NewMockParagraphRepo(ctrl)
	mockSentenceRepo := dbMocks.NewMockSentenceRepo(ctrl)
	processor := &workUploadProcessorImpl{
		textMapper:    mockTextMapper,
		paragraphRepo: mockParagraphRepo,
		sentenceRepo:  mockSentenceRepo,
	}

	ctx := context.Background()
	testErr := &errors.Error{
		Msg:    errors.UPLOAD_GO_ERR,
		Params: nil,
	}

	testCases := []struct {
		name      string
		workId    int32
		text      string
		err       *errors.Error
		mockCalls func()
	}{
		{
			name:   "Transform returns an error",
			text:   "test text",
			workId: 3,
			err:    testErr,
			mockCalls: func() {
				mockTextMapper.EXPECT().FindParagraphs(gomock.Any(), gomock.Any()).Return([]model.Paragraph{}, testErr)
			},
		},
		{
			name:   "delete sentences returns an error",
			text:   "test text",
			workId: 4,
			err: &errors.Error{
				Msg:    errors.UPLOAD_GO_ERR,
				Params: []string{"deleteSentences error"},
			},
			mockCalls: func() {
				mockTextMapper.EXPECT().FindParagraphs(gomock.Any(), gomock.Any()).Return([]model.Paragraph{}, nil)
				mockSentenceRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(fmt.Errorf("deleteSentences error"))
			},
		},
		{
			name:   "delete paragraphs returns an error",
			text:   "test text",
			workId: 4,
			err: &errors.Error{
				Msg:    errors.UPLOAD_GO_ERR,
				Params: []string{"deleteParagraphs error"},
			},
			mockCalls: func() {
				mockTextMapper.EXPECT().FindParagraphs(gomock.Any(), gomock.Any()).Return([]model.Paragraph{}, nil)
				mockSentenceRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(nil)
				mockParagraphRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(fmt.Errorf("deleteParagraphs error"))
			},
		},
		{
			name:   "persistParagraphs returns an error",
			text:   "test text",
			workId: 4,
			err: &errors.Error{
				Msg:    errors.UPLOAD_GO_ERR,
				Params: []string{"persistParagraphs error"},
			},
			mockCalls: func() {
				mockTextMapper.EXPECT().FindParagraphs(gomock.Any(), gomock.Any()).Return([]model.Paragraph{{}}, nil)
				mockParagraphRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(nil)
				mockSentenceRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(nil)
				mockParagraphRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(int32(0), fmt.Errorf("persistParagraphs error"))
			},
		},
		{
			name:   "FindSentences returns an error",
			text:   "test text",
			workId: 5,
			err:    testErr,
			mockCalls: func() {
				mockTextMapper.EXPECT().FindParagraphs(gomock.Any(), gomock.Any()).Return([]model.Paragraph{{}}, nil)
				mockSentenceRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(nil)
				mockParagraphRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(nil)
				mockParagraphRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(int32(1), nil)
				mockTextMapper.EXPECT().FindSentences(gomock.Any()).Return(nil, testErr)
			},
		},
		{
			name:   "persistSentences returns an error",
			text:   "test text",
			workId: 6,
			err: &errors.Error{
				Msg:    errors.UPLOAD_GO_ERR,
				Params: []string{"persistSentences error"},
			},
			mockCalls: func() {
				mockTextMapper.EXPECT().FindParagraphs(gomock.Any(), gomock.Any()).Return([]model.Paragraph{{}}, nil)
				mockSentenceRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(nil)
				mockParagraphRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(nil)
				mockParagraphRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(int32(1), nil)
				mockTextMapper.EXPECT().FindSentences(gomock.Any()).Return([]model.Sentence{{}}, nil)
				mockSentenceRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("persistSentences error"))
			},
		},
		{
			name:   "success",
			text:   "test text",
			workId: 6,
			mockCalls: func() {
				mockTextMapper.EXPECT().FindParagraphs(gomock.Any(), gomock.Any()).Return([]model.Paragraph{{}}, nil)
				mockTextMapper.EXPECT().FindSentences(gomock.Any()).Return([]model.Sentence{{}}, nil)
				mockSentenceRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(nil)
				mockParagraphRepo.EXPECT().DeleteByWorkId(gomock.Any(), gomock.Any()).Return(nil)
				mockParagraphRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(int32(1), nil)
				mockSentenceRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).Return([]int32{1}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockCalls()
			err := processor.Process(ctx, tc.workId, tc.text)
			assert.Equal(t, tc.err, err)
		})
	}
}
