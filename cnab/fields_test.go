package cnab

import (
	"math"
	"testing"
	"time"
)

func TestValidateIntField(t *testing.T) {
	field := FieldSpec{
		Name: "TestField",
	}
	err := validateIntField(field)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestParseIntField(t *testing.T) {
	field := FieldSpec{
		Name: "TestIntField",
	}

	// Test empty rawValue
	rawValue := []byte("   ")
	_, err := parseIntField(field, rawValue)
	if !IsErrFieldIsEmpty(err) {
		t.Errorf("Expected error for empty rawValue, got %v", err)
	}

	// Test valid integer
	rawValue = []byte("  12345 ")
	value, err := parseIntField(field, rawValue)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if value != 12345 {
		t.Errorf("Expected 12345, got %v", value)
	}

	// Test invalid integer
	rawValue = []byte(" 12a34 ")
	_, err = parseIntField(field, rawValue)
	if !IsErrFieldValueIsNotAnInt(err) {
		t.Errorf("Expected error for invalid integer input, got %v", err)
	}
}

func TestFormatIntField(t *testing.T) {
	field := FieldSpec{
		Name:   "TestIntField",
		Length: 5,
	}

	// Test with integer value
	value := 123
	str, err := formatIntField(field, value)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if str != "00123" {
		t.Errorf("Expected '00123', got '%s'", str)
	}

	// Test with string value that can be converted to int
	str, err = formatIntField(field, "456")
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if str != "00456" {
		t.Errorf("Expected '00456', got '%s'", str)
	}

	// Test with value that cannot be converted to int
	_, err = formatIntField(field, "abc")
	if !IsErrFieldValueIsNotAnInt(err) {
		t.Errorf("Expected error for invalid int value, got %v", err)
	}
}

func TestValidateFloatField(t *testing.T) {
	// Test with valid decimal
	field := FieldSpec{
		Name:    "TestFloatField",
		Decimal: 2,
	}
	err := validateFloatField(field)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	// Test with invalid decimal
	field.Decimal = -1
	err = validateFloatField(field)
	if !IsErrInvalidDecimalValue(err) {
		t.Errorf("Expected error for negative decimal, got %v", err)
	}
}

func TestParseFloatField(t *testing.T) {
	field := FieldSpec{
		Name:    "TestFloatField",
		Decimal: 2,
	}

	// Test with empty rawValue
	rawValue := []byte("  ")
	_, err := parseFloatField(field, rawValue)
	if !IsErrFieldIsEmpty(err) {
		t.Errorf("Expected error for empty rawValue, got %v", err)
	}

	// Test with valid float value
	rawValue = []byte(" 12345 ")
	value, err := parseFloatField(field, rawValue)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	expectedValue := 123.45 // because decimal is 2
	if value != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, value)
	}

	// Test with invalid integer in rawValue
	rawValue = []byte("12a34")
	_, err = parseFloatField(field, rawValue)
	if !IsErrFieldValueIsNotAnFloat(err) {
		t.Errorf("Expected error for invalid integer input, got %v", err)
	}
}

func TestFormatFloatField(t *testing.T) {
	field := FieldSpec{
		Name:    "TestFloatField",
		Length:  7,
		Decimal: 2,
	}

	// Test with float64 value
	value := 123.45
	str, err := formatFloatField(field, value)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if str != "0012345" {
		t.Errorf("Expected '0012345', got '%s'", str)
	}

	// Test with int value
	value = 123
	str, err = formatFloatField(field, value)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if str != "0012300" { // 123 * 100 = 12300
		t.Errorf("Expected '0012300', got '%s'", str)
	}

	// Test with string value that can be converted to float
	str, err = formatFloatField(field, "678.90")
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if str != "0067890" {
		t.Errorf("Expected '0067890', got '%s'", str)
	}

	// Test with invalid value
	_, err = formatFloatField(field, "abc")
	if !IsErrFieldValueIsNotAnFloat(err) {
		t.Errorf("Expected error for invalid float value, got %v", err)
	}
}

func TestValidateDateField(t *testing.T) {
	field := FieldSpec{
		Name:   "TestDateField",
		Format: "YYYYMMDD",
	}
	err := validateDateField(field)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}

	field.Format = ""
	err = validateDateField(field)
	if !IsErrMissingDateFormat(err) {
		t.Errorf("Expected error for missing date format, got %v", err)
	}
}

func TestParseDateField(t *testing.T) {
	field := FieldSpec{
		Name:   "TestDateField",
		Format: "YYYYMMDD",
	}

	// Test empty rawValue
	rawValue := []byte("   ")
	_, err := parseDateField(field, rawValue)
	if !IsErrFieldIsEmpty(err) {
		t.Errorf("Expected error for empty rawValue, got %v", err)
	}

	// Test rawValue with invalid length
	rawValue = []byte("202101")
	_, err = parseDateField(field, rawValue)
	if !IsErrInvalidDateLength(err) {
		t.Errorf("Expected error for invalid date length, got %v", err)
	}

	// Test valid date
	rawValue = []byte("20210315")
	value, err := parseDateField(field, rawValue)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	expectedDate := time.Date(2021, 3, 15, 0, 0, 0, 0, time.UTC)
	if !value.(time.Time).Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, value)
	}

	// Test invalid date
	rawValue = []byte("20211315")
	_, err = parseDateField(field, rawValue)
	if !IsErrFieldValueIsNotAnDate(err) {
		t.Errorf("Expected error for invalid date, got %v", err)
	}
}

func TestFormatDateField(t *testing.T) {
	field := FieldSpec{
		Name:   "TestDateField",
		Format: "YYYYMMDD",
	}

	// Test with time.Time value
	value := time.Date(2021, 3, 15, 0, 0, 0, 0, time.UTC)
	str, err := formatDateField(field, value)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if str != "20210315" {
		t.Errorf("Expected '20210315', got '%s'", str)
	}

	// Test with invalid value type
	valueInvalid := "20210315"
	_, err = formatDateField(field, valueInvalid)
	if !IsErrFieldValueIsNotAnDate(err) {
		t.Errorf("Expected error for invalid date value, got %v", err)
	}
}

func TestValidateStringField(t *testing.T) {
	field := FieldSpec{
		Name: "TestStringField",
	}
	err := validateStringField(field)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestParseStringField(t *testing.T) {
	field := FieldSpec{
		Name: "TestStringField",
	}
	rawValue := []byte("  Hello World  ")
	value, err := parseStringField(field, rawValue)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if value != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", value)
	}
}

func TestFormatStringField(t *testing.T) {
	field := FieldSpec{
		Name: "TestStringField",
	}

	// Test with string value
	value := "Hello World"
	str, err := formatStringField(field, value)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if str != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", str)
	}

	// Test with invalid value type
	valueInvalid := 12345
	_, err = formatStringField(field, valueInvalid)
	if !IsErrFieldValueIsNotAnString(err) {
		t.Errorf("Expected error for invalid string value, got %v", err)
	}
}

func TestAtoiUnsafe(t *testing.T) {
	// Test with valid digits
	b := []byte("12345")
	value, err := atoiUnsafe(b)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if value != 12345 {
		t.Errorf("Expected 12345, got %d", value)
	}

	// Test with invalid digits
	b = []byte("12a45")
	_, err = atoiUnsafe(b)
	if !IsErrFieldValueIsNotAnInt(err) {
		t.Errorf("Expected error for invalid integer input, got %v", err)
	}
}

func TestParseFloatBytes(t *testing.T) {
	// Valid input
	b := []byte("12345")
	decimal := 2
	value, err := parseFloatBytes(b, decimal)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	expectedValue := 123.45
	if value != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, value)
	}

	// Invalid input
	b = []byte("12a45")
	_, err = parseFloatBytes(b, decimal)
	if !IsErrFieldValueIsNotAnFloat(err) {
		t.Errorf("Expected error for invalid integer input, got %v", err)
	}

	// decimal = 0
	b = []byte("12345")
	decimal = 0
	value, err = parseFloatBytes(b, decimal)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	expectedValue = 12345
	if value != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, value)
	}
}

func TestParseDateBytes(t *testing.T) {
	// Valid date
	b := []byte("20210315")
	format := "YYYYMMDD"
	value, err := parseDateBytes(b, format)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	expectedDate := time.Date(2021, 3, 15, 0, 0, 0, 0, time.UTC)
	if !value.Equal(expectedDate) {
		t.Errorf("Expected %v, got %v", expectedDate, value)
	}

	// Invalid length
	b = []byte("2021031")
	_, err = parseDateBytes(b, format)
	if !IsErrInvalidDateLength(err) {
		t.Errorf("Expected error for invalid date length, got %v", err)
	}

	// Invalid date string
	b = []byte("20211315")
	_, err = parseDateBytes(b, format)
	if !IsErrFieldValueIsNotAnDate(err) {
		t.Errorf("Expected error for invalid date, got %v", err)
	}
}

func TestPow10(t *testing.T) {
	// n = 0
	if pow10(0) != 1.0 {
		t.Errorf("Expected pow10(0) == 1.0, got %v", pow10(0))
	}

	// n = 3
	if pow10(3) != 1000.0 {
		t.Errorf("Expected pow10(3) == 1000.0, got %v", pow10(3))
	}

	// n negative
	if pow10(-2) != 1.0 {
		t.Errorf("Expected pow10(-2) == 1.0 (since loop not entered), got %v", pow10(-2))
	}
}

func TestConvertDateFormat(t *testing.T) {
	// All tokens
	format := "YYYYMMDD"
	expected := "20060102"
	result := convertDateFormat(format)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Some tokens
	format = "YYYY-MM-DD"
	expected = "2006-01-02"
	result = convertDateFormat(format)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// No tokens
	format = "ABCD"
	expected = "ABCD"
	result = convertDateFormat(format)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFormatDate(t *testing.T) {
	value := time.Date(2021, 3, 15, 0, 0, 0, 0, time.UTC)
	format := "YYYYMMDD"
	expected := "20210315"
	result := formatDate(value, format)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestToInt(t *testing.T) {
	// int
	value := 123
	result, err := toInt(value)
	if err != nil || result != 123 {
		t.Errorf("Expected 123, got %d, error %v", result, err)
	}

	// int64
	valueInt64 := int64(456)
	result, err = toInt(valueInt64)
	if err != nil || result != 456 {
		t.Errorf("Expected 456, got %d, error %v", result, err)
	}

	// float64
	valueFloat64 := float64(789.0)
	result, err = toInt(valueFloat64)
	if err != nil || result != 789 {
		t.Errorf("Expected 789, got %d, error %v", result, err)
	}

	// string representing int
	valueStr := "321"
	result, err = toInt(valueStr)
	if err != nil || result != 321 {
		t.Errorf("Expected 321, got %d, error %v", result, err)
	}

	// string not representing int
	valueStr = "abc"
	_, err = toInt(valueStr)
	if !IsErrFieldValueIsNotAnInt(err) {
		t.Errorf("Expected error for invalid string, got %v", err)
	}

	// unsupported type
	valueUnsupported := []int{1, 2, 3}
	_, err = toInt(valueUnsupported)
	if !IsErrFieldValueIsNotAnInt(err) {
		t.Errorf("Expected error for unsupported type, got %v", err)
	}
}

func TestToFloat(t *testing.T) {
	// float64
	value := float64(123.45)
	result, err := toFloat(value)
	if err != nil || result != 123.45 {
		t.Errorf("Expected 123.45, got %v, error %v", result, err)
	}

	// float32
	valueFloat32 := float32(678.90)
	const tolerance = 1e-4
	result, err = toFloat(valueFloat32)
	if err != nil || math.Abs(result-678.90) > tolerance {
		t.Errorf("Expected 678.90, got %v, error %v", result, err)
	}

	// int
	valueInt := 123
	result, err = toFloat(valueInt)
	if err != nil || result != 123.0 {
		t.Errorf("Expected 123.0, got %v, error %v", result, err)
	}

	// string representing float
	valueStr := "456.78"
	result, err = toFloat(valueStr)
	if err != nil || result != 456.78 {
		t.Errorf("Expected 456.78, got %v, error %v", result, err)
	}

	// string not representing float
	valueStr = "abc"
	_, err = toFloat(valueStr)
	if !IsErrFieldValueIsNotAnFloat(err) {
		t.Errorf("Expected error for invalid string, got %v", err)
	}

	// unsupported type
	valueUnsupported := []int{1, 2, 3}
	_, err = toFloat(valueUnsupported)
	if !IsErrFieldValueIsNotAnFloat(err) {
		t.Errorf("Expected error for unsupported type, got %v", err)
	}
}

func TestFieldToError(t *testing.T) {
	err := fieldToError("Name", "TestField")
	expected := "field Name TestField"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}
