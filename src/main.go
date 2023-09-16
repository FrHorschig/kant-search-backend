package main

import (
	"os"
	"strings"

	"github.com/FrHorschig/kant-search-backend/api/handlers"
	"github.com/FrHorschig/kant-search-backend/core/upload"
	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/FrHorschig/kant-search-backend/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

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
	db := util.InitDbConnection()
	defer db.Close()

	workRepo := repository.NewWorkRepo(db)
	paragraphRepo := repository.NewParagraphRepo(db)
	sentenceRepo := repository.NewSentenceRepo(db)
	searchRepo := repository.NewSearchRepo(db)

	workProcessor := upload.NewWorkProcessor(workRepo, paragraphRepo, sentenceRepo)

	workHandler := handlers.NewWorkHandler(workProcessor, workRepo)
	paragraphHandler := handlers.NewParagraphHandler(paragraphRepo)
	searchHandler := handlers.NewSearchHandler(searchRepo)

	e := initEchoServer()
	registerHandlers(e, workHandler, paragraphHandler, searchHandler)
	e.Logger.Fatal(e.StartTLS(":3000", "ssl/cert.pem", "ssl/key.pem"))
}
