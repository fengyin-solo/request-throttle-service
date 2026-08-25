package service

import (
	"sort"
	"time"
)

type StatsOverview struct {
	RuleCount        int            `json:"rule_count"`
	ClientCount      int            `json:"client_count"`
	WindowRequests   int            `json:"window_requests"`
	AllowRate        float64        `json:"allow_rate"`
	AlertLevelCounts map[string]int `json:"alert_level_counts"`
}

type TopClient struct {
	ClientID string `json:"client_id"`
	AppID    string `json:"app_id"`
	Count    int    `json:"count"`
}

func (s *Service) StatsOverview() (*StatsOverview, error) {
	rules := s.store.ListRules()
	clients := s.store.ListClients()
	records := s.store.ListAccessRecords()
	alerts := s.store.ListAlerts()

	window := 60 * time.Second
	now := time.Now().UTC()
	windowStart := now.Add(-window)
	windowRequests := 0
	allowedCount := 0
	for _, r := range records {
		if !r.CreatedAt.Before(windowStart) {
			windowRequests++
			if r.Allowed {
				allowedCount++
			}
		}
	}
	allowRate := 0.0
	if windowRequests > 0 {
		allowRate = float64(allowedCount) / float64(windowRequests)
	}
	alertLevelCounts := make(map[string]int)
	for _, a := range alerts {
		alertLevelCounts[a.Level]++
	}

	return &StatsOverview{
		RuleCount:        len(rules),
		ClientCount:      len(clients),
		WindowRequests:   windowRequests,
		AllowRate:        allowRate,
		AlertLevelCounts: alertLevelCounts,
	}, nil
}

func (s *Service) TopClients(n int) ([]TopClient, error) {
	records := s.store.ListAccessRecords()
	countMap := make(map[string]int)
	clientMap := make(map[string]string)
	for _, r := range records {
		countMap[r.ClientID]++
	}
	clients := s.store.ListClients()
	for _, c := range clients {
		clientMap[c.ID] = c.AppID
	}

	type pair struct {
		clientID string
		count    int
	}
	pairs := make([]pair, 0, len(countMap))
	for cid, cnt := range countMap {
		pairs = append(pairs, pair{clientID: cid, count: cnt})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].clientID > pairs[j].clientID
		}
		return pairs[i].count > pairs[j].count
	})

	result := make([]TopClient, 0, n)
	for i := 0; i < len(pairs) && i < n; i++ {
		result = append(result, TopClient{
			ClientID: pairs[i].clientID,
			AppID:    clientMap[pairs[i].clientID],
			Count:    pairs[i].count,
		})
	}
	return result, nil
}
