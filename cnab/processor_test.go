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
				"start": 1,
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
				"start": 1,
				"length": 3
			}
		]
	}`

	err := p.LoadSpec(ctx, strings.NewReader(specJSON))
	if !IsErrFieldHasNoTypeSpecified(err) {
		t.Fatalf("Expected ErrFieldHasNoTypeSpecified, got %v", err)
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
				"start": 1,
				"length": 3
			},
			{
				"name": "payment_date",
				"type": "date",
				"start": 4,
				"length": 8,
				"format": "YYYYMMDD"
			},
			{
				"name": "payment_amount",
				"type": "float",
				"start": 12,
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

func TestProcessor_PackRecord(t *testing.T) {
	ctx := context.Background()
	p := NewProcessor()

	specJSON := `
	{
		"fields": [
			{
				"name": "bank_code",
				"type": "int",
				"start": 1,
				"length": 3
			},
			{
				"name": "payment_date",
				"type": "date",
				"start": 4,
				"length": 8,
				"format": "YYYYMMDD"
			},
			{
				"name": "payment_amount",
				"type": "float",
				"start": 12,
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
				"start": 1,
				"length": 3
			},
			{
				"name": "payment_date",
				"type": "date",
				"start": 4,
				"length": 8,
				"format": "YYYYMMDD"
			},
			{
				"name": "payment_amount",
				"type": "float",
				"start": 12,
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
				"start": 1,
				"length": 3
			},
			{
				"name": "payment_date",
				"type": "date",
				"start": 4,
				"length": 8,
				"format": "YYYYMMDD"
			},
			{
				"name": "payment_amount",
				"type": "float",
				"start": 12,
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
