package handler

import (
	"net/http"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/httpx"
)

func (s *Server) registerRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rules", s.createRule)
	mux.HandleFunc("GET /api/rules", s.listRules)
	mux.HandleFunc("GET /api/rules/{id}", s.getRule)
	mux.HandleFunc("PUT /api/rules/{id}", s.updateRule)
	mux.HandleFunc("DELETE /api/rules/{id}", s.deleteRule)
}

type createRuleRequest struct {
	Key       string `json:"key"`
	Algorithm string `json:"algorithm"`
	Limit     int    `json:"limit"`
	WindowSec int    `json:"window_sec"`
	Status    string `json:"status"`
}

func (s *Server) createRule(w http.ResponseWriter, r *http.Request) {
	var req createRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.CreateRule(model.Rule{Key: req.Key, Algorithm: req.Algorithm, Limit: req.Limit, WindowSec: req.WindowSec, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, result)
}

func (s *Server) listRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RuleFilter{
		Algorithm: r.URL.Query().Get("algorithm"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.GetRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type updateRuleRequest struct {
	Key       string `json:"key"`
	Algorithm string `json:"algorithm"`
	Limit     int    `json:"limit"`
	WindowSec int    `json:"window_sec"`
	Status    string `json:"status"`
}

func (s *Server) updateRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.UpdateRule(id, model.Rule{Key: req.Key, Algorithm: req.Algorithm, Limit: req.Limit, WindowSec: req.WindowSec, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
