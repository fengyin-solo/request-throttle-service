package handler

import (
	"net/http"
	"time"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/httpx"
)

func (s *Server) registerAlertRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alerts", s.createAlert)
	mux.HandleFunc("GET /api/alerts", s.listAlerts)
	mux.HandleFunc("GET /api/alerts/{id}", s.getAlert)
	mux.HandleFunc("PUT /api/alerts/{id}", s.updateAlert)
	mux.HandleFunc("DELETE /api/alerts/{id}", s.deleteAlert)
}

type createAlertRequest struct {
	ClientID    string    `json:"client_id"`
	RuleID      string    `json:"rule_id"`
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	TriggeredAt time.Time `json:"triggered_at"`
}

func (s *Server) createAlert(w http.ResponseWriter, r *http.Request) {
	var req createAlertRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.CreateAlert(model.Alert{ClientID: req.ClientID, RuleID: req.RuleID, Level: req.Level, Message: req.Message, TriggeredAt: req.TriggeredAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, result)
}

func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AlertFilter{
		ClientID: r.URL.Query().Get("client_id"),
		Level:    r.URL.Query().Get("level"),
	}
	items, total, err := s.svc.ListAlerts(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.GetAlert(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type updateAlertRequest struct {
	ClientID    string    `json:"client_id"`
	RuleID      string    `json:"rule_id"`
	Level       string    `json:"level"`
	Message     string    `json:"message"`
	TriggeredAt time.Time `json:"triggered_at"`
}

func (s *Server) updateAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAlertRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.UpdateAlert(id, model.Alert{ClientID: req.ClientID, RuleID: req.RuleID, Level: req.Level, Message: req.Message, TriggeredAt: req.TriggeredAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) deleteAlert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAlert(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
