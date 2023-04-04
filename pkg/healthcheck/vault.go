package healthcheck

import (
	"sync"
	"time"
)

type StatusVault struct {
	status ProbeStatus
	mu     sync.Mutex
}

func NewStatusVault(status ProbeStatus) *StatusVault {
	return &StatusVault{status: status}
}

func (v *StatusVault) Load() ProbeStatus {
	v.mu.Lock()
	defer v.mu.Unlock()

	return v.status
}

func (v *StatusVault) Set(status ProbeStatus) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.status = status
	v.status.UpdatedAt = time.Now()
}

func (v *StatusVault) Update(fn func(s *ProbeStatus)) {
	v.mu.Lock()
	defer v.mu.Unlock()

	old := v.status

	fn(&v.status)

	if v.status.Ready != old.Ready || v.status.Status != old.Status {
		v.status.UpdatedAt = time.Now()
	}
}
