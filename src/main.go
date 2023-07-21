package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/FrHorschig/kant-search-backend/database/repository"
	"github.com/FrHorschig/kant-search-backend/handlers"
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

func registerHandlers(e *echo.Echo, workHandler handlers.WorkHandler, sectionHandler handlers.ParagraphHandler) {
	e.POST("/api/v1/works", func(ctx echo.Context) error {
		return workHandler.PostWork(ctx)
	})
	e.GET("/api/v1/works", func(ctx echo.Context) error {
		return workHandler.GetWork(ctx)
	})
	e.GET("/api/v1/work/:id/paragraphs", func(ctx echo.Context) error {
		return sectionHandler.GetParagraphs(ctx)
	})
}

func main() {
	db := initDbConnection()
	defer db.Close()
	workRepo := repository.NewWorkRepo(db)
	paragraphRepo := repository.NewParagraphRepo(db)
	sentenceRepo := repository.NewSentenceRepo(db)

	workHandler := handlers.NewWorkHandler(workRepo, paragraphRepo, sentenceRepo)
	paragraphHandler := handlers.NewParagraphHandler(paragraphRepo)

	e := initEchoServer()
	registerHandlers(e, workHandler, paragraphHandler)
	e.Logger.Fatal(e.StartTLS(":3000", "ssl/cert.pem", "ssl/key.pem"))
}
