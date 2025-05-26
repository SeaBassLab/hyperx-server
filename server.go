package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StartDev(port string) error {
	r := chi.NewRouter()

	LoadPages(r, "apps/playground/pages")
	LoadAPI(r)
	ServeStatic(r, "apps/playground/public")

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("🚀 HyperX running at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, r)
}

