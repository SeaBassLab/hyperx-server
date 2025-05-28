package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/websocket"
)

var reloadClients = make(map[*websocket.Conn]bool)

func StartServer(env string, port string) error {
	r := chi.NewRouter()

	if env == "prod" {
		fmt.Println("🌐 Running in production mode")
		// TODO: Agregar middlewares útiles para prod
	} else {
		fmt.Println("🔧 Running in development mode")
	}
	if env != "prod" {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				w.Header().Set("Pragma", "no-cache")
				w.Header().Set("Expires", "0")
				next.ServeHTTP(w, r)
			})
		})
	}
	// Carga de páginas, APIs y archivos estáticos
	LoadPages(r, env == "prod")
	LoadAPI(r)
	ServeStatic(r)

		r.Handle("/__reload", websocket.Handler(func(ws *websocket.Conn) {
		reloadClients[ws] = true
		defer ws.Close()
		for {
			time.Sleep(1 * time.Hour) // idle
		}
	}))

	go watchAndReload([]string{"pages", "public"})

	// Dirección completa
	addr := fmt.Sprintf(":%s", port)

	// Mostrar URL amigable
	host := "localhost"
	if env == "prod" {
		host = "0.0.0.0"
	}
	fmt.Printf("🚀 HyperX running at http://%s%s\n", host, addr)

	return http.ListenAndServe(addr, r)
}

func watchAndReload(dirs []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, dir := range dirs {
		if err := watcher.Add(dir); err != nil {
			log.Println("Error viendo", dir, ":", err)
		}
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Println("🌀 File changed:", event.Name)
			for ws := range reloadClients {
				_ = websocket.Message.Send(ws, "reload")
			}
		case err := <-watcher.Errors:
			log.Println("Watcher error:", err)
		}
	}
}

