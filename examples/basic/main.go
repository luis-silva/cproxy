package main

import (
	"log"
	"net/http"
	"os"

	"github.com/luis-silva/cproxy/v2"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
	proxyAddress := "https://localhost:30443"
	proxyAuth := "Bearer sometoken"
	handler := cproxy.New(cproxy.Options.ProxyAddress(proxyAddress), cproxy.Options.Logger(logger), cproxy.Options.ProxyAuth(proxyAuth))
	log.Println("Listening on:", "*:2080")
	_ = http.ListenAndServe(":2080", handler)
}
