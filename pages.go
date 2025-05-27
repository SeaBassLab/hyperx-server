package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	templates "github.com/SeaBassLab/hyperx-templates"
	"github.com/go-chi/chi/v5"
)

var paramRegex = regexp.MustCompile(`\[[^\]]+\]`)

func LoadPages(r chi.Router, isProd bool) {
	renderer := templates.NewRenderer("pages", isProd)

	err := filepath.WalkDir("pages", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Ignorar directorios
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".html") {
			return nil
		}

		relPath, _ := filepath.Rel("pages", path)
		route := buildRouteFromFile(relPath)

		page := relPath // necesario para el closure
		r.Get(route, func(w http.ResponseWriter, r *http.Request) {
			params := map[string]string{}
			for _, segment := range strings.Split(route, "/") {
				if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
					key := strings.Trim(segment, "{}")
					params[key] = chi.URLParam(r, key)
				}
			}

			renderer.Render(w, page, params)
		})

		fmt.Println("📄 Dynamic route registered:", route)
		return nil
	})

	if err != nil {
		panic("❌ Error recorriendo templates: " + err.Error())
	}
}

func ServeStatic(r chi.Router) {
	r.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intentamos servir archivo estático
		filePath := filepath.Join("public", r.URL.Path)
		if _, err := os.Stat(filePath); err == nil {
			http.ServeFile(w, r, filePath)
			return
		}

		// Si no existe el archivo, devolvemos 404
		http.NotFound(w, r)
	}))
}

func buildRouteFromFile(path string) string {
	route := "/" + strings.TrimSuffix(path, filepath.Ext(path)) // remove .html
	// Reemplazar [param] por {param}
	route = strings.ReplaceAll(route, "\\", "/") // Windows fix
	route = replaceParams(route)

	if route == "/index" {
		return "/"
	}

	return route
}

func replaceParams(route string) string {
	return paramRegex.ReplaceAllStringFunc(route, func(match string) string {
		param := strings.Trim(match, "[]")
		return "{" + param + "}"
	})
}
