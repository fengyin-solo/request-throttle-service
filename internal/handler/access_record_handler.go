package handler

import (
	"net/http"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/httpx"
)

func (s *Server) registerAccessRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/access-records", s.createAccessRecord)
	mux.HandleFunc("GET /api/access-records", s.listAccessRecords)
	mux.HandleFunc("GET /api/access-records/{id}", s.getAccessRecord)
}

type createAccessRecordRequest struct {
	ClientID  string `json:"client_id"`
	RuleID    string `json:"rule_id"`
	Path      string `json:"path"`
	Allowed   bool   `json:"allowed"`
	Remaining int    `json:"remaining"`
	LatencyMs int    `json:"latency_ms"`
}

func (s *Server) createAccessRecord(w http.ResponseWriter, r *http.Request) {
	var req createAccessRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.CreateAccessRecord(model.AccessRecord{ClientID: req.ClientID, RuleID: req.RuleID, Path: req.Path, Allowed: req.Allowed, Remaining: req.Remaining, LatencyMs: req.LatencyMs})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, result)
}

func (s *Server) listAccessRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AccessRecordFilter{
		ClientID: r.URL.Query().Get("client_id"),
		RuleID:   r.URL.Query().Get("rule_id"),
		Path:     r.URL.Query().Get("path"),
	}
	if v := r.URL.Query().Get("allowed"); v != "" {
		allowed := v == "true"
		filter.Allowed = &allowed
	}
	items, total, err := s.svc.ListAccessRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAccessRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.GetAccessRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
