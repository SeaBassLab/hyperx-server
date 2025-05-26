package server

import (
	"net/http"
	"sync"
)

// handlerRegistry guarda todas las rutas que se van registrando con Register()
var handlerRegistry = make(map[string]http.HandlerFunc)
var mu sync.Mutex

// Register permite a los devs registrar handlers de la API (lo usan en init())
func Register(path string, handler http.HandlerFunc) {
	mu.Lock()
	defer mu.Unlock()
	handlerRegistry[path] = handler
}

// GetRegisteredHandlers devuelve el mapa con todos los handlers
func GetRegisteredHandlers() map[string]http.HandlerFunc {
	mu.Lock()
	defer mu.Unlock()

	// Clonar el mapa para evitar data races
	cloned := make(map[string]http.HandlerFunc)
	for k, v := range handlerRegistry {
		cloned[k] = v
	}
	return cloned
}
