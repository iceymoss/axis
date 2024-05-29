package main

import (
	"github.com/iceymoss/axis/framework"
	"net/http"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	serve := &http.Server{
		Addr:    ":8000",
		Handler: core,
	}
	serve.ListenAndServe()
}
