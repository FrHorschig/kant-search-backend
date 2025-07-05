package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/healthstatus"
	apiread "github.com/frhorschig/kant-search-backend/api/read"
	apisearch "github.com/frhorschig/kant-search-backend/api/search"
	apiupload "github.com/frhorschig/kant-search-backend/api/upload"
	coreread "github.com/frhorschig/kant-search-backend/core/read"
	coresearch "github.com/frhorschig/kant-search-backend/core/search"
	coreupload "github.com/frhorschig/kant-search-backend/core/upload"
	db "github.com/frhorschig/kant-search-backend/dataaccess"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func initEsConnection() *elasticsearch.TypedClient {
	esCert, err := os.ReadFile(os.Getenv("KSDB_CERT"))
	if err != nil {
		panic(err)
	}
	es, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{fmt.Sprintf(
			"%s:%s", os.Getenv("KSDB_URL"), os.Getenv("KSDB_PORT"),
		)},
		Username: os.Getenv("KSDB_USERNAME"),
		Password: os.Getenv("KSDB_PASSWORD"),
		CACert:   esCert,
	})
	if err != nil {
		panic(err)
	}

	retryCount := readIntConfig("KSGO_RETRY_COUNT")
	retryInterval := readIntConfig("KSGO_RETRY_INTERVAL")
	for i := 0; i < retryCount; i++ {
		health, err := es.Cluster.Health().Do(context.Background())
		if err == nil && health.Status == healthstatus.Green {
			return es
		}
		log.Info().Msgf("waiting for ES to start after attempt %d", i)
		time.Sleep(time.Duration(retryInterval) * time.Second)
	}
	panic("failed to connect to Elasticsearch after maximum number of attempts")
}

func readIntConfig(name string) int {
	str := strings.TrimSpace(os.Getenv(name))
	if str == "" {
		panic("unknown environment variable " + name)
	}
	num, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		panic("unable to convert value '" + str + "' of environment variable '" + name + "' to a number")
	}
	return int(num)
}

func initEchoServer() *echo.Echo {
	e := echo.New()
	if os.Getenv("KSGO_DISABLE_SSL") != "true" {
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
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "UP")
	})
	e.POST("/api/v1/upload", func(ctx echo.Context) error {
		return uploadHandler.PostVolume(ctx)
	})

	e.GET(("/api/v1/volumes"), func(ctx echo.Context) error {
		return readHandler.ReadVolumes(ctx)
	})
	e.GET(("/api/v1/works/:workCode/footnotes"), func(ctx echo.Context) error {
		return readHandler.ReadFootnotes(ctx)
	})
	e.GET(("/api/v1/works/:workCode/headings"), func(ctx echo.Context) error {
		return readHandler.ReadHeadings(ctx)
	})
	e.GET(("/api/v1/works/:workCode/paragraphs"), func(ctx echo.Context) error {
		return readHandler.ReadParagraphs(ctx)
	})
	e.GET(("/api/v1/works/:workCode/summaries"), func(ctx echo.Context) error {
		return readHandler.ReadSummaries(ctx)
	})

	e.POST(("/api/v1/search"), func(ctx echo.Context) error {
		return searchHandler.Search(ctx)
	})
}

func main() {
	es := initEsConnection()

	volumeRepo := db.NewVolumeRepo(es)
	contentRepo := db.NewContentRepo(es)

	uploadProcessor := coreupload.NewUploadProcessor(volumeRepo, contentRepo, os.Getenv("KSGO_CONFIG_PATH"))
	readProcessor := coreread.NewReadProcessor(volumeRepo, contentRepo)
	searchProcessor := coresearch.NewSearchProcessor(contentRepo)

	uploadHandler := apiupload.NewUploadHandler(uploadProcessor)
	readHandler := apiread.NewReadHandler(readProcessor)
	searchHandler := apisearch.NewSearchHandler(searchProcessor)

	e := initEchoServer()
	registerHandlers(e, uploadHandler, readHandler, searchHandler)
	if os.Getenv("KSGO_DISABLE_SSL") == "true" {
		e.Logger.Fatal(e.Start(":3000"))
	} else {
		e.Logger.Fatal(e.StartTLS(":443", os.Getenv("KSGO_CERT"), os.Getenv("KSGO_KEY")))
	}
}
