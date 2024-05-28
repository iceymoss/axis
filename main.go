package main

import (
	"github.com/iceymoss/axis/framework"
	"net/http"
)

func main() {
	serve := http.Server{
		Addr:    "8080",
		Handler: framework.NewCore(),
	}
	serve.ListenAndServe()
}
