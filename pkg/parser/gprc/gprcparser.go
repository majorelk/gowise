package grpcparser

import (
	"strings"
)

// GRPCInfo represents information parsed from a gRPC message
type GRPCInfo struct {
	ServiceName string
	MethodName  string
	Params      map[string]string
}

// ParseGRPCMessage parses a gRPC message and returns relevant information.
func ParseGRPCMessage(message string) *GRPCInfo {
	parts := strings.Split(message, "/")

	grpcInfo := &GRPCInfo{}

	if len(parts) >= 2 {
		grpcInfo.ServiceName = parts[0]

		// Split the method part into the method name and parameters
		methodAndParams := strings.Split(parts[1], "?")
		grpcInfo.MethodName = methodAndParams[0]

		// Extract parameters if present
		if len(methodAndParams) > 1 {
			grpcInfo.Params = make(map[string]string)

			for _, param := range strings.Split(methodAndParams[1], "&") {
				keyVal := strings.Split(param, "=")

				if len(keyVal) == 2 {
					grpcInfo.Params[keyVal[0]] = keyVal[1]
				}
			}
		}
	}

	return grpcInfo
}
