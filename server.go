package main

import (
	"github.com/Tamrin007/codecheck-3535/ws"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Static("/", "app")

	hub := ws.NewHub()
	go hub.Run()

	e.GET("/chat", echo.HandlerFunc(func(c echo.Context) error {
		return ws.Handler(hub, c)
	}))
	e.Logger.Fatal(e.Start(":1323"))
}
