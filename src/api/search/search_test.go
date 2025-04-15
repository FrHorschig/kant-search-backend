//go:build unit
// +build unit

package search

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frhorschig/kant-search-api/src/go/models"
	"github.com/frhorschig/kant-search-backend/core/search/mocks"
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
		"Search bind error":          testSearchBindError,
		"Search empty search string": testSearchEmptySearchString,
		"Search empty workIds":       testSearchEmptyWorkIds,
		"Search syntax error":        testSearchSyntaxError,
		"Search database error":      testSearchDatabaseError,
		"Search no result":           testSearchNotFound,
		"Search success":             testSearchSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, searchProcessor)
		})
	}
}

func testSearchBindError(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.Volume{VolumeNumber: 1, Section: 1})
	if err != nil {
		t.Fatal(err)
	}
	// GIVEN
	req := httptest.NewRequest(echo.GET, "/api/v1/search", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	// WHEN
	sut.Search(ctx)
	// THEN
	assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	assertErrorResponse(t, res, string(models.BAD_REQUEST_EMPTY_WORKS_SELECTION))
}

func testSearchEmptySearchString(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{WorkIds: []string{"id1"}, SearchString: "\t \n"})
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

func testSearchEmptyWorkIds(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{WorkIds: []string{}, SearchString: "& test"})
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

func testSearchSyntaxError(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{WorkIds: []string{"id1"}, SearchString: "& test", Options: models.SearchOptions{Scope: models.SearchScope("PARAGRAPH")}})
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
	assertErrorResponse(t, res, string(models.BAD_REQUEST_VALIDATION_WRONG_STARTING_CHAR))
}

func testSearchDatabaseError(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{WorkIds: []string{"workId"}, SearchString: "test"})
	if err != nil {
		t.Fatal(err)
	}
	matches := []model.SearchResult{}
	err = errors.New("database error")
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	searchProcessor.EXPECT().Search(gomock.Any(), gomock.Any()).Return(matches, err)
	// WHEN
	sut.Search(ctx)
	// THEN
	assert.Equal(t, http.StatusInternalServerError, ctx.Response().Status)
	assertErrorResponse(t, res, "")
}

func testSearchNotFound(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{WorkIds: []string{"id1"}, SearchString: "test"})
	if err != nil {
		t.Fatal(err)
	}
	matches := []model.SearchResult{}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	searchProcessor.EXPECT().Search(gomock.Any(), gomock.Any()).Return(matches, nil)
	// WHEN
	sut.Search(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
}

func testSearchSuccess(t *testing.T, sut *searchHandlerImpl, searchProcessor *mocks.MockSearchProcessor) {
	body, err := json.Marshal(models.SearchCriteria{
		WorkIds:      []string{"workId"},
		SearchString: "string",
		Options:      models.SearchOptions{Scope: models.SearchScope("PARAGRAPH")}})
	if err != nil {
		t.Fatal(err)
	}
	matches := []model.SearchResult{{
		Snippet:   "Test",
		Pages:     []int32{1},
		ContentId: "contentId",
		WorkId:    "workId",
	}}
	// GIVEN
	req := httptest.NewRequest(echo.POST, "/api/v1/search", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, res)
	searchProcessor.EXPECT().Search(gomock.Any(), gomock.Any()).Return(matches, nil)
	// WHEN
	sut.Search(ctx)
	// THEN
	assert.Equal(t, http.StatusOK, ctx.Response().Status)
	assert.Contains(t, res.Body.String(), "workId")
	assert.Contains(t, res.Body.String(), "matches")
	assert.Contains(t, res.Body.String(), "snippet")
	assert.Contains(t, res.Body.String(), "pages")
	assert.Contains(t, res.Body.String(), "contentId")
}

func assertErrorResponse(t *testing.T, res *httptest.ResponseRecorder, expectedMsg string) {
	assert.Contains(t, res.Body.String(), "code")
	assert.Contains(t, res.Body.String(), "message")
	assert.Contains(t, res.Body.String(), expectedMsg)
}
