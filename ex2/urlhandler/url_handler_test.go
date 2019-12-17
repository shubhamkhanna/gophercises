package urlhandler

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"
)

var pathsToUrls = map[string]string{
	"/mongo-db": "https://www.mongodb.com",
	"/mux":      "https://github.com/gorilla/mux",
}

var invalidPathsToUrls = map[string]string{}

var yamlFilename = flag.String("yaml", "test_path_urls.yml", "a csv file in the format of 'question,answer'")
var invalidyamlFilename = flag.String("invalid_yaml", "invalid_path_urls.yml", "a csv file in the format of 'question,answer'")
var mux = http.NewServeMux()

func TestUrlHandler(t *testing.T) {
	t.Run("it checks valid scenario", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/mongo-db", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := UrlHandler(pathsToUrls, mux)
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusFound {
			t.Errorf("Record not found")
		}
	})

	t.Run("it checks invalid scenario", func(t *testing.T) {
		req, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		invalidHandler := UrlHandler(invalidPathsToUrls, mux)
		invalidHandler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusMovedPermanently {
			t.Errorf("In valid scenario")
		}
	})
}

func TestParseYaml(t *testing.T) {
	t.Run("it returns parsed YML map", func(t *testing.T) {
		parsedMap := ParseYaml(*yamlFilename)
		if parsedMap == nil {
			t.Errorf("Map is blank!")
		}
	})

	t.Run("it returns an error", func(t *testing.T) {
		parsedMap := ParseYaml(*invalidyamlFilename)
		if parsedMap != nil {
			t.Errorf("Error should be  occurred")
		}
	})
}
