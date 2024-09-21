package main

import (
	"net/http"
	"os"
	"strings"
	"tomb_mates/internal/engine"
	"tomb_mates/internal/hub"
	"tomb_mates/internal/protos"
	"tomb_mates/web"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

var world *engine.World

func init() {
	world = &engine.World{
		Replica: false,
		Units:   map[string]*protos.Unit{},
	}
}

func main() {
	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(60))))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(getEnv("AUTH_SECRET", "jdkljskldjslk")))))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "ws") // Change "metrics" for your own path
		},
	}))
	e.Renderer = web.UiTemplates

	e.Static("/static", "assets")
	e.Static("/dist", "./.dist")

	go world.Evolve()

	h := hub.NewHub()
	go h.Run()

	e.GET("/", func(c echo.Context) error {

		return c.Render(http.StatusOK, "IndexPage", "HakaHata")
	})

	e.GET("/ws", wsHandler(h, world))

	e.Logger.Fatal(e.Start(":3000"))
}

func wsHandler(h *hub.Hub, world *engine.World) echo.HandlerFunc {
	return func(c echo.Context) error {
		hub.ServeWs(h, world, c.Response().Writer, c.Request())
		return nil
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
