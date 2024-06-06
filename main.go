package main

import (
	"github.com/iceymoss/axis/framework"
	"net/http"
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
	serve.ListenAndServe()
}
