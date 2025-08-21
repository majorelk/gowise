// Package assertions provides assertion functions for testing.
//
// Assertions is a utility package that offers functions to perform common
// assertion checks in tests.
//
// Example:
//
//	assert := assertions.New(t)
//	assert.True(true)
//	if assert.Error() != "" {
//	    t.Errorf("Expected no error, but got: %s", assert.Error())
//	}
package assertions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Assert is a struct that holds the testing context and error message.
type Assert struct {
	t        interface{}
	errorMsg string
}

// New creates a new Assert instance with the given testing context.
func New(t interface{}) *Assert {
	return &Assert{t: t}
}

// Equal asserts that two values are equal.
func (a *Assert) Equal(expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		a.reportError(expected, actual, "expected to be equal")
	}
}

// NotEqual asserts that two values are not equal.
func (a *Assert) NotEqual(expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		a.reportError(expected, actual, "expected to be not equal")
	}
}

// True asserts that a value is true.
func (a *Assert) True(value bool) {
	if !value {
		a.reportError(true, value, "expected to be true")
	}
}

// False asserts that a value is false.
func (a *Assert) False(value bool) {
	if value {
		a.reportError(false, value, "expected to be false")
	}
}

// reportError is a helper function to report test failures.
func (a *Assert) reportError(expected, actual interface{}, message string) {
	a.errorMsg = fmt.Sprintf("%s - Expected %v, Actual: %v", message, expected, actual)
}

// Error returns the error message if the assertion failed.
func (a *Assert) Error() string {
	return a.errorMsg
}

// Nil asserts that a value is nil.
func (a *Assert) Nil(value interface{}) {
	if !isNil(value) {
		a.reportError(nil, value, "expected to be nil")
	}
}

// NotNil asserts that a value is not nil.
func (a *Assert) NotNil(value interface{}) {
	if isNil(value) {
		a.reportError("not nil", value, "expected to be not nil")
	}
}

// isNil is a helper function to check if a value is nil.
func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	valueType := reflect.TypeOf(value)
	if valueType.Kind() == reflect.Ptr || valueType.Kind() == reflect.Slice || valueType.Kind() == reflect.Map {
		return reflect.ValueOf(value).IsNil()
	}

	return false
}

// Contains asserts that a slice, array, or string contains a certain element or substring.
func (a *Assert) Contains(container, item interface{}) {
	containerValue := reflect.ValueOf(container)

	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerValue.Len(); i++ {
			if reflect.DeepEqual(containerValue.Index(i).Interface(), item) {
				return
			}
		}
	case reflect.String:
		if strings.Contains(container.(string), item.(string)) {
			return
		}
	}

	a.reportError(container, item, "expected to contain")
}

// Greater asserts that the first value is greater than the second.
func (a *Assert) Greater(v1, v2 float64) {
	if v1 <= v2 {
		a.reportError(v1, v2, "expected to be greater")
	}
}

// Less asserts that the first value is less than the second.
func (a *Assert) Less(v1, v2 float64) {
	if v1 >= v2 {
		a.reportError(v1, v2, "expected to be less")
	}
}

// HasPrefix asserts that a string starts with a certain substring.
func (a *Assert) HasPrefix(s, prefix string) {
	if !strings.HasPrefix(s, prefix) {
		a.reportError(prefix, s, "expected to have prefix")
	}
}

// HasSuffix asserts that a string ends with a certain substring.
func (a *Assert) HasSuffix(s, suffix string) {
	if !strings.HasSuffix(s, suffix) {
		a.reportError(suffix, s, "expected to have suffix")
	}
}

// InDelta asserts that the difference between two numeric values is within a certain range.
func (a *Assert) InDelta(expected, actual, delta float64) {
	if math.Abs(expected-actual) > delta {
		a.reportError(expected, actual, fmt.Sprintf("expected difference to be within %v", delta))
	}
}

// InEpsilon asserts that the difference between two numeric values is within a certain percentage.
func (a *Assert) InEpsilon(expected, actual, epsilon float64) {
	if math.Abs((expected-actual)/((expected+actual)/2)) > epsilon {
		a.reportError(expected, actual, fmt.Sprintf("expected difference to be within %v percent", epsilon*100))
	}
}

// Regexp asserts that a string matches a regular expression.
func (a *Assert) Regexp(pattern, str string) {
	if matched, _ := regexp.MatchString(pattern, str); !matched {
		a.reportError(pattern, str, "expected to match regular expression")
	}
}

// NoError asserts that a function call returns no error.
func (a *Assert) NoError(err error) {
	if err != nil {
		a.reportError("no error", err, "expected no error")
	}
}

// ErrorType asserts that a function call returns a specific error type.
func (a *Assert) ErrorType(expected, actual error) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		a.reportError(expected, actual, "expected error of different type")
	}
}

// IsEmpty asserts that a given array, slice, map, or string is empty.
func (a *Assert) IsEmpty(value interface{}) {
	valueType := reflect.TypeOf(value)
	switch valueType.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		if reflect.ValueOf(value).Len() != 0 {
			a.reportError(0, reflect.ValueOf(value).Len(), "expected to be empty")
		}
	default:
		a.reportError(nil, value, "invalid type for IsEmpty")
	}
}

// IsNotEmpty asserts that a given array, slice, map, or string is not empty.
func (a *Assert) IsNotEmpty(value interface{}) {
	valueType := reflect.TypeOf(value)
	switch valueType.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		if reflect.ValueOf(value).Len() == 0 {
			a.reportError("not empty", reflect.ValueOf(value).Len(), "expected to be not empty")
		}
	default:
		a.reportError(nil, value, "invalid type for IsNotEmpty")
	}
}

// Len asserts that a given array, slice, map, or string has a specific length.
func (a *Assert) Len(value interface{}, expectedLen int) {
	valueType := reflect.TypeOf(value)
	switch valueType.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		if reflect.ValueOf(value).Len() != expectedLen {
			a.reportError(expectedLen, reflect.ValueOf(value).Len(), "expected different length")
		}
	default:
		a.reportError(nil, value, "invalid type for Len")
	}
}

// Implements asserts that an object implements a certain interface.
func (a *Assert) Implements(object, interfaceObj interface{}) {
	objectType := reflect.TypeOf(object)
	interfaceType := reflect.TypeOf(interfaceObj).Elem()
	if !objectType.Implements(interfaceType) {
		a.reportError(interfaceType, objectType, "expected to implement interface")
	}
}

// IsZero asserts that a given numeric value or time.Time is zero.
func (a *Assert) IsZero(value interface{}) {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		if v != 0 {
			a.reportError(0, v, "expected to be zero")
		}
	case uint, uint8, uint16, uint32, uint64:
		if v != 0 {
			a.reportError(0, v, "expected to be zero")
		}
	case float32, float64:
		if v != 0 {
			a.reportError(0, v, "expected to be zero")
		}
	case time.Time:
		if !v.IsZero() {
			a.reportError(time.Time{}, v, "expected to be zero")
		}
	default:
		a.reportError(nil, value, "invalid type for IsZero")
	}
}

// IsNotZero asserts that a given numeric value or time.Time is not zero.
func (a *Assert) IsNotZero(value interface{}) {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		if v == 0 {
			a.reportError("not zero", v, "expected to be not zero")
		}
	case uint, uint8, uint16, uint32, uint64:
		if v == 0 {
			a.reportError("not zero", v, "expected to be not zero")
		}
	case float32, float64:
		if v == 0 {
			a.reportError("not zero", v, "expected to be not zero")
		}
	case time.Time:
		if v.IsZero() {
			a.reportError("not zero", v, "expected to be not zero")
		}
	default:
		a.reportError(nil, value, "invalid type for IsNotZero")
	}
}

// IsWithinDuration asserts that a given time.Time is within a certain duration from another time.Time.
func (a *Assert) IsWithinDuration(t1, t2 time.Time, d time.Duration) {
	if abs := t1.Sub(t2); abs > d {
		a.reportError(d, abs, "expected to be within duration")
	}
}

// MatchesPattern asserts that a string matches a certain pattern.
func (a *Assert) MatchesPattern(pattern, s string) {
	if matched, _ := regexp.MatchString(pattern, s); !matched {
		a.reportError(pattern, s, "expected to match pattern")
	}
}

// Panics asserts that a certain function panics.
func (a *Assert) Panics(f func()) {
	defer func() {
		if r := recover(); r == nil {
			a.reportError("panic", nil, "expected to panic")
		}
	}()

	f()
}

// PanicsWith asserts that a certain function panics with a specific value.
func (a *Assert) PanicsWith(f func(), expected interface{}) {
	defer func() {
		if r := recover(); r != expected {
			a.reportError(expected, r, "expected to panic with")
		}
	}()

	f()
}

// NotPanics asserts that a certain function does not panic.
func (a *Assert) NotPanics(f func()) {
	defer func() {
		if r := recover(); r != nil {
			a.reportError("no panic", r, "expected not to panic")
		}
	}()

	f()
}

// Condition asserts that a certain condition is true.	// Condition asserts that a certain condition is true.
func (a *Assert) Condition(condition bool) {
	if !condition {
		a.reportError(true, condition, "expected condition to be true")
	}
}

// Conditionf asserts that a certain condition is true with a formatted message.
func (a *Assert) Conditionf(condition bool, format string, args ...interface{}) {
	if !condition {
		a.reportError(true, condition, fmt.Sprintf(format, args...))
	}
}

// HttpStatus asserts that a HTTP response has a certain status code.
func (a *Assert) HttpStatus(response *http.Response, expected int) {
	if response.StatusCode != expected {
		a.reportError(expected, response.StatusCode, "expected different HTTP status")
	}
}

// JsonEqual asserts that a JSON string is equivalent to another JSON string or object.
func (a *Assert) JsonEqual(expected, actual string) {
	var obj1, obj2 interface{}
	json.Unmarshal([]byte(expected), &obj1)
	json.Unmarshal([]byte(actual), &obj2)
	a.Equal(obj1, obj2)
}

// HasHeader asserts that a HTTP response has a certain header.
func (a *Assert) HasHeader(response *http.Response, header string) {
	if _, ok := response.Header[header]; !ok {
		a.reportError(header, nil, "expected to have header")
	}
}

// HeaderEqual asserts that a HTTP response header has a certain value.
func (a *Assert) HeaderEqual(response *http.Response, header, expected string) {
	if value, ok := response.Header[header]; !ok || len(value) == 0 || value[0] != expected {
		a.reportError(expected, response.Header[header], "expected different header value")
	}
}

// BodyContains asserts that a HTTP response body contains a certain string.
func (a *Assert) BodyContains(response *http.Response, expected string) {
	body, _ := ioutil.ReadAll(response.Body)
	if !strings.Contains(string(body), expected) {
		a.reportError(expected, string(body), "expected body to contain")
	}
}

// BodyMatches asserts that a HTTP response body matches a certain regular expression.
func (a *Assert) BodyMatches(response *http.Response, pattern string) {
	body, _ := ioutil.ReadAll(response.Body)
	if matched, _ := regexp.MatchString(pattern, string(body)); !matched {
		a.reportError(pattern, string(body), "expected body to match")
	}
}

// IsRedirect asserts that a HTTP response is a redirect.
func (a *Assert) IsRedirect(response *http.Response) {
	if response.StatusCode < 300 || response.StatusCode >= 400 {
		a.reportError("3xx", response.StatusCode, "expected redirect status code")
	}
}

// IsSuccess asserts that a HTTP response is a success.
func (a *Assert) IsSuccess(response *http.Response) {
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		a.reportError("2xx", response.StatusCode, "expected success status code")
	}
}

// IsClientError asserts that a HTTP response is a client error.
func (a *Assert) IsClientError(response *http.Response) {
	if response.StatusCode < 400 || response.StatusCode >= 500 {
		a.reportError("4xx", response.StatusCode, "expected client error status code")
	}
}

// IsServerError asserts that a HTTP response is a server error.
func (a *Assert) IsServerError(response *http.Response) {
	if response.StatusCode < 500 || response.StatusCode >= 600 {
		a.reportError("5xx", response.StatusCode, "expected server error status code")
	}
}

// BodyJsonEqual asserts that a HTTP response body is equivalent to a given JSON object.
func (a *Assert) BodyJsonEqual(response *http.Response, expected interface{}) {
	body, _ := ioutil.ReadAll(response.Body)
	var actual interface{}
	json.Unmarshal(body, &actual)
	a.Equal(expected, actual)
}

// HasCookie asserts that a HTTP response has a certain cookie.
func (a *Assert) HasCookie(response *http.Response, name string) {
	var hasCookie bool
	for _, cookie := range response.Cookies() {
		if cookie.Name == name {
			hasCookie = true
			break
		}
	}
	if !hasCookie {
		a.reportError(name, nil, "expected to have cookie")
	}
}

// HeaderContains asserts that a HTTP response header contains a certain value.
func (a *Assert) HeaderContains(response *http.Response, header, expected string) {
	values, ok := response.Header[header]
	if !ok {
		a.reportError(expected, nil, "expected to have header")
	} else {
		var found bool
		for _, value := range values {
			if value == expected {
				found = true
				break
			}
		}
		if !found {
			a.reportError(expected, values, "expected header to contain value")
		}
	}
}

func (a *Assert) ResponseTime(url string, maxTime time.Duration) {
	start := time.Now()
	_, err := http.Get(url)
	if err != nil {
		a.reportError(maxTime, err, "error making request")
	}
	elapsed := time.Since(start)
	if elapsed > maxTime {
		a.reportError(maxTime, elapsed, "response time exceeded")
	}
}

func (a *Assert) IsSorted(slice []int) {
	if !sort.IntsAreSorted(slice) {
		a.reportError(nil, slice, "slice is not sorted")
	}
}

func (a *Assert) IsSortedFloat64(slice []float64) {
	if !sort.Float64sAreSorted(slice) {
		a.reportError(nil, slice, "slice is not sorted")
	}
}

func (a *Assert) FileExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		a.reportError(path, nil, "file does not exist")
	}
}

func (a *Assert) DirectoryExists(path string) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) || !info.IsDir() {
		a.reportError(path, nil, "directory does not exist")
	}
}
