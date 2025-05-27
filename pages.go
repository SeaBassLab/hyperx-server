package server

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	templates "github.com/SeaBassLab/hyperx-templates"
	"github.com/go-chi/chi/v5"
)

func LoadPages(r chi.Router, pagesDir string) {
	renderer := templates.NewRenderer(pagesDir)

	files, err := filepath.Glob(filepath.Join(pagesDir, "*.html"))
	if err != nil {
		panic("❌ Error leyendo templates: " + err.Error())
	}

	for _, file := range files {
		name := filepath.Base(file)
		route := "/" + strings.TrimSuffix(name, filepath.Ext(name))
		if route == "/index" {
			route = "/"
		}

		tmplName := name // se pasa el nombre exacto del template

		r.Get(route, func(w http.ResponseWriter, r *http.Request) {
			renderer.Render(w, tmplName, nil)
		})

		fmt.Println("📄 Page route registered:", route)
	}
}

func ServeStatic(r chi.Router, publicDir string) {
	fs := http.FileServer(http.Dir(publicDir))

	// Servir todos los archivos desde /public como si fueran en raíz
	r.Handle("/*", http.StripPrefix("/", fs))
}

