package pkg

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sandbox.go/jwt/internal/handler"
	"sandbox.go/jwt/internal/middleware"
	"syscall"
	"time"
)

func server(port int) *http.Server {
	router := gin.Default()
	router.POST("/token", handler.CreateTokenHandler)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	api := router.Group("/v1")
	api.Use(middleware.AuthorizeJWT)
	api.GET("/auth", func(ctx *gin.Context) {
		name, _ := ctx.Get("userName")
		role, _ := ctx.Get("userRole")
		data := gin.H{"name": name, "role": role}
		ctx.JSON(http.StatusOK, data)
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

}

func StartApp() {
	port := 8080
	log.Printf("Starting server at Port %v", port)
	srv := server(port)
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv.RegisterOnShutdown(func() {
		log.Println("Shutting down server")
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %srv\n", err)
		}
	}()
	log.Printf("Server Started at %d", port)
	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
