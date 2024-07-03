package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/AndriiPets/chirpy/api"
	"github.com/joho/godotenv"
)

func main() {
	
	godotenv.Load()
	const filepathRoot = "."
	port := ":8080"


	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		_, err := os.Open("database.json")
		if err == nil {
			os.Remove("database.json")
		}

	}

	api := api.NewApi(os.Getenv("JWT_SECRET"),)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", api.HandleReady)

	mux.HandleFunc("GET /api/chirps", api.GetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", api.GetChirp)
	mux.HandleFunc("POST /api/chirps", api.PostChirp)

	mux.HandleFunc("POST /api/users", api.PostUser)
	mux.HandleFunc("PUT /api/users", api.UpdateUser)
	mux.HandleFunc("POST /api/login", api.LoginUser)
	
	mux.HandleFunc("GET /admin/metrics", api.HandleMetrics)
	mux.HandleFunc("GET /admin/reset", api.HandleResetMetrics)
	mux.Handle("GET /app/*", api.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))) )

	server := &http.Server{
		Addr: port,
		Handler: mux,
	}

	log.Printf("Server listening on port:%s\n",port)
	log.Fatal(server.ListenAndServe())
}

