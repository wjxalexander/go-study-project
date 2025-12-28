package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/jingxinwangdev/go-prject/internal/app"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "The port the server will listen on")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	app.Logger.Printf("Starting server on port %d", port)

	// 1. 先注册路由
	http.HandleFunc("/health", HealthCheckHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}

// http.ResponseWriter 实际上是一个接口（Interface），而不是一个结构体（Struct）。
// 接口变量：http.ResponseWriter 就像是一个盒子，里面已经装着那个指针了。
// 3. 调用方法会变得很麻烦 接口的设计初衷是让你直接调用它的方法。

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
