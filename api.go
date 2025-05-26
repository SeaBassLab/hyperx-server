package server

import (
	"strings"

	"github.com/go-chi/chi/v5"
)

func LoadAPI(r chi.Router) {
	for path, handler := range GetRegisteredHandlers() {
		// Convierte rutas :id en {id} para chi
		chiPath := convertColonToBraces(path)

		// Registro catch-all, el handler decide qué método aceptar
		r.HandleFunc(chiPath, handler)
	}
}

// convierte "/api/post/:id" → "/api/post/{id}"
func convertColonToBraces(path string) string {
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if strings.HasPrefix(seg, ":") {
			segments[i] = "{" + seg[1:] + "}"
		}
	}
	return strings.Join(segments, "/")
}
