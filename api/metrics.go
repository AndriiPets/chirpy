package api

import (
	"html/template"
	"net/http"
	"strconv"
)

func (cfg *Api) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits ++
		next.ServeHTTP(w,r)
	})
}

func (cfg *Api) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("metrics.html"))
	hits := strconv.Itoa(cfg.fileserverHits)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl.Execute(w, hits)
}

func (cfg *Api) HandleResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}