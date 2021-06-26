package main

import (
	"log"
	"net/http"
)

func StartHTTPServer() {
	http.Handle("/faker", f)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatalf("HTTP Server Start error: %s\n", err.Error())
	}
}
