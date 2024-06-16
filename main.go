package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	port := ":8080"
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", handleReady)
	mux.HandleFunc("GET /api/chirps", getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", getChirp)
	mux.HandleFunc("POST /api/chirps", postChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)
	mux.HandleFunc("GET /admin/reset", apiCfg.handleResetMetrics)
	mux.Handle("GET /app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))) )

	server := &http.Server{
		Addr: port,
		Handler: mux,
	}

	log.Printf("Server listening on port:%s\n",port)
	log.Fatal(server.ListenAndServe())
}

