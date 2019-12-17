package main

import (
	"flag"
	"fmt"
	"gophercises/ex2/urlhandler"
	"net/http"
)

var pathsToUrls = map[string]string{
	"/mongo-db": "https://www.mongodb.com",
	"/mux":      "https://github.com/gorilla/mux",
}

var yamlFilename = flag.String("yaml", "", "a csv file in the format of 'question,answer'")
var mux = http.NewServeMux()

func main() {
	flag.Parse()
	urlHandler := urlhandler.UrlHandler(pathsToUrls, mux)
	if *yamlFilename != "" {
		urlHandler = urlhandler.UrlHandler(urlhandler.ParseYaml(*yamlFilename), mux)
	}
	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(response, "Home Page")
	})
	http.ListenAndServe(":3000", urlHandler)
}
