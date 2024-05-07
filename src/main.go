package main

import (
	"database/sql"
	"os"
	"strings"

	"github.com/frhorschig/kant-search-backend/api/handlers"
	"github.com/frhorschig/kant-search-backend/core/search"
	"github.com/frhorschig/kant-search-backend/core/upload"
	"github.com/frhorschig/kant-search-backend/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initDbConnection() *sql.DB {
	connStr := "host=" + os.Getenv("KSGO_DB_HOST") +
		" port=" + os.Getenv("KSDB_PORT") +
		" user=" + os.Getenv("KSDB_USER") +
		" password=" + os.Getenv("KSDB_PASSWORD") +
		" dbname=" + os.Getenv("KSDB_NAME") +
		" sslmode=" + os.Getenv("KSGO_DB_SSLMODE")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func initEchoServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("KSGO_ALLOW_ORIGINS"), ","),
		AllowMethods: []string{echo.GET, echo.POST},
		AllowHeaders: []string{"*"},
	}))
	return e
}

func registerHandlers(e *echo.Echo, workHandler handlers.WorkHandler, paragraphHandler handlers.ParagraphHandler, searchHandler handlers.SearchHandler) {
	e.GET("/api/v1/volumes", func(ctx echo.Context) error {
		return workHandler.GetVolumes(ctx)
	})
	e.GET("/api/v1/works", func(ctx echo.Context) error {
		return workHandler.GetWorks(ctx)
	})
	e.POST("/api/v1/works/:workId", func(ctx echo.Context) error {
		return workHandler.PostWork(ctx)
	})
	e.GET("/api/v1/works/:workId/paragraphs", func(ctx echo.Context) error {
		return paragraphHandler.GetParagraphs(ctx)
	})
	e.POST("/api/v1/search", func(ctx echo.Context) error {
		return searchHandler.Search(ctx)
	})
}

func main() {
	db := initDbConnection()
	defer db.Close()

	volumeRepo := database.NewVolumeRepo(db)
	workRepo := database.NewWorkRepo(db)
	paragraphRepo := database.NewParagraphRepo(db)
	sentenceRepo := database.NewSentenceRepo(db)

	uploadProcessor := upload.NewWorkProcessor(paragraphRepo, sentenceRepo)
	searchProcessor := search.NewSearchProcessor(paragraphRepo, sentenceRepo)

	workHandler := handlers.NewWorkHandler(volumeRepo, workRepo, uploadProcessor)
	paragraphHandler := handlers.NewParagraphHandler(paragraphRepo)
	searchHandler := handlers.NewSearchHandler(searchProcessor)

	e := initEchoServer()
	registerHandlers(e, workHandler, paragraphHandler, searchHandler)
	if os.Getenv("KSGO_DISABLE_SSL") == "true" {
		e.Logger.Fatal(e.Start(":3000"))
	} else {
		e.Logger.Fatal(e.StartTLS(":3000", os.Getenv("KSGO_CERT_PATH"), os.Getenv("KSGO_KEY_PATH")))
	}
}
