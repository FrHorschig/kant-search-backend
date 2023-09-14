package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/FrHorschig/kant-search-backend/api/handlers"
	"github.com/FrHorschig/kant-search-backend/core/read"
	"github.com/FrHorschig/kant-search-backend/core/search"
	"github.com/FrHorschig/kant-search-backend/core/upload"
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func initDbConnection() *sql.DB {
	connStr := "user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func initEchoServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(os.Getenv("ALLOW_ORIGINS"), ","),
		AllowMethods: []string{echo.GET},
	}))
	return e
}

func registerHandlers(e *echo.Echo, workHandler handlers.WorkHandler, sectionHandler handlers.ParagraphHandler, searchHander handlers.SearchHandler) {
	e.GET("/api/v1/volumes", func(ctx echo.Context) error {
		return workHandler.GetVolumes(ctx)
	})
	e.GET("/api/v1/works", func(ctx echo.Context) error {
		return workHandler.GetWorks(ctx)
	})
	e.POST("/api/v1/works", func(ctx echo.Context) error {
		return workHandler.PostWork(ctx)
	})
	e.GET("/api/v1/works/:workId/paragraphs/:paragraphId", func(ctx echo.Context) error {
		return sectionHandler.GetParagraph(ctx)
	})
	e.GET("/api/v1/works/:id/paragraphs", func(ctx echo.Context) error {
		return sectionHandler.GetParagraphs(ctx)
	})
	e.POST("/api/v1/search/paragraphs", func(ctx echo.Context) error {
		return searchHander.SearchParagraphs(ctx)
	})
}

func main() {
	db := initDbConnection()
	defer db.Close()

	workRepo := repository.NewWorkRepo(db)
	paragraphRepo := repository.NewParagraphRepo(db)
	sentenceRepo := repository.NewSentenceRepo(db)
	searchRepo := repository.NewSearchRepo(db)

	workProcessor := upload.NewWorkProcessor(workRepo, paragraphRepo, sentenceRepo)
	workReader := read.NewWorkReader(workRepo)
	paragraphReader := read.NewParagraphReader(paragraphRepo)
	searcher := search.NewSearcher(searchRepo)

	workHandler := handlers.NewWorkHandler(workProcessor, workReader)
	paragraphHandler := handlers.NewParagraphHandler(paragraphReader)
	searchHandler := handlers.NewSearchHandler(searcher)

	e := initEchoServer()
	registerHandlers(e, workHandler, paragraphHandler, searchHandler)
	e.Logger.Fatal(e.StartTLS(":3000", "ssl/cert.pem", "ssl/key.pem"))
}
