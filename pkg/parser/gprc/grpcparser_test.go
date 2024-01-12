package grpcparser

import (
	"fmt"
	"net/url"
	"testing"
)

func TestParseGRPCMessage(t *testing.T) {
	testCases := []struct {
		input           string
		expectedService string
		expectedMethod  string
		expectedParams  url.Values
		expectError     bool
	}{
		{"grpc.Service/Method", "grpc.Service", "Method", nil, false},
		{"another.Service/AnotherMethod", "another.Service", "AnotherMethod", nil, false},
		{"grpc.Service/Method?param1=value1&param2=value2", "grpc.Service", "Method", url.Values{"param1": []string{"value1"}, "param2": []string{"value2"}}, false},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			fmt.Println("Running test case: ", tc.input)

			grpcInfo, err := ParseGRPCMessage(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			} else if err != nil {
				t.Errorf("Did not expect error but got: %v", err)
				return
			}

			if grpcInfo.ServiceName != tc.expectedService {
				t.Errorf("Expected service: %s, got: %s", tc.expectedService, grpcInfo.ServiceName)
			}

			if grpcInfo.MethodName != tc.expectedMethod {
				t.Errorf("Expected method: %s, got: %s", tc.expectedMethod, grpcInfo.MethodName)
			}

			if !equalParams(grpcInfo.Params, tc.expectedParams) {
				t.Errorf("Expected params: %v, got: %v", tc.expectedParams, grpcInfo.Params)
			}
		})
	}
}

// equalParams checks if two url.Values are equal.
func equalParams(a, b url.Values) bool {
	if len(a) != len(b) {
		return false
	}

	for key := range a {
		if !equalSlice(a[key], b[key]) {
			return false
		}
	}

	return true
}

// equalSlice checks if two slices are equal.
func equalSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
