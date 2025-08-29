// Package custom demonstrates usage of custom domain-specific assertions.
package custom

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDomainAssertions(t *testing.T) {
	t.Run("EmailValidation", func(t *testing.T) {
		assert := NewDomainAssert(t)

		// Valid emails
		assert.IsValidEmail("user@example.com").
			IsValidEmail("test.user+label@domain.co.uk").
			IsValidEmail("123@numbers.org")

		// Domain-specific validation
		assert.HasEmailDomain("alice@company.com", "company.com").
			HasEmailDomain("bob@internal.corp", "internal.corp")
	})

	t.Run("HTTPResponseValidation", func(t *testing.T) {
		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "success",
				"data":    map[string]string{"id": "123"},
			})
		}))
		defer server.Close()

		// Make request and test response
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		assert := NewDomainAssert(t)
		assert.HasStatusCode(resp, http.StatusOK).
			HasContentType(resp, "application/json").
			HasHeader(resp, "Content-Type", "application/json")
	})

	t.Run("JSONValidation", func(t *testing.T) {
		assert := NewDomainAssert(t)

		jsonData := `{
			"user": {
				"id": 123,
				"name": "Alice",
				"profile": {
					"email": "alice@example.com",
					"active": true
				}
			},
			"metadata": {
				"version": "1.0",
				"timestamp": "2024-01-01T00:00:00Z"
			}
		}`

		assert.IsValidJSON(jsonData).
			JSONHasKey(jsonData, "user").
			JSONHasKey(jsonData, "user.profile").
			JSONHasKey(jsonData, "user.profile.email").
			JSONEquals(jsonData, "user", map[string]interface{}{
				"id":   float64(123), // JSON numbers are float64
				"name": "Alice",
				"profile": map[string]interface{}{
					"email":  "alice@example.com",
					"active": true,
				},
			})
	})

	t.Run("URLValidation", func(t *testing.T) {
		assert := NewDomainAssert(t)

		assert.IsValidURL("https://api.example.com/v1/users").
			URLHasScheme("https://api.example.com/v1/users", "https").
			URLHasHost("https://api.example.com/v1/users", "api.example.com")
	})

	t.Run("TimeValidation", func(t *testing.T) {
		assert := NewDomainAssert(t)

		now := time.Now()
		recent := now.Add(-30 * time.Minute)
		future := now.Add(2 * time.Hour)

		// UK timezone for working hours test
		uk, _ := time.LoadLocation("Europe/London")
		workingTime := time.Date(2024, 1, 15, 14, 30, 0, 0, uk) // Monday 2:30 PM

		assert.IsRecentTime(recent, time.Hour).
			IsFutureTime(future).
			IsWorkingHours(workingTime, uk)
	})

	t.Run("StringPatternValidation", func(t *testing.T) {
		assert := NewDomainAssert(t)

		assert.MatchesPattern("user123", `^user\d+$`).
			HasLength("username", 5, 20).
			IsAlphanumeric("user123").
			IsAlphanumeric("ABC123")
	})

	t.Run("BusinessLogicValidation", func(t *testing.T) {
		assert := NewDomainAssert(t)

		assert.IsValidUserAge(25).
			IsValidUserAge(0).
			IsValidUserAge(100).
			IsValidPrice(19.99).
			IsValidPrice(0.01).
			HasValidInventoryCount(50).
			HasValidInventoryCount(0)
	})

	t.Run("CollectionBusinessLogic", func(t *testing.T) {
		assert := NewDomainAssert(t)

		uniqueStrings := []string{"apple", "banana", "cherry"}
		sortedNumbers := []int{1, 5, 10, 15, 20}
		sortedStrings := []string{"apple", "banana", "cherry", "date"}

		assert.HasUniqueElements(uniqueStrings).
			IsSortedAscending(sortedNumbers).
			IsSortedAscending(sortedStrings)
	})
}

func TestComprehensiveDomainValidation(t *testing.T) {
	assert := NewDomainAssert(t)

	t.Run("CompleteUserValidation", func(t *testing.T) {
		user := User{
			ID:       1,
			Username: "alice123",
			Email:    "alice@company.com",
			Age:      28,
			Active:   true,
			Created:  time.Now().Add(-1 * time.Hour),
		}

		assert.IsValidUser(user)
	})

	t.Run("CompleteProductValidation", func(t *testing.T) {
		product := Product{
			ID:          1,
			Name:        "Wireless Headphones",
			Price:       99.99,
			Category:    "Electronics",
			InStock:     50,
			Description: "High-quality wireless headphones with noise cancellation",
		}

		assert.IsValidProduct(product)
	})

	t.Run("CompleteOrderValidation", func(t *testing.T) {
		products := []Product{
			{
				ID:       1,
				Name:     "Product 1",
				Price:    29.99,
				Category: "Category A",
				InStock:  10,
			},
			{
				ID:       2,
				Name:     "Product 2",
				Price:    49.99,
				Category: "Category B",
				InStock:  5,
			},
		}

		order := Order{
			ID:       1001,
			UserID:   123,
			Products: products,
			Total:    79.98,
			Status:   "processing",
			Created:  time.Now().Add(-2 * time.Hour),
		}

		assert.IsValidOrder(order)
	})
}

func TestFailureScenarios(t *testing.T) {
	// These tests demonstrate what happens when custom assertions fail
	// Using a mock testing.T to capture failures without stopping the test

	t.Run("DemonstrateDomainFailures", func(t *testing.T) {
		mockT := &mockTestingT{}
		assert := NewDomainAssert(mockT)

		// Email validation failures
		assert.IsValidEmail("invalid-email")
		assert.HasEmailDomain("user@wrong.com", "expected.com")

		// Business logic failures
		assert.IsValidUserAge(-5)
		assert.IsValidUserAge(200)
		assert.IsValidPrice(-10.50)

		// Collection failures
		duplicateStrings := []string{"apple", "banana", "apple"}
		assert.HasUniqueElements(duplicateStrings)

		unsortedNumbers := []int{1, 5, 3, 10}
		assert.IsSortedAscending(unsortedNumbers)

		// Verify that failures were captured
		if len(mockT.errors) == 0 {
			t.Error("Expected assertion failures to be captured")
		}

		t.Logf("Captured %d assertion failures:", len(mockT.errors))
		for i, err := range mockT.errors {
			t.Logf("  %d: %s", i+1, err)
		}
	})
}

func Example_domainAssertions() {
	// Create a custom assertion context
	mockT := &mockTestingT{}
	assert := NewDomainAssert(mockT)

	// Email validation
	assert.IsValidEmail("user@example.com").
		HasEmailDomain("user@example.com", "example.com")

	// Business object validation
	user := User{
		ID:       1,
		Username: "alice",
		Email:    "alice@company.com",
		Age:      25,
		Active:   true,
		Created:  time.Now(),
	}
	assert.IsValidUser(user)

	// URL validation
	assert.IsValidURL("https://api.example.com/v1/users").
		URLHasScheme("https://api.example.com/v1/users", "https")

	// JSON validation
	jsonData := `{"name": "Alice", "active": true}`
	assert.IsValidJSON(jsonData).
		JSONHasKey(jsonData, "name").
		JSONEquals(jsonData, "active", true)

	// Collection validation
	uniqueItems := []string{"apple", "banana", "cherry"}
	assert.HasUniqueElements(uniqueItems)

	sortedNumbers := []int{1, 2, 3, 4, 5}
	assert.IsSortedAscending(sortedNumbers)

	// Output shows how to use domain-specific assertions
}

// mockTestingT captures test failures for demonstration
type mockTestingT struct {
	errors []string
}

func (m *mockTestingT) Errorf(format string, args ...interface{}) {
	m.errors = append(m.errors, fmt.Sprintf(format, args...))
}

func (m *mockTestingT) FailNow() {
	// Don't actually fail in examples
}

func (m *mockTestingT) Helper() {
	// Optional method
}
