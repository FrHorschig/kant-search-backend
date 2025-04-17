package search

import (
	"testing"

	dbMocks "github.com/frhorschig/kant-search-backend/dataaccess/mocks"
	"github.com/golang/mock/gomock"
)

func TestSearchProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	contentRepo := dbMocks.NewMockContentRepo(ctrl)
	sut := NewSearchProcessor(contentRepo).(*searchProcessorImpl)

	for scenario, fn := range map[string]func(t *testing.T, sut *searchProcessorImpl, searchProcessor *dbMocks.MockContentRepo){
		"Search syntax error": testSearchSyntaxError,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t, sut, contentRepo)
		})
	}
}

func testSearchSyntaxError(t *testing.T, sut *searchProcessorImpl, contentRepo *dbMocks.MockContentRepo) {
	// body, err := json.Marshal(models.SearchCriteria{WorkIds: []string{"id1"}, SearchString: "& test", Options: models.SearchOptions{Scope: models.SearchScope("PARAGRAPH")}})
	//
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	// // GIVEN
	// req := httptest.NewRequest(echo.POST, "/api/v1/search", bytes.NewReader(body))
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// res := httptest.NewRecorder()
	// ctx := echo.New().NewContext(req, res)
	// // WHEN
	// sut.Search(ctx)
	// // THEN
	// assert.Equal(t, http.StatusBadRequest, ctx.Response().Status)
	// assertErrorResponse(t, res, string(models.BAD_REQUEST_VALIDATION_WRONG_STARTING_CHAR))
}
