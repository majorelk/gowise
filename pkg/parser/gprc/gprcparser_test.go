package grpcparser

import (
	"testing"
	"strings"
	"fmt"
)

// getMethodNameWithoutParams extracts the method name without parameters.
func getMethodNameWithoutParams(methodWithParams string) string {
	parts := strings.Split(methodWithParams, "?")
	return parts[0]
}

func TestParseGRPCMessage(t *testing.T) {
	testCases := []struct {
		input           string
		expectedService string
		expectedMethod  string
		expectedParams  map[string]string
	}{
		{"grpc.Service/Method", "grpc.Service", "Method", nil},
		{"another.Service/AnotherMethod", "another.Service", "AnotherMethod", nil},
		{"grpc.Service/Method?param1=value1&param2=value2", "grpc.Service", "Method", map[string]string{"param1": "value1", "param2": "value2"}},
		// {"grpc.Service/Method?param1=value1&param2=value2", "grpc.Service", "Method", map[string]string{"param1": "value1", "param2": "value2"}},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			fmt.Println("Running test case: ", tc.input)

			grpcInfo := ParseGRPCMessage(tc.input)

			if grpcInfo.ServiceName != tc.expectedService {
				t.Errorf("Expected service: %s, got: %s", tc.expectedService, grpcInfo.ServiceName)
			}

			expectedMethodWithoutParams := getMethodNameWithoutParams(tc.expectedMethod)
			actualMethodWithoutParams := getMethodNameWithoutParams(grpcInfo.MethodName)

			if actualMethodWithoutParams != expectedMethodWithoutParams {
				t.Errorf("Expected method: %s, got: %s", expectedMethodWithoutParams, actualMethodWithoutParams)
			}

			if !equalParams(grpcInfo.Params, tc.expectedParams) {
				t.Errorf("Expected params: %v, got: %v", tc.expectedParams, grpcInfo.Params)
			}
		})
	}
}

// equalParams checks if two maps of parameters are equal.
func equalParams(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for key, val := range a {
		if bVal, ok := b[key]; !ok || bVal != val {
			return false
		}
	}
	
	return true
}

