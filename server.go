package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StartDev(port string) error {
	r := chi.NewRouter()

	LoadPages(r, "pages")
	LoadAPI(r)
	ServeStatic(r, "public")

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("🚀 HyperX running at http://localhost%s\n", addr)
	return http.ListenAndServe(addr, r)
}

