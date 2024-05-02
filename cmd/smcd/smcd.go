package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cebarks/smcd"

	"github.com/gin-gonic/gin"
)

func main() {
	smcd.WorkingDir = os.Getenv("SMCD_DIR")
	if smcd.WorkingDir == "" {
		log.Fatalf("Could not get working dir: %v", smcd.WorkingDir)
	} else {
		log.Default().Printf("Working dir: %v", smcd.WorkingDir)
	}

	servers := smcd.DiscoverServers()

	if len(servers) == 0 {
		log.Fatal("No servers detected. Exiting.")
	}

	log.Default().Printf("Found %v Servers:\n", len(servers))
	for _, s := range servers {
		log.Default().Printf("- %s\n", s.Folder)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":25542")

	log.Default().Println("Done. Exiting.")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func StartServers(servers []*smcd.Server) {
	log.Default().Println("Starting servers...")
	for _, server := range servers {
		server.Start()
	}
}

func StopServers(servers []*smcd.Server) {
	log.Default().Println("Stopping servers...")
	for _, server := range servers {
		server.Stop()
	}
}
