package main

import (
	"net/http"
	"os"
	"strings"
	"time"
	"tomb_mates/internal/game"
	"tomb_mates/internal/hub"
	"tomb_mates/internal/protos"
	"tomb_mates/web"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"
)

const tickRate = time.Second / 30

func main() {
	e := echo.New()
	w := game.New(false, make(map[uint32]*protos.Unit))
	h := hub.New(w)
	go w.Run(tickRate)

	// e.Use(middleware.Logger())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		LogLevel:          log.ERROR,
		DisableStackAll:   true,
		DisablePrintStack: true,
	}))

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

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "IndexPage", "HakaHata")
	})

	e.GET("/ws", h.WsHandler(w))

	e.Logger.Fatal(e.Start(":3000"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
