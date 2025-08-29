// Package basic demonstrates fundamental GoWise assertion usage patterns.
// This example shows how to use core assertions in typical testing scenarios.
package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"gowise/pkg/assertions"
)

// User represents a typical domain model for testing examples.
type User struct {
	ID       int
	Username string
	Email    string
	Active   bool
	Roles    []string
	Metadata map[string]interface{}
}

// mockT implements TestingT interface for demonstration purposes.
type mockT struct {
	errors []string
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func (m *mockT) FailNow() {
	panic("test failed")
}

func (m *mockT) Helper() {} // Optional method

func main() {
	fmt.Println("=== GoWise Basic Usage Examples ===")

	// Example 1: Core Equality Assertions
	fmt.Println("1. Core Equality Assertions")
	demonstrateEquality()
	fmt.Println()

	// Example 2: Nil Checking
	fmt.Println("2. Nil Checking")
	demonstrateNilChecking()
	fmt.Println()

	// Example 3: Boolean Assertions
	fmt.Println("3. Boolean Assertions")
	demonstrateBooleanAssertions()
	fmt.Println()

	// Example 4: Collection Assertions
	fmt.Println("4. Collection Assertions")
	demonstrateCollectionAssertions()
	fmt.Println()

	// Example 5: Error Assertions
	fmt.Println("5. Error Assertions")
	demonstrateErrorAssertions()
	fmt.Println()

	// Example 6: Numeric and Time Assertions
	fmt.Println("6. Numeric and Time Assertions")
	demonstrateNumericAssertions()
	fmt.Println()

	// Example 7: Multiple Assertions
	fmt.Println("7. Multiple Assertions")
	demonstrateMultipleAssertions()
	fmt.Println()
}

func demonstrateEquality() {
	// Create mock testing context for demonstration
	mock := &mockT{}
	assert := assertions.New(mock)

	// Basic equality - these will pass
	fmt.Println("✓ Testing basic equality:")
	assert.Equal(42, 42)
	assert.Equal("hello", "hello")
	assert.Equal(true, true)
	fmt.Printf("  All basic equality tests passed: %t\n", len(mock.errors) == 0)

	// Complex type equality
	user1 := User{ID: 1, Username: "alice", Active: true}
	user2 := User{ID: 1, Username: "alice", Active: true}
	assert.Equal(user1, user2)
	fmt.Printf("  Struct equality test passed: %t\n", len(mock.errors) == 0)

	// NotEqual assertion
	assert.NotEqual(42, 24)
	assert.NotEqual("hello", "world")
	fmt.Printf("  NotEqual tests passed: %t\n", len(mock.errors) == 0)

	// DeepEqual for complex nested structures
	users1 := []User{{ID: 1, Username: "alice"}, {ID: 2, Username: "bob"}}
	users2 := []User{{ID: 1, Username: "alice"}, {ID: 2, Username: "bob"}}
	assert.DeepEqual(users1, users2)
	fmt.Printf("  Deep equality test passed: %t\n", len(mock.errors) == 0)

	// Pointer identity with Same
	value := 42
	ptr1 := &value
	ptr2 := ptr1
	assert.Same(ptr1, ptr2)
	fmt.Printf("  Pointer identity test passed: %t\n", len(mock.errors) == 0)
}

func demonstrateNilChecking() {
	mock := &mockT{}
	assert := assertions.New(mock)

	// Test nil pointers
	var nilPtr *string
	assert.Nil(nilPtr)
	fmt.Printf("✓ Nil pointer test passed: %t\n", len(mock.errors) == 0)

	// Test non-nil pointers
	value := "not nil"
	assert.NotNil(&value)
	fmt.Printf("✓ NotNil pointer test passed: %t\n", len(mock.errors) == 0)

	// Test nil slices
	var nilSlice []int
	assert.Nil(nilSlice)
	fmt.Printf("✓ Nil slice test passed: %t\n", len(mock.errors) == 0)

	// Test empty but not nil slice
	emptySlice := make([]int, 0)
	assert.NotNil(emptySlice)
	fmt.Printf("✓ Empty slice (not nil) test passed: %t\n", len(mock.errors) == 0)

	// Test nil maps
	var nilMap map[string]int
	assert.Nil(nilMap)
	fmt.Printf("✓ Nil map test passed: %t\n", len(mock.errors) == 0)

	// Test nil interfaces (like error)
	var nilError error
	assert.Nil(nilError)
	fmt.Printf("✓ Nil error test passed: %t\n", len(mock.errors) == 0)
}

func demonstrateBooleanAssertions() {
	mock := &mockT{}
	assert := assertions.New(mock)

	// True assertions
	assert.True(2 > 1)
	assert.True(len("hello") == 5)
	assert.True(math.Pi > 3.0)
	fmt.Printf("✓ True assertions passed: %t\n", len(mock.errors) == 0)

	// False assertions
	assert.False(2 < 1)
	assert.False(len("") > 0)
	assert.False(false)
	fmt.Printf("✓ False assertions passed: %t\n", len(mock.errors) == 0)

	// Boolean expressions with variables
	user := User{Active: true, Username: "alice"}
	assert.True(user.Active)
	assert.True(len(user.Username) > 0)
	fmt.Printf("✓ Boolean expressions with structs passed: %t\n", len(mock.errors) == 0)
}

func demonstrateCollectionAssertions() {
	mock := &mockT{}
	assert := assertions.New(mock)

	// Length assertions
	assert.Len("hello", 5)
	assert.Len([]int{1, 2, 3}, 3)
	assert.Len(map[string]int{"a": 1, "b": 2}, 2)
	fmt.Printf("✓ Length assertions passed: %t\n", len(mock.errors) == 0)

	// Contains assertions for strings
	assert.Contains("hello world", "world")
	assert.Contains("testing", "test")
	fmt.Printf("✓ String contains assertions passed: %t\n", len(mock.errors) == 0)

	// Contains assertions for slices
	numbers := []int{1, 2, 3, 4, 5}
	assert.Contains(numbers, 3)
	assert.Contains(numbers, 1)
	fmt.Printf("✓ Slice contains assertions passed: %t\n", len(mock.errors) == 0)

	// Contains assertions for maps (checks keys)
	userMap := map[string]User{
		"alice": {ID: 1, Username: "alice"},
		"bob":   {ID: 2, Username: "bob"},
	}
	assert.Contains(userMap, "alice")
	assert.Contains(userMap, "bob")
	fmt.Printf("✓ Map contains assertions passed: %t\n", len(mock.errors) == 0)
}

func demonstrateErrorAssertions() {
	mock := &mockT{}
	assert := assertions.New(mock)

	// NoError assertion
	successOperation := func() error { return nil }
	err := successOperation()
	assert.NoError(err)
	fmt.Printf("✓ NoError assertion passed: %t\n", len(mock.errors) == 0)

	// HasError assertion
	failingOperation := func() error { return fmt.Errorf("operation failed") }
	err = failingOperation()
	assert.HasError(err)
	fmt.Printf("✓ HasError assertion passed: %t\n", len(mock.errors) == 0)

	// ErrorContains assertion
	specificErr := fmt.Errorf("database connection failed: timeout after 30s")
	assert.ErrorContains(specificErr, "connection failed")
	assert.ErrorContains(specificErr, "timeout")
	fmt.Printf("✓ ErrorContains assertions passed: %t\n", len(mock.errors) == 0)

	// ErrorMatches with regex
	assert.ErrorMatches(specificErr, `database .* failed`)
	assert.ErrorMatches(specificErr, `timeout after \d+s`)
	fmt.Printf("✓ ErrorMatches assertions passed: %t\n", len(mock.errors) == 0)
}

func demonstrateNumericAssertions() {
	mock := &mockT{}
	assert := assertions.New(mock)

	// InDelta for floating-point comparisons
	assert.InDelta(math.Pi, 3.14159, 0.0001)
	assert.InDelta(1.0/3.0, 0.3333, 0.0001)
	fmt.Printf("✓ InDelta assertions passed: %t\n", len(mock.errors) == 0)

	// WithinDuration for time comparisons
	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	end := time.Now()
	assert.IsWithinDuration(end, start, 50*time.Millisecond)
	fmt.Printf("✓ WithinDuration assertion passed: %t\n", len(mock.errors) == 0)
}

func demonstrateMultipleAssertions() {
	mock := &mockT{}
	assert := assertions.New(mock)

	// Multiple related assertions
	user := User{
		ID:       123,
		Username: "alice",
		Email:    "alice@example.com",
		Active:   true,
		Roles:    []string{"user", "admin"},
		Metadata: map[string]interface{}{
			"created": time.Now(),
			"source":  "api",
		},
	}

	// Multiple assertions for comprehensive validation
	assert.Equal(user.ID, 123)
	assert.Equal(user.Username, "alice")
	assert.True(user.Active)
	assert.Contains(user.Email, "@")
	assert.Len(user.Roles, 2)
	assert.Contains(user.Roles, "admin")
	assert.Contains(user.Metadata, "source")

	fmt.Printf("✓ Multiple assertions test passed: %t\n", len(mock.errors) == 0)

	// Demonstrate comprehensive validation
	users := []User{
		{ID: 1, Username: "alice", Active: true},
		{ID: 2, Username: "bob", Active: false},
		{ID: 3, Username: "charlie", Active: true},
	}

	assert.Len(users, 3)
	assert.True(users[0].Active)
	assert.False(users[1].Active)
	assert.Equal(users[2].Username, "charlie")

	fmt.Printf("✓ Comprehensive validation passed: %t\n", len(mock.errors) == 0)

	if len(mock.errors) == 0 {
		fmt.Println("🎉 All assertions completed successfully!")
	} else {
		fmt.Printf("❌ Some assertions failed: %v\n", mock.errors)
	}
}

// simulateFailure demonstrates what happens when assertions fail.
func simulateFailure() {
	fmt.Println("\n=== Demonstrating Assertion Failures ===")

	mock := &mockT{}
	assert := assertions.New(mock)

	// These assertions will fail and generate error messages
	assert.Equal(42, 24)                // Numbers don't match
	assert.True(false)                  // Boolean is false
	assert.Contains("hello", "goodbye") // String doesn't contain substring
	assert.Len([]int{1, 2, 3}, 5)       // Wrong length

	fmt.Printf("Failed assertions captured: %d\n", len(mock.errors))
	for i, err := range mock.errors {
		fmt.Printf("Error %d: %s\n", i+1, err)
	}
}

func init() {
	log.SetPrefix("GoWise Example: ")
	log.SetFlags(log.Ltime | log.Lshortfile)
}
