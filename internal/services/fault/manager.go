package fault

import (
	"sync"
)

type Factory func(options Options) (Fault, error)

type Manager struct {
	mtx       sync.RWMutex
	factories map[string]Factory
}

func (m *Manager) Register(factory Factory, name string) bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	if _, ok := m.factories[name]; ok {
		return false
	}

	m.factories[name] = factory

	return true
}

func (m *Manager) Lookup(name string) (Factory, bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	factory, ok := m.factories[name]
	if !ok {
		return nil, ok
	}

	return factory, ok
}

func NewManager() *Manager {
	return &Manager{
		mtx:       sync.RWMutex{},
		factories: map[string]Factory{},
	}
}
