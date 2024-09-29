package cnab

import "sync"

// FieldHandlerStore is a store for field handlers.
type FieldHandlerStore interface {
	// GetFieldHandler returns the field handler for the given field type.
	GetFieldHandler(fieldType string) *FieldHandler
}

type store struct {
	// m map[string]*FieldHandler
	m sync.Map
}

type StoreOptions func(s *store)

// WithFieldHandler adds a field handler to the store.
func WithFieldHandler(fieldName string, handler *FieldHandler) StoreOptions {
	return func(s *store) {
		s.m.Store(fieldName, handler)
	}
}

// NewFieldHandlerStore creates a new FieldHandlerStore.
func NewFieldHandlerStore(storeOpts ...StoreOptions) FieldHandlerStore {
	store := &store{}
	for _, opt := range storeOpts {
		opt(store)
	}
	return store
}

// GetFieldHandler returns the field handler for the given field type.
func (s *store) GetFieldHandler(fieldType string) *FieldHandler {
	if v, ok := s.m.Load(fieldType); ok {
		return v.(*FieldHandler)
	}
	return nil
}
