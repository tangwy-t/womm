package badge

import "sync"

// Registry stores badges. Methods that return *Badge share the internal pointer —
// callers MUST NOT mutate returned Badges. Badges are registered once at startup
// and treated as immutable thereafter.
type Registry struct {
	mu     sync.RWMutex
	badges map[string]*Badge
}

func NewRegistry() *Registry {
	return &Registry{badges: make(map[string]*Badge)}
}

func (r *Registry) Register(b *Badge) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.badges[b.ID] = b
}

func (r *Registry) Lookup(id string) (*Badge, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.badges[id]
	return b, ok
}

func (r *Registry) ListAll() []*Badge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*Badge, 0, len(r.badges))
	for _, b := range r.badges {
		result = append(result, b)
	}
	return result
}

func (r *Registry) ListByType(bt BadgeType) []*Badge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*Badge
	for _, b := range r.badges {
		if b.Type == bt {
			result = append(result, b)
		}
	}
	return result
}

// RegisterAll registers all declarative and certified badges. Call once at startup.
func RegisterAll(r *Registry) {
	for _, b := range declarativeBadges {
		r.Register(b)
	}
	for _, b := range certifiedBadges {
		r.Register(b)
	}
}
