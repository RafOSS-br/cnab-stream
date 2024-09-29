package cnab

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestProcessor_LoadSpec(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	if len(p.(*processor).spec.Fields) != 1 {
		t.Fatalf("Expected 1 field, got %d", len(p.(*processor).spec.Fields))
	}
}

func TestProcessor_LoadSpec_FieldLengthZero(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 0
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if !IsErrLengthMustBeGreaterThanZero(err) {
		t.Fatalf("Expected ErrLengthMustBeGreaterThanZero, got %v", err)
	}
}

func TestProcessor_LoadSpec_InvalidJSON(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `invalid`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if !IsErrFailedToDecodeSpecJSON(err) {
		t.Fatalf("Expected ErrFailedToDecodeSpecJSON, got %v", err)
	}
}

func TestProcessor_LoadSpec_InvalidField(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "",
				"start": 0,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if !IsErrFieldHasNoTypeSpecified(err) {
		t.Fatalf("Expected ErrFieldHasNoTypeSpecified, got %v", err)
	}
}

func TestProcessor_LoadSpec_InvalidFieldStart(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": -1,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if !IsErrStartMustBeGreaterOrEqualZero(err) {
		t.Fatalf("Expected ErrLengthMustBeGreaterThanZero, got %v", err)
	}
}

func TestProcessor_LoadSpec_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	cancel()
	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if !IsErrCancelledContext(err) {
		t.Fatalf("Expected context error, got %v", err)
	}
}

func TestProcessor_ParseRecord(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			},
			{
				"name": "payment_date",
				"type": "date",
				"start": 3,
				"length": 8,
				"format": "YYYYMMDD"
			},
			{
				"name": "payment_amount",
				"type": "float",
				"start": 11,
				"length": 8,
				"decimal": 2
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	record := []byte("001202101011234567890") // bank_code=1, payment_date=2021-01-01, payment_amount=123456.78

	data, err := p.ParseRecord(ctx, record)
	if err != nil {
		t.Fatalf("Failed to parse record: %v", err)
	}

	if data["bank_code"] != 1 {
		t.Errorf("Expected bank_code 1, got %v", data["bank_code"])
	}

	expectedDate := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	if !data["payment_date"].(time.Time).Equal(expectedDate) {
		t.Errorf("Expected payment_date %v, got %v", expectedDate, data["payment_date"])
	}

	expectedAmount := 123456.78
	if data["payment_amount"] != expectedAmount {
		t.Errorf("Expected payment_amount %v, got %v", expectedAmount, data["payment_amount"])
	}
}

func TestProcessor_ParseRecord_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	record := []byte("001")

	cancel()
	_, err = p.ParseRecord(ctx, record)
	if !IsErrCancelledContext(err) {
		t.Fatalf("Expected context error, got %v", err)
	}
}

func TestProcessor_ParseRecord_WithErr(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	record := []byte("")

	_, err = p.ParseRecord(ctx, record)
	if err == nil {
		t.Fatalf("Expected ErrFailedToParseRecord, got %v", err)
	}
}

func TestProcessor_PackRecord(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			},
			{
				"name": "payment_date",
				"type": "date",
				"start": 3,
				"length": 8,
				"format": "YYYYMMDD"
			},
			{
				"name": "payment_amount",
				"type": "float",
				"start": 11,
				"length": 10,
				"decimal": 2
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	data := map[string]interface{}{
		"bank_code":      1,
		"payment_date":   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		"payment_amount": 12345678.90,
	} // bank_code=1, payment_date=2021-01-01, payment_amount=12345678.90

	record, err := p.PackRecord(ctx, data)
	if err != nil {
		t.Fatalf("Failed to pack record: %v", err)
	}

	expectedRecord := "001202101011234567890" // bank_code=1, payment_date=2021-01-01, payment_amount=12345678.90
	if string(record) != expectedRecord {
		t.Errorf("Expected record %s, got %s", expectedRecord, string(record))
	}
}

func TestProcessor_PackRecord_InvalidData(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	data := map[string]interface{}{}

	_, err = p.PackRecord(ctx, data)
	if err == nil {
		t.Fatalf("Expected ErrFailedToPackRecord, got %v", err)
	}
}

func TestProcessor_PackRecord_FieldExceedsSpecifiedLength(t *testing.T) {
	ctx := context.Background()

	mockedFieldHandlerStore := &MockFieldHandlerStore{}
	p := NewProcessor(
		WithFieldHandlerStore(mockedFieldHandlerStore),
	)

	// Mock GetFieldHandler
	mockedFieldHandler := &FieldHandler{
		Format: func(field *FieldSpec, value interface{}) (string, error) {
			return "ab", nil
		},
		Validate: func(field *FieldSpec) error {
			return nil
		},
	}
	mockedFieldHandlerStore.On("GetFieldHandler", "int").Return(mockedFieldHandler)

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 1
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	data := map[string]interface{}{
		"bank_code": 10,
	}

	_, err = p.PackRecord(ctx, data)
	if !IsErrFieldExceedsSpecifiedLength(err) {
		t.Fatalf("Expected ErrFieldValueExceedsSpecifiedLength, got %v", err)
	}
}

func TestProcessor_PackRecord_WithErrorOnFormatField(t *testing.T) {
	ctx := context.Background()

	mockedFieldHandlerStore := &MockFieldHandlerStore{}
	// Mock GetFieldHandler
	mockedFieldHandler := &FieldHandler{
		Format: func(field *FieldSpec, value interface{}) (string, error) {
			return "", ErrFailedToFormatField
		},
		Validate: func(field *FieldSpec) error {
			return nil
		},
	}
	mockedFieldHandlerStore.On("GetFieldHandler", "int").Return(mockedFieldHandler)

	p := NewProcessor(
		WithFieldHandlerStore(mockedFieldHandlerStore),
	)

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	data := map[string]interface{}{
		"bank_code": 1,
	}

	_, err = p.PackRecord(ctx, data)
	if !IsErrFailedToFormatField(err) {
		t.Fatalf("Expected ErrFailedToFormatField, got %v", err)
	}
}

func TestProcessor_PackRecord_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	data := map[string]interface{}{
		"bank_code": 1,
	}

	cancel()
	_, err = p.PackRecord(ctx, data)
	if !IsErrCancelledContext(err) {
		t.Fatalf("Expected context error, got %v", err)
	}
}

func TestProcessor_PackRecord_WithErr(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	data := map[string]interface{}{}

	_, err = p.PackRecord(ctx, data)
	if err == nil {
		t.Fatalf("Expected ErrFailedToPackRecord, got %v", err)
	}
}

func TestProcessor_formatFieldValue_Int(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "int",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	fieldSpec := FieldSpec{
		Name:   "int",
		Type:   "int",
		Start:  0,
		Length: 3,
	}

	p.LoadSpec(ctx, strings.NewReader(specJSON))

	value, err := p.(*processor).formatFieldValue(&fieldSpec, 001)
	if err != nil {
		t.Fatalf("Failed to format field value: %v", err)
	}

	if value != "001" {
		t.Errorf("Expected 001, got %v", value)
	}
}

func Test_Spec(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "int",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		t.Fatalf("Failed to load spec: %v", err)
	}

	expected := "&{Name:int Type:int Start:0 Length:3 Format: Decimal:0 End:2}"

	if p.Spec() != expected {
		t.Errorf("Expected %s, got %s", expected, "'"+p.Spec()+"'")
	}
}

func Test_formatFieldValue_InvalidType(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "invalid",
				"type": "invalid",
				"start": 0,
				"length": 3
			}
			]
		}`

	fieldSpec := FieldSpec{
		Name:   "invalid",
		Type:   "invalid",
		Start:  1,
		Length: 3,
	}

	p.LoadSpec(ctx, strings.NewReader(specJSON))

	_, err := p.(*processor).formatFieldValue(&fieldSpec, "invalid")
	if err == nil {
		t.Fatalf("Expected error for invalid field type, got nil")
	}
}

func Test_parseRecord_UnsupportedFieldType(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "unsupported",
				"type": "unsupported",
				"start": 0,
				"length": 3
			}
		]
	}`

	field := &FieldSpec{
		Name:   "unsupported",
		Type:   "unsupported",
		Start:  0,
		Length: 3,
	}

	m := make(map[string]interface{})
	m["unsupported"] = "001"

	p.LoadSpec(ctx, strings.NewReader(specJSON))

	err := p.(*processor).parseRecord([]byte("001"), field, m)
	if err == nil {
		t.Fatalf("Expected error for unsupported field type, got nil")
	}
}

func Test_parseRecord_MockedErrorOnParseField(t *testing.T) {
	ctx := context.Background()

	mockedFieldHandlerStore := &MockFieldHandlerStore{}
	p := NewProcessor(
		WithFieldHandlerStore(mockedFieldHandlerStore),
	)

	// Mock GetFieldHandler
	mockedFieldHandler := &FieldHandler{
		Validate: func(field *FieldSpec) error {
			return nil
		},
		Parse: func(field *FieldSpec, rawValue []byte) (interface{}, error) {
			return nil, ErrFailedToParseField
		},
	}
	mockedFieldHandlerStore.On("GetFieldHandler", "int").Return(mockedFieldHandler)

	specJSON := `
	{
		"fields": [
			{
				"name": "int",
				"type": "int",
				"start": 0,
				"length": 3
			}
		]
	}`

	field := &FieldSpec{
		Name:   "int",
		Type:   "int",
		Start:  0,
		Length: 3,
	}

	m := make(map[string]interface{})
	m["int"] = 1

	p.LoadSpec(ctx, strings.NewReader(specJSON))

	err := p.(*processor).parseRecord([]byte("001"), field, m)
	if !IsErrFailedToParseField(err) {
		t.Fatalf("Expected ErrFailedToFormatField, got %v", err)
	}
}

// Benchmark Tests

func BenchmarkProcessor_ParseRecord(b *testing.B) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			},
			{
				"name": "payment_date",
				"type": "date",
				"start": 3,
				"length": 8,
				"format": "YYYYMMDD"
			},
			{
				"name": "payment_amount",
				"type": "float",
				"start": 11,
				"length": 10,
				"decimal": 2
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		b.Fatalf("Failed to load spec: %v", err)
	}

	record := []byte("001202101011234567890")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.ParseRecord(ctx, record)
		if err != nil {
			b.Fatalf("Failed to parse record: %v", err)
		}
	}
}

func BenchmarkProcessor_PackRecord(b *testing.B) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 0,
				"length": 3
			},
			{
				"name": "payment_date",
				"type": "date",
				"start": 3,
				"length": 8,
				"format": "YYYYMMDD"
			},
			{
				"name": "payment_amount",
				"type": "float",
				"start": 11,
				"length": 10,
				"decimal": 2
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if err != nil {
		b.Fatalf("Failed to load spec: %v", err)
	}

	data := map[string]interface{}{
		"bank_code":      1,
		"payment_date":   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		"payment_amount": 123456.78,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.PackRecord(ctx, data)
		if err != nil {
			b.Fatalf("Failed to pack record: %v", err)
		}
	}
}
