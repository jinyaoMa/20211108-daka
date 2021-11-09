package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var server *http.Server

var IsRunning = false
var HasSwagger = false

func Start(withSwagger bool) {
	HasSwagger = withSwagger

	router = gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	router.StaticFile("/", "./client/dist/index.html")
	router.Static("/css", "./client/dist/css")
	router.Static("/fonts", "./client/dist/fonts")
	router.Static("/js", "./client/dist/js")

	Apis(router, HasSwagger)

	server = &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func Run(withSwagger bool) {
	if IsRunning {
		return
	}
	IsRunning = true
	if withSwagger {
		log.Println("Server With Swagger Start!")
	} else {
		log.Println("Server Start!")
	}

	go Start(withSwagger)
}

func Stop() {
	if !IsRunning {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server Stop!")
	IsRunning = false
}
