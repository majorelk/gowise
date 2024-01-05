package grpcparser

import (
	"strings"
	"fmt"
)

// GRPCInfo represents information parsed from a gRPC message
type GRPCInfo struct {
	ServiceName 	string
	MethodName	string
	Params		map[string]string
}

// ParseGRPCMessage parses a gRPC message and returns relevant information.
func ParseGRPCMessage(message string) *GRPCInfo {
	parts := strings.Split(message, "/")

	grpcInfo := &GRPCInfo{}

	if len(parts) >= 2 {
		grpcInfo.ServiceName = parts[0]
		grpcInfo.MethodName = parts[1]
	}

	// Extract parameters from the method
	if len(parts) >= 3 {
		paramsAndMethod := strings.Split(parts[2], "?")

		// Print statements for debugging
		fmt.Println("Message:", message)
		fmt.Println("Parts:", parts)
		fmt.Println("ParamsAndMethod:", paramsAndMethod)

		// Extract Parameters if present
		if len(paramsAndMethod) > 1 {
			grpcInfo.Params = make(map[string]string)

			for _, param := range strings.Split(paramsAndMethod[1], "&") {
				keyVal := strings.Split(param, "=")

				if len(keyVal) == 2 {
					grpcInfo.Params[keyVal[0]] = keyVal[1]
				}
			}
		}
		

		// Print statements for debugging
		fmt.Println("Message:", message)
		fmt.Println("Parts:", parts)
		fmt.Println("ParamsAndMethod:", paramsAndMethod)

	}
	
	return grpcInfo
}
