package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FrHorschig/kant-search-api/models"
	"github.com/FrHorschig/kant-search-backend/database/model"
	"github.com/FrHorschig/kant-search-backend/database/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSearchHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchRepo := mocks.NewMockSearchRepo(ctrl)
	sut := &searchHandlerImpl{
		searchRepo: searchRepo,
	}

	for scenario, fn := range map[string]func(t *testing.T, sut *searchHandlerImpl, searchRepo *mocks.MockSearchRepo){
		"SearchParagraphs bind error":        testSearchParagraphsBindError,
		"SearchParagraphs empty searchTerms": testSearchParagraphsEmptySearchTerms,
		"SearchParagraphs empty workIds":     testSearchParagraphsEmptyWorkIds,
		"SearchParagraphs database error":    testSearchParagraphsDatabaseError,
		"SearchParagraphs no result":         testSearchParagraphsNotFound,
		"SearchParagraphs success":           testSearchParagraphsSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, searchRepo)
		})
	}
}

func testSearchParagraphsBindError(t *testing.T, sut *searchHandlerImpl, searchRepo *mocks.MockSearchRepo) {
	// GIVEN
	body, err := json.Marshal(models.Volume{Id: 1, Title: "title", Section: 1})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/search/paragraphs", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	// WHEN
	sut.SearchParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testSearchParagraphsEmptySearchTerms(t *testing.T, sut *searchHandlerImpl, searchRepo *mocks.MockSearchRepo) {
	// GIVEN
	body, err := json.Marshal(models.SearchCriteria{SearchTerms: []string{}, WorkIds: []int32{1}, Scope: models.PARAGRAPH})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search/paragraphs", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	// WHEN
	sut.SearchParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testSearchParagraphsEmptyWorkIds(t *testing.T, sut *searchHandlerImpl, searchRepo *mocks.MockSearchRepo) {
	// GIVEN
	body, err := json.Marshal(models.SearchCriteria{SearchTerms: []string{"test"}, WorkIds: []int32{}, Scope: models.PARAGRAPH})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search/paragraphs", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	// WHEN
	sut.SearchParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testSearchParagraphsDatabaseError(t *testing.T, sut *searchHandlerImpl, searchRepo *mocks.MockSearchRepo) {
	body, err := json.Marshal(models.SearchCriteria{SearchTerms: []string{"test"}, WorkIds: []int32{1}, Scope: models.PARAGRAPH})
	if err != nil {
		t.Fatal(err)
	}
	matches := []model.SearchResult{}
	err = errors.New("database error")
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search/paragraphs", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	searchRepo.EXPECT().SearchParagraphs(gomock.Any(), gomock.Any()).Return(matches, err)
	// WHEN
	sut.SearchParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testSearchParagraphsNotFound(t *testing.T, sut *searchHandlerImpl, searchRepo *mocks.MockSearchRepo) {
	body, err := json.Marshal(models.SearchCriteria{SearchTerms: []string{"test"}, WorkIds: []int32{1}, Scope: models.PARAGRAPH})
	if err != nil {
		t.Fatal(err)
	}
	matches := []model.SearchResult{}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search/paragraphs", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	searchRepo.EXPECT().SearchParagraphs(gomock.Any(), gomock.Any()).Return(matches, nil)
	// WHEN
	sut.SearchParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusNotFound, ctx.Response().Status)
	assertErrorResponse(t, res)
}

func testSearchParagraphsSuccess(t *testing.T, sut *searchHandlerImpl, searchRepo *mocks.MockSearchRepo) {
	body, err := json.Marshal(models.SearchCriteria{SearchTerms: []string{"string"}, WorkIds: []int32{1}, Scope: models.PARAGRAPH})
	if err != nil {
		t.Fatal(err)
	}
	matches := []model.SearchResult{{
		Volume:    1,
		WorkTitle: "Test",
		Snippet:   "Test",
		Pages:     []int32{1},
		WorkId:    1,
		ElementId: 1,
	}}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search/paragraphs", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	searchRepo.EXPECT().SearchParagraphs(gomock.Any(), gomock.Any()).Return(matches, nil)
	// WHEN
	sut.SearchParagraphs(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "workId")
	assert.Contains(t, res.Body.String(), "matches")
	assert.Contains(t, res.Body.String(), "snippet")
	assert.Contains(t, res.Body.String(), "pages")
	assert.Contains(t, res.Body.String(), "elementId")
}
