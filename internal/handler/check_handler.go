package handler

import (
	"net/http"

	"ratelimiter/pkg/httpx"
)

func (s *Server) registerCheckRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/check", s.check)
}

type checkRequest struct {
	ClientID string `json:"client_id"`
	RuleKey  string `json:"rule_key"`
	Path     string `json:"path"`
}

func (s *Server) check(w http.ResponseWriter, r *http.Request) {
	var req checkRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.Check(req.ClientID, req.RuleKey, req.Path)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
