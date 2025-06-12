//go:build unit
// +build unit

package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/core/search/errors"
	"github.com/frhorschig/kant-search-backend/core/search/mocks"
	"github.com/frhorschig/kant-search-backend/dataaccess/esmodel"
	"github.com/frhorschig/kant-search-backend/dataaccess/model"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSearchHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchProcessor := mocks.NewMockSearchProcessor(ctrl)
	sut := NewSearchHandler(searchProcessor).(*searchHandlerImpl)

	for scenario, fn := range map[string]func(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor){
		"Search empty search string": testSearchEmptySearchTerms,
		"Search empty workCodes":     testSearchEmptyWorkCodes,
		"Search database error":      testSearchDatabaseError,
		"Search no result":           testSearchNotFound,
		"Search success":             testSearchSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, searchProcessor)
		})
	}
}

func testSearchEmptySearchTerms(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{SearchTerms: "\t \n", Options: models.SearchOptions{WorkCodes: []string{"code"}}})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	// WHEN
	sut.Search(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res, string(models.BAD_REQUEST_EMPTY_SEARCH_TERMS))
}

func testSearchEmptyWorkCodes(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{SearchTerms: "& test", Options: models.SearchOptions{WorkCodes: []string{}}})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	// WHEN
	sut.Search(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res, string(models.BAD_REQUEST_EMPTY_WORKS_SELECTION))
}

func testSearchDatabaseError(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{SearchTerms: "test", Options: models.SearchOptions{WorkCodes: []string{"code"}}})
	if err != nil {
		t.Fatal(err)
	}
	matches := []model.SearchResult{}
	testErr := errors.New(nil, fmt.Errorf("database error"))
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	searchProcessor.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Return(matches, testErr)
	// WHEN
	sut.Search(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assertErrorResponse(t, res, "")
}

func testSearchNotFound(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{SearchTerms: "test", Options: models.SearchOptions{WorkCodes: []string{"code"}}})
	if err != nil {
		t.Fatal(err)
	}
	matches := []model.SearchResult{}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	searchProcessor.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Return(matches, errors.Nil())
	// WHEN
	sut.Search(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
}

func testSearchSuccess(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{
		SearchTerms: "string",
		Options: models.SearchOptions{
			Scope:     models.SearchScope("PARAGRAPH"),
			WorkCodes: []string{"workCode"},
		}})
	if err != nil {
		t.Fatal(err)
	}
	matches := []model.SearchResult{{
		HighlightText: "highlightText",
		FmtText:       "fmtText",
		PageByIndex:   []esmodel.IndexNumberPair{{I: 12, Num: 37}},
		LineByIndex:   []esmodel.IndexNumberPair{{I: 8, Num: 2481}},
		Ordinal:       1,
		WorkCode:      "workCode",
	}}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	searchProcessor.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).Return(matches, errors.Nil())
	// WHEN
	sut.Search(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "workCode")
	assert.Contains(t, res.Body.String(), "hits")
	assert.Contains(t, res.Body.String(), "highlightText")
	assert.Contains(t, res.Body.String(), "fmtText")
	assert.Contains(t, res.Body.String(), "pageByIndex")
	assert.Contains(t, res.Body.String(), "lineByIndex")
	assert.Contains(t, res.Body.String(), "1")
}

func assertErrorResponse(t *testing.T, res *httptest.ResponseRecorder, expectedMsg string) {
	assert.Contains(t, res.Body.String(), "code")
	assert.Contains(t, res.Body.String(), "message")
	assert.Contains(t, res.Body.String(), expectedMsg)
}
