package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mrecachinas/dcserver/internal/config"

	"net/http"

	"github.com/mrecachinas/dcserver/internal/api"
)

// Run is the main entrypoint into the DCServer app.
// It handles parsing the command-line, setting up
// connections to MongoDB and RabbitMQ, and
// instantiates and runs the echo server.
func Run(cfg *config.Config) {
	dcapi, err := api.NewDCAPI(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer dcapi.DB.Disconnect(ctx)
	defer dcapi.AMQPClient.Close()
	defer dcapi.AMQPChannel.Close()

	e := SetupEchoServer(dcapi)

	// Run server
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	go func() {
		if err := e.Start(address); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// SetupEchoServer sets up the actual webserver and connects
// the routes to the route handler functions.
func SetupEchoServer(dcapi *api.Api) *echo.Echo {
	// Setup server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// If we're in debug mode, allow CORS
	if dcapi.Cfg.Debug {
		e.Use(middleware.CORS())
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.INFO)
	}

	// Setup routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/api/status", dcapi.GetAllStatus)
	e.GET("/api/status/:id", dcapi.GetStatus)
	e.GET("/api/tasks", dcapi.GetTasks)
	e.POST("/api/tasks/create", dcapi.CreateTask)
	e.POST("/api/tasks/:id/stop", dcapi.StopTask)
	e.GET("/ws", dcapi.UpdaterWebsocket)
	e.File("/", "public/index.html")
	return e
}
