package main

import (
	"log"
	"net/http"

	"github.com/fengjx/xin"
	"github.com/fengjx/xin/middleware"
)

func main() {
	app := xin.New()
	app.Use(middleware.Logger)
	app.GET("/", func(w http.ResponseWriter, r *http.Request) {
		xin.WriteString(w, http.StatusOK, "Hello World!")
	})
	log.Println("Server starting on :8080...")
	app.Run(":8080")
}
