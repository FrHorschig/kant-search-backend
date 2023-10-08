//go:build unit
// +build unit

package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestParagraphHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	paragraphRepo := mocks.NewMockParagraphRepo(ctrl)
	sut := NewParagraphHandler(paragraphRepo).(*paragraphHandlerImpl)

	for scenario, fn := range map[string]func(t *testing.T, sut *paragraphHandlerImpl, paragraphRepo *mocks.MockParagraphRepo){
		"GetParagraphs parse workId error": testGetParagraphsParseWorkIdError,
		"GetParagraphs database error":     testGetParagraphsDatabaseError,
		"GetParagraphs no result":          testGetParagraphsNotFound,
		"GetParagraphs success":            testGetParagraphsSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, paragraphRepo)
		})
	}
}

func testGetParagraphsParseWorkIdError(t *testing.T, sut *paragraphHandlerImpl, paragraphRepo *mocks.MockParagraphRepo) {
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/x/paragraphs", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("x")
	// WHEN
	sut.GetParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testGetParagraphsDatabaseError(t *testing.T, sut *paragraphHandlerImpl, paragraphRepo *mocks.MockParagraphRepo) {
	paragraphs := []model.Paragraph{}
	err := fmt.Errorf("error")
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/1/paragraphs", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	paragraphRepo.EXPECT().SelectAll(gomock.Any(), gomock.Any()).Return(paragraphs, err)
	// WHEN
	sut.GetParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testGetParagraphsNotFound(t *testing.T, sut *paragraphHandlerImpl, paragraphRepo *mocks.MockParagraphRepo) {
	paragraphs := []model.Paragraph{}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/1/paragraphs", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	paragraphRepo.EXPECT().SelectAll(gomock.Any(), gomock.Any()).Return(paragraphs, nil)
	// WHEN
	sut.GetParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusNotFound, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testGetParagraphsSuccess(t *testing.T, sut *paragraphHandlerImpl, paragraphRepo *mocks.MockParagraphRepo) {
	paragraphs := []model.Paragraph{{
		Id:     1,
		Text:   "text",
		Pages:  []int32{1, 2, 3},
		WorkId: 1,
	}}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/1/paragraphs", nil)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	ctx.SetParamNames("workId")
	ctx.SetParamValues("1")
	paragraphRepo.EXPECT().SelectAll(gomock.Any(), gomock.Any()).Return(paragraphs, nil)
	// WHEN
	sut.GetParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "id")
	assert.Contains(t, res.Body.String(), "text")
	assert.Contains(t, res.Body.String(), "pages")
	assert.Contains(t, res.Body.String(), "workId")
}
