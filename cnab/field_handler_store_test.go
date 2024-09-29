package cnab

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockFieldHandlerStore is a mock type for FieldHandlerStore.
type MockFieldHandlerStore struct {
	mock.Mock
}

// GetFieldHandler returns the field handler for the given field type.
func (_m *MockFieldHandlerStore) GetFieldHandler(fieldType string) *FieldHandler {
	ret := _m.Called(fieldType)

	var r0 *FieldHandler
	if rf, ok := ret.Get(0).(func(string) *FieldHandler); ok {
		r0 = rf(fieldType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*FieldHandler)
		}
	}

	return r0
}

// TestGetFieldHandler tests the GetFieldHandler method.
func TestGetFieldHandler(t *testing.T) {
	store := NewFieldHandlerStore(WithFieldHandler("test", &FieldHandler{}))
	handler := store.GetFieldHandler("test")
	if handler == nil {
		t.Fatal("Expected a field handler")
	}
	handler = store.GetFieldHandler("missing")
	if handler != nil {
		t.Fatal("Expected no field handler")
	}
}
