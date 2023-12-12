package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/izaakdale/lib/router"
)

func main() {
	crt, err := tls.LoadX509KeyPair(os.Getenv("SERVER_CRT"), os.Getenv("SERVER_KEY"))
	if err != nil {
		panic(err)
	}

	srv := http.Server{
		Addr:      fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		Handler:   router.New(),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{crt}},
	}

	srv.ListenAndServeTLS("", "")
}
