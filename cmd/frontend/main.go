package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	staticDir := getenvDefault("FRONTEND_DIR", "simple-frontend")
	addr := ":" + getenvDefault("FRONTEND_PORT", "3000")
	fs := http.FileServer(http.Dir(staticDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// serve static files, fallback to index.html for SPA-like routing
		p := filepath.Join(staticDir, filepath.Clean(r.URL.Path))
		if info, err := os.Stat(p); err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	// basic CORS passthrough for dev if requested
	corsOrigins := os.Getenv("FRONTEND_CORS_ALLOW")
	if corsOrigins != "" {
		origins := strings.Split(corsOrigins, ",")
		http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, o := range origins {
				if origin == strings.TrimSpace(o) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
					break
				}
			}
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		})
	}

	log.Printf("frontend server listening on %s, serving %s", addr, staticDir)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getenvDefault(k, d string) string {
	v := os.Getenv(k)
	if v == "" { return d }
	return v
}