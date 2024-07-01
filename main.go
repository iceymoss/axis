package main

import (
	"context"
	"github.com/iceymoss/axis/framework"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	core := framework.NewCore()
	//core.Use(middleware.Test1(), middleware.Test2())
	//subjectApi := core.Group("/test")
	//subjectApi.Use(middleware.Test3())
	registerRouter(core)
	serve := &http.Server{
		Addr:    ":8000",
		Handler: core,
	}
	go func() {
		serve.ListenAndServe()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit //主协程堵塞，等待结束命令

	// 调用Server.Shutdown graceful结束
	if err := serve.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
