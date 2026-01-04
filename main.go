package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/jingxinwangdev/go-prject/internal/app"
	"github.com/jingxinwangdev/go-prject/internal/routes"
)

func main() {
	var port int
	// go run main.go -port 8080
	// os.getenv("PORT")
	flag.IntVar(&port, "port", 8080, "The port the server will listen on")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	app.Logger.Printf("Starting server on port %d", port)

	// 1. 先注册路由 replace with chi
	// http.HandleFunc("/health", app.HealthCheckHandler)
	router := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
