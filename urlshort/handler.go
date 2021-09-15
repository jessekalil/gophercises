package urlshort

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

type urlShort struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func httpRedirect(rw http.ResponseWriter, url string) {
	rw.Header().Set("Location", url)
	rw.WriteHeader(301)
}

func MapHandler(pathsoUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		reqPath := r.URL.Path

		for path, url := range pathsoUrls {
			if path == reqPath {
				httpRedirect(rw, url)
				return
			}
		}

		fallback.ServeHTTP(rw, r)
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsoUrls, err := parseYAML(yml)

	if err != nil {
		return nil, err
	}

	return MapHandler(pathsoUrls, fallback), nil
}

func parseYAML(data []byte) (map[string]string, error) {
	var sliceUrlshort []urlShort

	err := yaml.Unmarshal(data, &sliceUrlshort)
	if err != nil {
		return nil, err
	}

	pathsoUrls := buildMap(sliceUrlshort)

	return pathsoUrls, nil
}

func buildMap(urlShort []urlShort) map[string]string {
	pathsoUrls := make(map[string]string)

	for _, urlshort := range urlShort {
		pathsoUrls[urlshort.Path] = urlshort.Url
	}

	return pathsoUrls
}
