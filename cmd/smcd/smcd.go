package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cebarks/smcd"

	"github.com/gin-gonic/gin"
)

const HTTP_SHUTDOWN_TIMEOUT = 30

func main() {
	smcd.WorkingDir = os.Getenv("SMCD_DIR")
	if smcd.WorkingDir == "" {
		log.Fatalf("Could not get working dir: %v", smcd.WorkingDir)
	} else {
		log.Printf("Working dir: %v", smcd.WorkingDir)
	}

	servers := smcd.DiscoverServers()

	if len(servers) == 0 {
		log.Fatal("No servers detected. Exiting.")
	}

	log.Printf("Found %v Servers:\n", len(servers))
	for _, s := range servers {
		log.Printf("- %s\n", s.Folder)
	}

	router := setupRouter(servers)

	srv := &http.Server{
		Addr:    ":25542",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	StartServers(servers)

	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	StopServers(servers)

	ctx, cancel := context.WithTimeout(context.Background(), HTTP_SHUTDOWN_TIMEOUT*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting.")
}

func setupRouter(servers []*smcd.Server) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/servers", func(ctx *gin.Context) {
		ctx.JSON(200, servers)
	})

	return r
}

func StartServers(servers []*smcd.Server) {
	log.Println("Starting servers...")
	for _, server := range servers {
		server.Start()
	}
}

func StopServers(servers []*smcd.Server) {
	log.Println("Stopping servers...")
	for _, server := range servers {
		server.Stop()
	}
}
