package handler

import (
	"net/http"

	"ratelimiter/internal/model"
	"ratelimiter/pkg/httpx"
)

func (s *Server) registerBucketRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/buckets", s.createBucket)
	mux.HandleFunc("GET /api/buckets", s.listBuckets)
	mux.HandleFunc("GET /api/buckets/{id}", s.getBucket)
	mux.HandleFunc("PUT /api/buckets/{id}", s.updateBucket)
	mux.HandleFunc("DELETE /api/buckets/{id}", s.deleteBucket)
	mux.HandleFunc("POST /api/buckets/{id}/refill", s.refillBucket)
}

type createBucketRequest struct {
	ClientID string  `json:"client_id"`
	RuleID   string  `json:"rule_id"`
	Tokens   float64 `json:"tokens"`
}

func (s *Server) createBucket(w http.ResponseWriter, r *http.Request) {
	var req createBucketRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.CreateBucket(model.Bucket{ClientID: req.ClientID, RuleID: req.RuleID, Tokens: req.Tokens})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, result)
}

func (s *Server) listBuckets(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.BucketFilter{
		ClientID: r.URL.Query().Get("client_id"),
		RuleID:   r.URL.Query().Get("rule_id"),
	}
	items, total, err := s.svc.ListBuckets(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBucket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.GetBucket(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type updateBucketRequest struct {
	ClientID string  `json:"client_id"`
	RuleID   string  `json:"rule_id"`
	Tokens   float64 `json:"tokens"`
}

func (s *Server) updateBucket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateBucketRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.UpdateBucket(id, model.Bucket{ClientID: req.ClientID, RuleID: req.RuleID, Tokens: req.Tokens})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

func (s *Server) deleteBucket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteBucket(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type refillBucketRequest struct {
	Tokens float64 `json:"tokens"`
}

func (s *Server) refillBucket(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req refillBucketRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.RefillBucket(id, req.Tokens)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
