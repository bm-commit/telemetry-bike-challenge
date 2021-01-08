package static

import (
	"net/http"
	"path/filepath"
)

// ServeStaticHTML handler to serve html frontend
func ServeStaticHTML(publicDir string) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, filepath.Join(publicDir, "/frontend.html"))
		return
	}
}
