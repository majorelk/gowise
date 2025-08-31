// Package fluent-chaining demonstrates the new method chaining functionality in GoWise.
// This example shows how fluent chaining enables more readable and concise test assertions.
package main

import (
	"fmt"
	"gowise/pkg/assertions"
)

// User represents a domain model for demonstration.
type User struct {
	ID       int
	Username string
	Email    string
	Active   bool
	Roles    []string
}

// mockT implements TestingT interface for demonstration.
type mockT struct {
	errors []string
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func (m *mockT) FailNow() {
	// In real testing, this would stop test execution
	fmt.Println("Test would fail here")
}

func (m *mockT) Helper() {}

func main() {
	fmt.Println("=== GoWise Fluent Chaining Examples ===\n")

	// Example 1: Traditional approach vs Fluent chaining
	fmt.Println("1. Before and After Comparison:")
	fmt.Println("   Before (individual calls):")
	fmt.Println("   assert.Equal(user.ID, 123)")
	fmt.Println("   assert.True(user.Active)")
	fmt.Println("   assert.Contains(user.Email, \"@\")")
	fmt.Println("")
	fmt.Println("   After (fluent chaining):")
	fmt.Println("   assert.Equal(user.ID, 123).True(user.Active).Contains(user.Email, \"@\")")
	fmt.Println("")

	// Example 2: Successful chaining
	fmt.Println("2. Successful Chaining:")
	user := User{
		ID:       123,
		Username: "john_doe",
		Email:    "john@example.com",
		Active:   true,
		Roles:    []string{"user", "admin"},
	}

	mock1 := &mockT{}
	assert1 := assertions.New(mock1)

	result1 := assert1.
		Equal(user.ID, 123).
		True(user.Active).
		Contains(user.Email, "@").
		Len(user.Roles, 2).
		Contains(user.Username, "john")

	fmt.Printf("   Chain succeeded: %v\n", !result1.HasFailed())
	fmt.Printf("   Error: %q\n\n", result1.Error())

	// Example 3: Fail-fast behavior
	fmt.Println("3. Fail-Fast Behavior:")
	mock2 := &mockT{}
	assert2 := assertions.New(mock2)

	result2 := assert2.
		Equal(user.ID, 999).        // This will fail
		True(false).                // This becomes no-op due to fail-fast
		Contains(user.Email, "xyz") // This also becomes no-op

	fmt.Printf("   Chain failed: %v\n", result2.HasFailed())
	fmt.Printf("   Error (only first failure): %q\n\n", result2.Error())

	// Example 4: Complex chaining with different assertion types
	fmt.Println("4. Complex Mixed Assertions:")
	var err error
	data := map[string]interface{}{
		"count": 42,
		"items": []string{"a", "b", "c"},
		"user":  &user,
	}

	mock3 := &mockT{}
	assert3 := assertions.New(mock3)

	result3 := assert3.
		NoError(err).                 // Error assertion
		NotNil(data["user"]).         // Nil assertion
		Contains(data, "count").      // Map contains
		Len(data["items"], 3).        // Length assertion
		True(data["count"].(int) > 0) // Boolean assertion

	fmt.Printf("   Complex chain succeeded: %v\n", !result3.HasFailed())
	fmt.Printf("   Error: %q\n\n", result3.Error())

	// Example 5: Error handling chain
	fmt.Println("5. Error Handling Chain:")
	mock4 := &mockT{}
	assert4 := assertions.New(mock4)

	result4 := assert4.
		NoError(err).
		Len(user.Roles, 2).
		HasError(fmt.Errorf("expected error")) // This will fail

	fmt.Printf("   Error handling failed as expected: %v\n", result4.HasFailed())
	fmt.Printf("   Error: %q\n\n", result4.Error())

	fmt.Println("=== Benefits of Fluent Chaining ===")
	fmt.Println("✓ More readable - assertions read like natural language")
	fmt.Println("✓ More concise - less repetition of 'assert' variable")
	fmt.Println("✓ Fail-fast - stops at first failure, preserving context")
	fmt.Println("✓ Backward compatible - existing code continues to work")
	fmt.Println("✓ Type-safe - all methods return *Assert for consistent chaining")
}

