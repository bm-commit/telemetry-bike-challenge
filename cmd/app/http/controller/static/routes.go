package static

import "net/http"

// RegisterRoute of static html file
func RegisterRoute(publicDir string) {
	http.HandleFunc("/", ServeStaticHTML(publicDir))
}
