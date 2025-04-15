package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	apiread "github.com/frhorschig/kant-search-backend/api/read"
	apisearch "github.com/frhorschig/kant-search-backend/api/search"
	apiupload "github.com/frhorschig/kant-search-backend/api/upload"
	coreread "github.com/frhorschig/kant-search-backend/core/read"
	coresearch "github.com/frhorschig/kant-search-backend/core/search"
	coreupload "github.com/frhorschig/kant-search-backend/core/upload"
	db "github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initEsConnection() *elasticsearch.TypedClient {
	es, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{fmt.Sprintf(
			"%s:%s", os.Getenv("KSDB_URL"), os.Getenv("KSDB_PORT"),
		)},
		Username:               os.Getenv("KSDB_USER"),
		Password:               os.Getenv("KSDB_PWD"),
		CertificateFingerprint: os.Getenv("KSDB_CERT_HASH"),
	})
	if err != nil {
		panic(err)
	}
	return es
}

func initEchoServer() *echo.Echo {
	e := echo.New()
	if !(os.Getenv("KSGO_DISABLE_SSL") == "true") {
		e.Pre(middleware.HTTPSRedirect())
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("KSGO_ALLOW_ORIGINS"), ","),
		AllowMethods: []string{echo.GET, echo.POST},
		AllowHeaders: []string{"*"},
	}))
	return e
}

func registerHandlers(e *echo.Echo, uploadHandler apiupload.UploadHandler, readHandler apiread.ReadHandler, searchHandler apisearch.SearchHandler) {
	e.POST("/api/v1/upload", func(ctx echo.Context) error {
		return uploadHandler.PostVolume(ctx)
	})

	e.GET(("/api/v1/volumes"), func(ctx echo.Context) error {
		return readHandler.ReadVolumes(ctx)
	})
	e.GET(("/api/v1/works/"), func(ctx echo.Context) error {
		return readHandler.ReadWork(ctx)
	})
	e.GET(("/api/v1/works/:workId/footnotes"), func(ctx echo.Context) error {
		return readHandler.ReadFootnotes(ctx)
	})
	e.GET(("/api/v1/works/:workId/headings"), func(ctx echo.Context) error {
		return readHandler.ReadHeadings(ctx)
	})
	e.GET(("/api/v1/works/:workId/paragraphs"), func(ctx echo.Context) error {
		return readHandler.ReadParagraphs(ctx)
	})
	e.GET(("/api/v1/works/:workId/summaries"), func(ctx echo.Context) error {
		return readHandler.ReadSummaries(ctx)
	})
}

func main() {
	es := initEsConnection()

	volumeRepo := db.NewVolumeRepo(es)
	workRepo := db.NewWorkRepo(es)
	contentRepo := db.NewContentRepo(es)

	uploadProcessor := coreupload.NewUploadProcessor(volumeRepo, workRepo, contentRepo)
	readProcessor := coreread.NewReadProcessor(volumeRepo, workRepo, contentRepo)
	searchProcessor := coresearch.NewSearchProcessor(contentRepo)

	uploadHandler := apiupload.NewUploadHandler(uploadProcessor)
	readHandler := apiread.NewReadHandler(readProcessor)
	searchHandler := apisearch.NewSearchHandler(searchProcessor)

	e := initEchoServer()
	registerHandlers(e, uploadHandler, readHandler, searchHandler)
	if os.Getenv("KSGO_DISABLE_SSL") == "true" {
		e.Logger.Fatal(e.Start(":3000"))
	} else {
		e.Logger.Fatal(e.StartTLS(":3000", os.Getenv("KSGO_CERT_PATH"), os.Getenv("KSGO_KEY_PATH")))
	}
}
