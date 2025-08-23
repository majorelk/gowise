package assertions

import (
	"math"
	"testing"
	"time"
)

// TestEqualityFastPath tests the fast-path optimisations for comparable types.
func TestEqualityFastPath(t *testing.T) {
	tests := []struct {
		name       string
		got, want  interface{}
		shouldPass bool
	}{
		// Comparable types (fast path)
		{"int equal", 42, 42, true},
		{"int not equal", 42, 24, false},
		{"string equal", "hello", "hello", true},
		{"string not equal", "hello", "world", false},
		{"bool equal true", true, true, true},
		{"bool equal false", false, false, true},
		{"bool not equal", true, false, false},
		{"float64 equal", 3.14, 3.14, true},
		{"float64 not equal", 3.14, 2.71, false},
		{"float64 NaN", math.NaN(), math.NaN(), false}, // NaN != NaN

		// Nil comparisons
		{"both nil", nil, nil, true},
		{"got nil want value", nil, 42, false},
		{"got value want nil", 42, nil, false},

		// These will be tested separately due to pointer complexity

		// Different types (not comparable with ==)
		{"different types", 42, "42", false},
		{"int vs float", 42, 42.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := New(t)

			assert.Equal(tt.got, tt.want)

			hasError := assert.Error() != ""
			if tt.shouldPass && hasError {
				t.Errorf("Expected pass but got error: %s", assert.Error())
			} else if !tt.shouldPass && !hasError {
				t.Errorf("Expected failure but assertion passed")
			}
		})
	}

	// Test pointer comparisons separately
	t.Run("pointer identity", func(t *testing.T) {
		x := 42
		p1 := &x
		p2 := p1
		y := 42
		p3 := &y

		assert := New(t)

		// Same pointer should pass
		assert.Equal(p1, p2)
		if assert.Error() != "" {
			t.Errorf("Expected same pointer to be equal: %s", assert.Error())
		}

		// Different pointers should fail
		assert = New(t) // reset
		assert.Equal(p1, p3)
		if assert.Error() == "" {
			t.Errorf("Expected different pointers to be not equal")
		}
	})
}

// TestNotEqualFastPath tests NotEqual with fast-path optimisations.
func TestNotEqualFastPath(t *testing.T) {
	tests := []struct {
		name       string
		got, want  interface{}
		shouldPass bool
	}{
		{"int not equal", 42, 24, true},
		{"int equal", 42, 42, false},
		{"string not equal", "hello", "world", true},
		{"string equal", "hello", "hello", false},
		{"nil vs value", nil, 42, true},
		{"both nil", nil, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := New(t)

			assert.NotEqual(tt.got, tt.want)

			hasError := assert.Error() != ""
			if tt.shouldPass && hasError {
				t.Errorf("Expected pass but got error: %s", assert.Error())
			} else if !tt.shouldPass && !hasError {
				t.Errorf("Expected failure but assertion passed")
			}
		})
	}
}

// TestDeepEqualComplexTypes tests deep equality for non-comparable types.
func TestDeepEqualComplexTypes(t *testing.T) {
	tests := []struct {
		name       string
		got, want  interface{}
		shouldPass bool
	}{
		{
			"slices equal",
			[]int{1, 2, 3},
			[]int{1, 2, 3},
			true,
		},
		{
			"slices not equal",
			[]int{1, 2, 3},
			[]int{1, 2, 4},
			false,
		},
		{
			"maps equal",
			map[string]int{"a": 1, "b": 2},
			map[string]int{"a": 1, "b": 2},
			true,
		},
		{
			"maps not equal",
			map[string]int{"a": 1, "b": 2},
			map[string]int{"a": 1, "b": 3},
			false,
		},
		{
			"structs equal",
			struct{ X, Y int }{1, 2},
			struct{ X, Y int }{1, 2},
			true,
		},
		{
			"structs not equal",
			struct{ X, Y int }{1, 2},
			struct{ X, Y int }{1, 3},
			false,
		},
		{
			"time equal",
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := New(t)

			assert.DeepEqual(tt.got, tt.want)

			hasError := assert.Error() != ""
			if tt.shouldPass && hasError {
				t.Errorf("Expected pass but got error: %s", assert.Error())
			} else if !tt.shouldPass && !hasError {
				t.Errorf("Expected failure but assertion passed")
			}
		})
	}
}

// TestSamePointerIdentity tests the Same assertion for pointer identity.
func TestSamePointerIdentity(t *testing.T) {
	x, y := 42, 42
	px, py := &x, &y
	samePx := px

	tests := []struct {
		name       string
		got, want  interface{}
		shouldPass bool
	}{
		{"same pointer", px, samePx, true},
		{"different pointers same value", px, py, false},
		{"same value not pointers", x, x, true}, // Same memory location
		{"different values", 100, 200, false},   // Use different literals to avoid potential reuse
		{"nil pointers", (*int)(nil), (*int)(nil), true},
		{"nil vs non-nil", (*int)(nil), px, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := New(t)

			assert.Same(tt.got, tt.want)

			hasError := assert.Error() != ""
			if tt.shouldPass && hasError {
				t.Errorf("Expected pass but got error: %s", assert.Error())
			} else if !tt.shouldPass && !hasError {
				t.Errorf("Expected failure but assertion passed")
			}
		})
	}
}

// TestBooleanAssertions tests True and False assertions.
func TestBooleanAssertions(t *testing.T) {
	t.Run("True assertion", func(t *testing.T) {
		tests := []struct {
			name       string
			condition  bool
			shouldPass bool
		}{
			{"true condition", true, true},
			{"false condition", false, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert := New(t)

				assert.True(tt.condition)

				hasError := assert.Error() != ""
				if tt.shouldPass && hasError {
					t.Errorf("Expected pass but got error: %s", assert.Error())
				} else if !tt.shouldPass && !hasError {
					t.Errorf("Expected failure but assertion passed")
				}
			})
		}
	})

	t.Run("False assertion", func(t *testing.T) {
		tests := []struct {
			name       string
			condition  bool
			shouldPass bool
		}{
			{"false condition", false, true},
			{"true condition", true, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert := New(t)

				assert.False(tt.condition)

				hasError := assert.Error() != ""
				if tt.shouldPass && hasError {
					t.Errorf("Expected pass but got error: %s", assert.Error())
				} else if !tt.shouldPass && !hasError {
					t.Errorf("Expected failure but assertion passed")
				}
			})
		}
	})
}
