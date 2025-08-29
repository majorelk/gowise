// Package custom demonstrates how to create domain-specific custom assertions
// that extend GoWise's core functionality for specific business domains.
package custom

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"gowise/pkg/assertions"
)

// DomainAssert extends the base Assert with domain-specific assertions
type DomainAssert struct {
	*assertions.Assert
	t assertions.TestingT
}

// NewDomainAssert creates a new domain-specific assertion context
func NewDomainAssert(t assertions.TestingT) *DomainAssert {
	return &DomainAssert{
		Assert: assertions.New(t),
		t:      t,
	}
}

// Email Validation Assertions
func (a *DomainAssert) IsValidEmail(email string) *DomainAssert {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		a.t.Errorf("IsValidEmail: invalid email format\n  got: %q\n  expected: valid email format", email)
	}
	return a
}

func (a *DomainAssert) HasEmailDomain(email, expectedDomain string) *DomainAssert {
	if !strings.Contains(email, "@") {
		a.t.Errorf("HasEmailDomain: invalid email format: %q", email)
		return a
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		a.t.Errorf("HasEmailDomain: invalid email format: %q", email)
		return a
	}

	domain := parts[1]
	if domain != expectedDomain {
		a.t.Errorf("HasEmailDomain: wrong email domain\n  got domain: %q\n  want domain: %q\n  full email: %q",
			domain, expectedDomain, email)
	}
	return a
}

// HTTP Response Assertions
func (a *DomainAssert) HasStatusCode(response *http.Response, expectedStatus int) *DomainAssert {
	if response.StatusCode != expectedStatus {
		a.t.Errorf("HasStatusCode: wrong HTTP status\n  got: %d (%s)\n  want: %d (%s)",
			response.StatusCode, http.StatusText(response.StatusCode),
			expectedStatus, http.StatusText(expectedStatus))
	}
	return a
}

func (a *DomainAssert) HasHeader(response *http.Response, headerName, expectedValue string) *DomainAssert {
	actualValue := response.Header.Get(headerName)
	if actualValue != expectedValue {
		a.t.Errorf("HasHeader: wrong header value\n  header: %q\n  got: %q\n  want: %q",
			headerName, actualValue, expectedValue)
	}
	return a
}

func (a *DomainAssert) HasContentType(response *http.Response, expectedContentType string) *DomainAssert {
	contentType := response.Header.Get("Content-Type")

	// Handle cases where content type might include charset
	if !strings.HasPrefix(contentType, expectedContentType) {
		a.t.Errorf("HasContentType: wrong content type\n  got: %q\n  want: %q (or with charset)",
			contentType, expectedContentType)
	}
	return a
}

// JSON Response Assertions
func (a *DomainAssert) IsValidJSON(data string) *DomainAssert {
	var js interface{}
	if err := json.Unmarshal([]byte(data), &js); err != nil {
		a.t.Errorf("IsValidJSON: invalid JSON format\n  error: %v\n  data: %q", err, data)
	}
	return a
}

func (a *DomainAssert) JSONHasKey(jsonData string, keyPath string) *DomainAssert {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		a.t.Errorf("JSONHasKey: invalid JSON format\n  error: %v", err)
		return a
	}

	// Simple key path traversal (supports dot notation like "user.profile.name")
	keys := strings.Split(keyPath, ".")
	current := data

	for i, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			if val, exists := v[key]; exists {
				current = val
			} else {
				a.t.Errorf("JSONHasKey: key not found\n  key path: %q\n  missing key: %q (at position %d)\n  available keys: %v",
					keyPath, key, i, getMapKeys(v))
				return a
			}
		default:
			a.t.Errorf("JSONHasKey: cannot traverse key path\n  key path: %q\n  stopped at: %q\n  current value type: %T",
				keyPath, key, current)
			return a
		}
	}

	return a
}

func (a *DomainAssert) JSONEquals(jsonData string, key string, expectedValue interface{}) *DomainAssert {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		a.t.Errorf("JSONEquals: invalid JSON format\n  error: %v", err)
		return a
	}

	actualValue, exists := data[key]
	if !exists {
		a.t.Errorf("JSONEquals: key not found\n  key: %q\n  available keys: %v", key, getMapKeys(data))
		return a
	}

	if !reflect.DeepEqual(actualValue, expectedValue) {
		a.t.Errorf("JSONEquals: wrong value for key\n  key: %q\n  got: %v (%T)\n  want: %v (%T)",
			key, actualValue, actualValue, expectedValue, expectedValue)
	}

	return a
}

// URL Validation Assertions
func (a *DomainAssert) IsValidURL(rawURL string) *DomainAssert {
	if _, err := url.Parse(rawURL); err != nil {
		a.t.Errorf("IsValidURL: invalid URL format\n  URL: %q\n  error: %v", rawURL, err)
	}
	return a
}

func (a *DomainAssert) URLHasScheme(rawURL, expectedScheme string) *DomainAssert {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		a.t.Errorf("URLHasScheme: invalid URL format\n  URL: %q\n  error: %v", rawURL, err)
		return a
	}

	if parsed.Scheme != expectedScheme {
		a.t.Errorf("URLHasScheme: wrong URL scheme\n  URL: %q\n  got scheme: %q\n  want scheme: %q",
			rawURL, parsed.Scheme, expectedScheme)
	}
	return a
}

func (a *DomainAssert) URLHasHost(rawURL, expectedHost string) *DomainAssert {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		a.t.Errorf("URLHasHost: invalid URL format\n  URL: %q\n  error: %v", rawURL, err)
		return a
	}

	if parsed.Host != expectedHost {
		a.t.Errorf("URLHasHost: wrong URL host\n  URL: %q\n  got host: %q\n  want host: %q",
			rawURL, parsed.Host, expectedHost)
	}
	return a
}

// Time Assertions
func (a *DomainAssert) IsRecentTime(timestamp time.Time, maxAge time.Duration) *DomainAssert {
	age := time.Since(timestamp)
	if age > maxAge {
		a.t.Errorf("IsRecentTime: timestamp too old\n  timestamp: %v\n  age: %v\n  max age: %v",
			timestamp, age, maxAge)
	}
	return a
}

func (a *DomainAssert) IsFutureTime(timestamp time.Time) *DomainAssert {
	if !timestamp.After(time.Now()) {
		a.t.Errorf("IsFutureTime: timestamp is not in the future\n  timestamp: %v\n  current time: %v",
			timestamp, time.Now())
	}
	return a
}

func (a *DomainAssert) IsWorkingHours(timestamp time.Time, timezone *time.Location) *DomainAssert {
	localTime := timestamp.In(timezone)
	hour := localTime.Hour()
	weekday := localTime.Weekday()

	// Define working hours (9 AM to 5 PM, Monday to Friday)
	isWeekday := weekday >= time.Monday && weekday <= time.Friday
	isWorkingHour := hour >= 9 && hour < 17

	if !isWeekday || !isWorkingHour {
		a.t.Errorf("IsWorkingHours: timestamp is outside working hours\n  timestamp: %v\n  local time: %v\n  hour: %d\n  weekday: %v",
			timestamp, localTime, hour, weekday)
	}
	return a
}

// String Pattern Assertions
func (a *DomainAssert) MatchesPattern(str, pattern string) *DomainAssert {
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		a.t.Errorf("MatchesPattern: invalid regex pattern\n  pattern: %q\n  error: %v", pattern, err)
		return a
	}

	if !matched {
		a.t.Errorf("MatchesPattern: string does not match pattern\n  string: %q\n  pattern: %q", str, pattern)
	}
	return a
}

func (a *DomainAssert) HasLength(str string, min, max int) *DomainAssert {
	length := len(str)
	if length < min || length > max {
		a.t.Errorf("HasLength: string length outside valid range\n  string: %q\n  length: %d\n  valid range: %d-%d",
			str, length, min, max)
	}
	return a
}

func (a *DomainAssert) IsAlphanumeric(str string) *DomainAssert {
	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumeric.MatchString(str) {
		a.t.Errorf("IsAlphanumeric: string contains non-alphanumeric characters\n  string: %q", str)
	}
	return a
}

// Business Logic Assertions
func (a *DomainAssert) IsValidUserAge(age int) *DomainAssert {
	if age < 0 || age > 150 {
		a.t.Errorf("IsValidUserAge: age outside realistic range\n  age: %d\n  valid range: 0-150", age)
	}
	return a
}

func (a *DomainAssert) IsValidPrice(price float64) *DomainAssert {
	if price < 0 {
		a.t.Errorf("IsValidPrice: price cannot be negative\n  price: %.2f", price)
	}
	if price > 1000000 {
		a.t.Errorf("IsValidPrice: price exceeds maximum allowed\n  price: %.2f\n  maximum: 1,000,000", price)
	}
	return a
}

func (a *DomainAssert) HasValidInventoryCount(count int) *DomainAssert {
	if count < 0 {
		a.t.Errorf("HasValidInventoryCount: inventory count cannot be negative\n  count: %d", count)
	}
	return a
}

// Collection Business Logic
func (a *DomainAssert) HasUniqueElements(slice interface{}) *DomainAssert {
	// This is a simplified version - in real implementation you'd use reflection
	// to handle different slice types properly
	switch s := slice.(type) {
	case []string:
		seen := make(map[string]bool)
		duplicates := []string{}

		for _, item := range s {
			if seen[item] {
				duplicates = append(duplicates, item)
			} else {
				seen[item] = true
			}
		}

		if len(duplicates) > 0 {
			a.t.Errorf("HasUniqueElements: found duplicate elements\n  duplicates: %v", duplicates)
		}

	case []int:
		seen := make(map[int]bool)
		duplicates := []int{}

		for _, item := range s {
			if seen[item] {
				duplicates = append(duplicates, item)
			} else {
				seen[item] = true
			}
		}

		if len(duplicates) > 0 {
			a.t.Errorf("HasUniqueElements: found duplicate elements\n  duplicates: %v", duplicates)
		}

	default:
		a.t.Errorf("HasUniqueElements: unsupported type %T", slice)
	}

	return a
}

func (a *DomainAssert) IsSortedAscending(slice interface{}) *DomainAssert {
	switch s := slice.(type) {
	case []int:
		for i := 1; i < len(s); i++ {
			if s[i] < s[i-1] {
				a.t.Errorf("IsSortedAscending: slice not sorted in ascending order\n  position %d: %d > %d",
					i, s[i-1], s[i])
				return a
			}
		}
	case []string:
		for i := 1; i < len(s); i++ {
			if s[i] < s[i-1] {
				a.t.Errorf("IsSortedAscending: slice not sorted in ascending order\n  position %d: %q > %q",
					i, s[i-1], s[i])
				return a
			}
		}
	default:
		a.t.Errorf("IsSortedAscending: unsupported type %T", slice)
	}

	return a
}

// Helper functions
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Example domain models for testing
type User struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Age      int       `json:"age"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	InStock     int     `json:"in_stock"`
	Description string  `json:"description"`
}

type Order struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Products []Product `json:"products"`
	Total    float64   `json:"total"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created"`
}

// Comprehensive domain assertion for complex business objects
func (a *DomainAssert) IsValidUser(user User) *DomainAssert {
	a.True(user.ID > 0)
	a.HasLength(user.Username, 3, 50)
	a.IsAlphanumeric(user.Username)
	a.IsValidEmail(user.Email)
	a.IsValidUserAge(user.Age)
	a.IsRecentTime(user.Created, 24*time.Hour)
	return a
}

func (a *DomainAssert) IsValidProduct(product Product) *DomainAssert {
	a.True(product.ID > 0)
	a.HasLength(product.Name, 1, 200)
	a.IsValidPrice(product.Price)
	a.HasLength(product.Category, 1, 50)
	a.HasValidInventoryCount(product.InStock)
	a.HasLength(product.Description, 0, 1000) // Optional field
	return a
}

func (a *DomainAssert) IsValidOrder(order Order) *DomainAssert {
	// Validate basic order properties
	a.True(order.ID > 0)
	a.True(order.UserID > 0)
	a.True(len(order.Products) > 0)
	a.IsValidPrice(order.Total)
	a.Contains([]string{"pending", "processing", "completed", "cancelled"}, order.Status)
	a.IsRecentTime(order.Created, 30*24*time.Hour) // Within last 30 days

	// Validate all products in the order
	for _, product := range order.Products {
		a.IsValidProduct(product)
	}

	return a
}
