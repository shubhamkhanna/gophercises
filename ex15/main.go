package main

import (
	"gophercises/ex15/handler"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/", handler.SourceCodeHandler)
	mux.HandleFunc("/panic/", handler.PanicDemo)
	mux.HandleFunc("/", handler.Hello)
	log.Fatal(http.ListenAndServe(":3000", handler.DevMiddleware(mux)))
}
