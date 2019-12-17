package main

import (
	"gophercises/ex18/primitive"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var r = mux.NewRouter()

func main() {
	// go func() {
	// 	t := time.NewTicker(5 * time.Minute)
	// 	for {
	// 		<-t.C
	// 	}
	// }()
	r.HandleFunc("/image_upload.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "."+r.URL.Path)
	})
	r.HandleFunc("/modify/{id}", primitive.Modify)
	r.HandleFunc("/upload", primitive.Upload).Methods("POST")
	fs := http.FileServer(http.Dir("./img/"))
	r.Handle("/img/{id}", http.StripPrefix("/img", fs))
	log.Fatal(http.ListenAndServe(":3000", r))
}
