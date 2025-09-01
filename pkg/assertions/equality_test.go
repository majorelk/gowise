package assertions

import (
	"math"
	"testing"
	"time"
)

// behaviorMockT is defined in assertions_passing_test.go - shared across test files

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
			// Test GoWise framework behavioral contract
			mock := &behaviorMockT{}
			assert := New(mock)

			assert.Equal(tt.got, tt.want)

			// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
			if tt.shouldPass && len(mock.errorCalls) != 0 {
				t.Errorf("Equal should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			} else if !tt.shouldPass && len(mock.errorCalls) != 1 {
				t.Errorf("Equal should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
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

		// Same pointer should pass
		mock1 := &behaviorMockT{}
		assert1 := New(mock1)
		assert1.Equal(p1, p2)
		if len(mock1.errorCalls) != 0 {
			t.Errorf("Same pointer should pass (no Errorf calls), got %d: %v", len(mock1.errorCalls), mock1.errorCalls)
		}

		// Different pointers should fail
		mock2 := &behaviorMockT{}
		assert2 := New(mock2)
		assert2.Equal(p1, p3)
		if len(mock2.errorCalls) != 1 {
			t.Errorf("Different pointers should fail (1 Errorf call), got %d: %v", len(mock2.errorCalls), mock2.errorCalls)
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
			// Test GoWise framework behavioral contract
			mock := &behaviorMockT{}
			assert := New(mock)

			assert.NotEqual(tt.got, tt.want)

			// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
			if tt.shouldPass && len(mock.errorCalls) != 0 {
				t.Errorf("NotEqual should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			} else if !tt.shouldPass && len(mock.errorCalls) != 1 {
				t.Errorf("NotEqual should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
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
			// Test GoWise framework behavioral contract
			mock := &behaviorMockT{}
			assert := New(mock)

			assert.DeepEqual(tt.got, tt.want)

			// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
			if tt.shouldPass && len(mock.errorCalls) != 0 {
				t.Errorf("DeepEqual should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			} else if !tt.shouldPass && len(mock.errorCalls) != 1 {
				t.Errorf("DeepEqual should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
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
			// Test GoWise framework behavioral contract
			mock := &behaviorMockT{}
			assert := New(mock)

			assert.Same(tt.got, tt.want)

			// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
			if tt.shouldPass && len(mock.errorCalls) != 0 {
				t.Errorf("Same should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
			} else if !tt.shouldPass && len(mock.errorCalls) != 1 {
				t.Errorf("Same should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
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
				// Test GoWise framework behavioral contract
				mock := &behaviorMockT{}
				assert := New(mock)

				assert.True(tt.condition)

				// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
				if tt.shouldPass && len(mock.errorCalls) != 0 {
					t.Errorf("True should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				} else if !tt.shouldPass && len(mock.errorCalls) != 1 {
					t.Errorf("True should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
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
				// Test GoWise framework behavioral contract
				mock := &behaviorMockT{}
				assert := New(mock)

				assert.False(tt.condition)

				// Framework behavior: PASS = no Errorf calls, FAIL = exactly 1 Errorf call
				if tt.shouldPass && len(mock.errorCalls) != 0 {
					t.Errorf("False should pass (no Errorf calls), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				} else if !tt.shouldPass && len(mock.errorCalls) != 1 {
					t.Errorf("False should fail (1 Errorf call), got %d: %v", len(mock.errorCalls), mock.errorCalls)
				}
			})
		}
	})
}
