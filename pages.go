package server

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func LoadPages(r chi.Router, pagesDir string) {
	files, _ := filepath.Glob(filepath.Join(pagesDir, "*.html"))

	for _, file := range files {
		name := filepath.Base(file)
		route := "/" + strings.TrimSuffix(name, filepath.Ext(name))
		if route == "/index" {
			route = "/"
		}

		page := file // clonar variable para el closure
		r.Get(route, func(w http.ResponseWriter, r *http.Request) {
			tmpl := template.Must(template.ParseFiles(page))
			_ = tmpl.Execute(w, nil)
		})
		fmt.Println("📄 Page route registered:", route)
	}
}

func ServeStatic(r chi.Router, publicDir string) {
	fs := http.FileServer(http.Dir(publicDir))

	// Servir archivos estáticos desde /static/*
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Opcional: servir raíz "/" si querés exponer archivos directo sin prefijo
	// r.Handle("/*", http.FileServer(http.Dir(publicDir)))
}
