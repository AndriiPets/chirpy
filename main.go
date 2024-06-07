package main

import (
	"log"
	"net/http"
	"strconv"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits ++
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits:" + " " + strconv.Itoa(cfg.fileserverHits)))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func (cfg *apiConfig) handleResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func main() {
	const filepathRoot = "."
	port := ":8080"
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleReady)
	mux.HandleFunc("GET /metrics", apiCfg.handleMetrics)
	mux.HandleFunc("POST /reset", apiCfg.handleResetMetrics)
	mux.Handle("GET /app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))) )

	server := &http.Server{
		Addr: port,
		Handler: mux,
	}

	log.Printf("Server listening on port:%s\n",port)
	log.Fatal(server.ListenAndServe())
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK!"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}