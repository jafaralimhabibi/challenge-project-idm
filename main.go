package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	config "challenge-project/config"
	localmw "challenge-project/middleware"
	rAuth "challenge-project/router/auth"
	rProduct "challenge-project/router/products"
	socketService "challenge-project/services/socket"

	echomw "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	if err := config.LoadConfig("config.yml"); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := config.GetConfig()
	ctx := context.Background()
	if _, err := config.InitMongo(ctx); err != nil {
		log.Fatalf("failed to init mongo: %v", err)
	}

	e := echo.New()
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())

	io := socketService.NewServer()
	e.Any("/socket.io/*", echo.WrapHandler(io))

	api := e.Group("/api")
	auth := e.Group("/auth")
	api.Use(localmw.ValidateMiddleware())
	auth.Use(localmw.ValidateMiddleware())

	rAuth.UserAuth(auth)
	rProduct.ProductRoute(api)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	go func() {
		log.Printf("server starting at %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := e.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("shutdown failed: %v", err)
	}
	io.Close()

	time.Sleep(3600 * time.Millisecond)
}
