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

func registerTextHandlers(e *echo.Echo, handler handlers.TextHandler) {
	e.GET("/api/v1/text/:id", func(ctx echo.Context) error {
		return handler.GetTextById(ctx)
	})
}

func main() {
	db := initDbConnection()
	defer db.Close()
	textRepo := repository.NewTextRepo(db)

	textHandler := handlers.NewTextHandler(textRepo)

	e := initEchoServer()
	registerTextHandlers(e, textHandler)
	e.Logger.Fatal(e.StartTLS(":3000", "ssl/cert.pem", "ssl/key.pem"))
}
