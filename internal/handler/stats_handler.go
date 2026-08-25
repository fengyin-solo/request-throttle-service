package handler

import (
	"net/http"
	"strconv"

	"ratelimiter/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
	mux.HandleFunc("GET /api/stats/top-clients", s.topClients)
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	result, err := s.svc.StatsOverview()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) topClients(w http.ResponseWriter, r *http.Request) {
	n := 10
	if v := r.URL.Query().Get("n"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			n = parsed
		}
	}
	result, err := s.svc.TopClients(n)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
