package healthcheck

import (
	"context"
	"encoding/json"
	"os"
	"time"
)

type ProbeFunc func(ctx context.Context) ProbeStatus

type ProbeStatus struct {
	Hostname  string
	Component string
	Status    string
	Ready     bool
	Critical  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProbeStatus(component string, critical bool, status string) ProbeStatus {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	return ProbeStatus{
		Hostname:  hostname,
		Component: component,
		Critical:  critical,
		Status:    status,
		CreatedAt: time.Now(),
	}
}

func (s *ProbeStatus) MarshalJSON() ([]byte, error) {
	type resp struct {
		Hostname  string `json:"hostname"`
		Critical  bool   `json:"critical,omitempty"`
		Status    string `json:"status,omitempty"`
		Ready     bool   `json:"ready"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at,omitempty"`
	}

	r := resp{
		Hostname:  s.Hostname,
		Status:    s.Status,
		Critical:  s.Critical,
		Ready:     s.Ready,
		CreatedAt: s.CreatedAt.In(time.UTC).Format(time.RFC3339),
	}
	if !s.UpdatedAt.IsZero() {
		r.UpdatedAt = s.UpdatedAt.In(time.UTC).Format(time.RFC3339)
	}

	return json.Marshal(&r)
}
