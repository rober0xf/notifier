package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Serve(router *gin.Engine) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", 3000),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on port: %d", 3000)

	return server.ListenAndServe()
}
