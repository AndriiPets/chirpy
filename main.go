package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/AndriiPets/chirpy/api"
)

func main() {
	const filepathRoot = "."
	port := ":8080"
	apiCfg := apiConfig{}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		_, err := os.Open("database.json")
		if err == nil {
			os.Remove("database.json")
		}

	}

	api := api.NewApi()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", api.HandleReady)
	mux.HandleFunc("GET /api/chirps", api.GetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", api.GetChirp)
	mux.HandleFunc("POST /api/chirps", api.PostChirp)
	mux.HandleFunc("POST /api/users", api.PostUser)
	mux.HandleFunc("POST /api/login", api.LoginUser)
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

