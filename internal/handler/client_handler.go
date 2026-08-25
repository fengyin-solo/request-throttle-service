package handler

import (
	"net/http"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/httpx"
)

func (s *Server) registerClientRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/clients", s.createClient)
	mux.HandleFunc("GET /api/clients", s.listClients)
	mux.HandleFunc("GET /api/clients/{id}", s.getClient)
	mux.HandleFunc("PUT /api/clients/{id}", s.updateClient)
	mux.HandleFunc("DELETE /api/clients/{id}", s.deleteClient)
	mux.HandleFunc("POST /api/clients/{id}/ban", s.banClient)
	mux.HandleFunc("POST /api/clients/{id}/unban", s.unbanClient)
}

type createClientRequest struct {
	AppID      string `json:"app_id"`
	Name       string `json:"name"`
	DailyQuota int64  `json:"daily_quota"`
	Status     string `json:"status"`
}

func (s *Server) createClient(w http.ResponseWriter, r *http.Request) {
	var req createClientRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.CreateClient(model.Client{AppID: req.AppID, Name: req.Name, DailyQuota: req.DailyQuota, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, result)
}

func (s *Server) listClients(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ClientFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListClients(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getClient(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.GetClient(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type updateClientRequest struct {
	AppID      string `json:"app_id"`
	Name       string `json:"name"`
	DailyQuota int64  `json:"daily_quota"`
	Status     string `json:"status"`
}

func (s *Server) updateClient(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateClientRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.UpdateClient(id, model.Client{AppID: req.AppID, Name: req.Name, DailyQuota: req.DailyQuota, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) deleteClient(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteClient(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) banClient(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.BanClient(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) unbanClient(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.UnbanClient(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
