package urlhandler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// YamlConfig is exported.
type YamlConfig struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func UrlHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		if redirect_to, ok := pathsToUrls[path]; ok {
			http.Redirect(response, request, redirect_to, http.StatusFound)
		}
		fallback.ServeHTTP(response, request)
	}

}

func ParseYaml(yamlFilename string) map[string]string {
	var yamlConfig []YamlConfig
	yamlFile, err := ioutil.ReadFile(yamlFilename)
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	yamlMap := make(map[string]string)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, val := range yamlConfig {
		yamlMap[val.Path] = val.Url
	}
	return yamlMap
}
