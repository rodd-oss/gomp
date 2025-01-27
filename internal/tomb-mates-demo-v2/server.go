/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"gomp/internal/tomb-mates-demo-v2/game"
	"gomp/internal/tomb-mates-demo-v2/hub"
	"gomp/internal/tomb-mates-demo-v2/web"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echopprof "github.com/sevenNt/echo-pprof"
	"golang.org/x/time/rate"
)

const gameTickRate = time.Second / 60

func main() {
	runtime.SetCPUProfileRate(1000)
	g := game.New(false)
	go g.Run(gameTickRate)

	e := echo.New()
	h := hub.New(g)
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
		return c.Render(http.StatusOK, "IndexPage", "Index")
	})
	e.GET("/game", func(c echo.Context) error {
		return c.Render(http.StatusOK, "GamePage", "Game")
	})

	e.GET("/ws", h.WsHandler(g))

	echopprof.Wrap(e)

	e.Logger.Fatal(e.Start(":3000"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
