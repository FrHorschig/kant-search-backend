package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

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

func registerTextHandlers(e *echo.Echo, debugHandler handlers.DebugHandler, addWorkHandler handlers.AddWorkHandler) {
	e.GET("/api/v1/debug", func(ctx echo.Context) error {
		return debugHandler.GetDebugInfo(ctx)
	})
	e.POST("/api/v1/upload/work", func(ctx echo.Context) error {
		return addWorkHandler.PostWork(ctx)
	})
}

func main() {
	db := initDbConnection()
	defer db.Close()

	debugHandler := handlers.NewDebugHandler()
	addWorkHandler := handlers.NewAddWorkHandler()

	e := initEchoServer()
	registerTextHandlers(e, debugHandler, addWorkHandler)
	e.Logger.Fatal(e.StartTLS(":3000", "ssl/cert.pem", "ssl/key.pem"))
}
