package generics

import (
	"testing"
)

func TestValidateAll(t *testing.T) {
	// Test successful validation
	err := ValidateAll("test",
		NotEmpty,
		MinLength(3),
		MaxLength(10),
	)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test validation failure
	err = ValidateAll("test",
		NotEmpty,
		MinLength(10), // Should fail
	)
	if err == nil {
		t.Error("Expected validation error")
	}
}

func TestValidateSlice(t *testing.T) {
	// Test successful validation
	slice := []string{"test1", "test2", "test3"}
	err := ValidateSlice(slice, NotEmpty)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test validation failure
	sliceWithEmpty := []string{"test1", "", "test3"}
	err = ValidateSlice(sliceWithEmpty, NotEmpty)
	if err == nil {
		t.Error("Expected validation error")
	}
}

func TestValidateMap(t *testing.T) {
	// Test successful validation
	m := map[string]string{"key1": "value1", "key2": "value2"}
	err := ValidateMap(m, NotEmpty)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test validation failure
	mWithEmpty := map[string]string{"key1": "value1", "key2": ""}
	err = ValidateMap(mWithEmpty, NotEmpty)
	if err == nil {
		t.Error("Expected validation error")
	}
}

func TestNotEmpty(t *testing.T) {
	// Test valid cases
	if err := NotEmpty("test"); err != nil {
		t.Errorf("Expected no error for 'test', got %v", err)
	}
	if err := NotEmpty(" test "); err != nil {
		t.Errorf("Expected no error for ' test ', got %v", err)
	}

	// Test invalid cases
	if err := NotEmpty(""); err == nil {
		t.Error("Expected error for empty string")
	}
	if err := NotEmpty("   "); err == nil {
		t.Error("Expected error for whitespace string")
	}
	if err := NotEmpty("\t\n"); err == nil {
		t.Error("Expected error for whitespace string")
	}
}

func TestMinLength(t *testing.T) {
	validator := MinLength(3)

	// Test valid cases
	if err := validator("test"); err != nil {
		t.Errorf("Expected no error for 'test', got %v", err)
	}
	if err := validator("testing"); err != nil {
		t.Errorf("Expected no error for 'testing', got %v", err)
	}

	// Test invalid cases
	if err := validator("ab"); err == nil {
		t.Error("Expected error for 'ab'")
	}
	if err := validator(""); err == nil {
		t.Error("Expected error for empty string")
	}
}

func TestMaxLength(t *testing.T) {
	validator := MaxLength(5)

	// Test valid cases
	if err := validator("test"); err != nil {
		t.Errorf("Expected no error for 'test', got %v", err)
	}
	if err := validator(""); err != nil {
		t.Errorf("Expected no error for empty string, got %v", err)
	}

	// Test invalid cases
	if err := validator("testing"); err == nil {
		t.Error("Expected error for 'testing'")
	}
}

func TestLength(t *testing.T) {
	validator := Length(4)

	// Test valid cases
	if err := validator("test"); err != nil {
		t.Errorf("Expected no error for 'test', got %v", err)
	}

	// Test invalid cases
	if err := validator("abc"); err == nil {
		t.Error("Expected error for 'abc'")
	}
	if err := validator("testing"); err == nil {
		t.Error("Expected error for 'testing'")
	}
}

func TestMin(t *testing.T) {
	validator := Min(10)

	// Test valid cases
	if err := validator(15); err != nil {
		t.Errorf("Expected no error for 15, got %v", err)
	}
	if err := validator(10); err != nil {
		t.Errorf("Expected no error for 10, got %v", err)
	}

	// Test invalid cases
	if err := validator(5); err == nil {
		t.Error("Expected error for 5")
	}
	if err := validator(-5); err == nil {
		t.Error("Expected error for -5")
	}
}

func TestMax(t *testing.T) {
	validator := Max(20)

	// Test valid cases
	if err := validator(15); err != nil {
		t.Errorf("Expected no error for 15, got %v", err)
	}
	if err := validator(20); err != nil {
		t.Errorf("Expected no error for 20, got %v", err)
	}

	// Test invalid cases
	if err := validator(25); err == nil {
		t.Error("Expected error for 25")
	}
}

func TestRange(t *testing.T) {
	validator := Range(10, 20)

	// Test valid cases
	if err := validator(15); err != nil {
		t.Errorf("Expected no error for 15, got %v", err)
	}
	if err := validator(10); err != nil {
		t.Errorf("Expected no error for 10, got %v", err)
	}
	if err := validator(20); err != nil {
		t.Errorf("Expected no error for 20, got %v", err)
	}

	// Test invalid cases
	if err := validator(5); err == nil {
		t.Error("Expected error for 5")
	}
	if err := validator(25); err == nil {
		t.Error("Expected error for 25")
	}
}

func TestNotNil(t *testing.T) {
	// Test valid cases
	str := "test"
	if err := NotNil(&str); err != nil {
		t.Errorf("Expected no error for non-nil pointer, got %v", err)
	}

	// Test invalid cases
	if err := NotNil[string](nil); err == nil {
		t.Error("Expected error for nil pointer")
	}
}

func TestNotZero(t *testing.T) {
	// Test valid cases
	if err := NotZero(42); err != nil {
		t.Errorf("Expected no error for 42, got %v", err)
	}
	if err := NotZero("test"); err != nil {
		t.Errorf("Expected no error for 'test', got %v", err)
	}

	// Test invalid cases
	if err := NotZero(0); err == nil {
		t.Error("Expected error for 0")
	}
	if err := NotZero(""); err == nil {
		t.Error("Expected error for empty string")
	}
}

func TestIn(t *testing.T) {
	validator := In("option1", "option2", "option3")

	// Test valid cases
	if err := validator("option1"); err != nil {
		t.Errorf("Expected no error for 'option1', got %v", err)
	}
	if err := validator("option2"); err != nil {
		t.Errorf("Expected no error for 'option2', got %v", err)
	}

	// Test invalid cases
	if err := validator("option4"); err == nil {
		t.Error("Expected error for 'option4'")
	}
	if err := validator(""); err == nil {
		t.Error("Expected error for empty string")
	}
}

func TestNotIn(t *testing.T) {
	validator := NotIn("forbidden1", "forbidden2")

	// Test valid cases
	if err := validator("allowed1"); err != nil {
		t.Errorf("Expected no error for 'allowed1', got %v", err)
	}
	if err := validator("allowed2"); err != nil {
		t.Errorf("Expected no error for 'allowed2', got %v", err)
	}

	// Test invalid cases
	if err := validator("forbidden1"); err == nil {
		t.Error("Expected error for 'forbidden1'")
	}
	if err := validator("forbidden2"); err == nil {
		t.Error("Expected error for 'forbidden2'")
	}
}

func TestPositive(t *testing.T) {
	validator := Positive[int]()

	// Test valid cases
	if err := validator(1); err != nil {
		t.Errorf("Expected no error for 1, got %v", err)
	}
	if err := validator(100); err != nil {
		t.Errorf("Expected no error for 100, got %v", err)
	}
	if err := validator(0); err != nil {
		t.Errorf("Expected no error for 0, got %v", err)
	}

	// Test invalid cases
	if err := validator(-1); err == nil {
		t.Error("Expected error for -1")
	}
}

func TestNegative(t *testing.T) {
	validator := Negative[int]()

	// Test valid cases
	if err := validator(-1); err != nil {
		t.Errorf("Expected no error for -1, got %v", err)
	}
	if err := validator(-100); err != nil {
		t.Errorf("Expected no error for -100, got %v", err)
	}

	// Test invalid cases
	if err := validator(0); err == nil {
		t.Error("Expected error for 0")
	}
	if err := validator(1); err == nil {
		t.Error("Expected error for 1")
	}
}

func TestSliceMinLength(t *testing.T) {
	validator := SliceMinLength[int](3)

	// Test valid cases
	if err := validator([]int{1, 2, 3}); err != nil {
		t.Errorf("Expected no error for slice with 3 elements, got %v", err)
	}
	if err := validator([]int{1, 2, 3, 4, 5}); err != nil {
		t.Errorf("Expected no error for slice with 5 elements, got %v", err)
	}

	// Test invalid cases
	if err := validator([]int{1, 2}); err == nil {
		t.Error("Expected error for slice with 2 elements")
	}
	if err := validator([]int{}); err == nil {
		t.Error("Expected error for empty slice")
	}
}

func TestSliceMaxLength(t *testing.T) {
	validator := SliceMaxLength[int](3)

	// Test valid cases
	if err := validator([]int{1, 2, 3}); err != nil {
		t.Errorf("Expected no error for slice with 3 elements, got %v", err)
	}
	if err := validator([]int{1, 2}); err != nil {
		t.Errorf("Expected no error for slice with 2 elements, got %v", err)
	}
	if err := validator([]int{}); err != nil {
		t.Errorf("Expected no error for empty slice, got %v", err)
	}

	// Test invalid cases
	if err := validator([]int{1, 2, 3, 4}); err == nil {
		t.Error("Expected error for slice with 4 elements")
	}
}

func TestMapMinSize(t *testing.T) {
	validator := MapMinSize[string, int](2)

	// Test valid cases
	m1 := map[string]int{"a": 1, "b": 2}
	if err := validator(m1); err != nil {
		t.Errorf("Expected no error for map with 2 elements, got %v", err)
	}
	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	if err := validator(m2); err != nil {
		t.Errorf("Expected no error for map with 3 elements, got %v", err)
	}

	// Test invalid cases
	m3 := map[string]int{"a": 1}
	if err := validator(m3); err == nil {
		t.Error("Expected error for map with 1 element")
	}
	m4 := map[string]int{}
	if err := validator(m4); err == nil {
		t.Error("Expected error for empty map")
	}
}

func TestMapMaxSize(t *testing.T) {
	validator := MapMaxSize[string, int](2)

	// Test valid cases
	m1 := map[string]int{"a": 1, "b": 2}
	if err := validator(m1); err != nil {
		t.Errorf("Expected no error for map with 2 elements, got %v", err)
	}
	m2 := map[string]int{"a": 1}
	if err := validator(m2); err != nil {
		t.Errorf("Expected no error for map with 1 element, got %v", err)
	}
	m3 := map[string]int{}
	if err := validator(m3); err != nil {
		t.Errorf("Expected no error for empty map, got %v", err)
	}

	// Test invalid cases
	m4 := map[string]int{"a": 1, "b": 2, "c": 3}
	if err := validator(m4); err == nil {
		t.Error("Expected error for map with 3 elements")
	}
}

func TestValidationError_Error(t *testing.T) {
	// Test error without field
	err := &ValidationError{Message: "test error"}
	expected := "test error"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}

	// Test error with field
	errWithField := &ValidationError{Field: "email", Message: "invalid format"}
	expectedWithField := "email: invalid format"
	if errWithField.Error() != expectedWithField {
		t.Errorf("Expected '%s', got '%s'", expectedWithField, errWithField.Error())
	}
}
