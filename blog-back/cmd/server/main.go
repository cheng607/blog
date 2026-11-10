package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"blog-back/internal/config"
	"blog-back/internal/router"
	"blog-back/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	cfg := config.Load()

	db, err := store.Open(cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := router.Setup(cfg, db)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("blog-back listening on %s", addr)
	server := &http.Server{Addr: addr, Handler: r, ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 15 * time.Second, WriteTimeout: 80 * time.Second, IdleTimeout: 60 * time.Second}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
