package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/alexyslozada/inbeat/backend"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS(), middleware.Logger())

	e.Static("/", "frontend")
	e.GET("/api/v1/influencer/:user_name", backend.Influencer)

	err := e.Start(":8080")
	if err != nil {
		log.Printf("error server: %v\n", err)
	}
}
