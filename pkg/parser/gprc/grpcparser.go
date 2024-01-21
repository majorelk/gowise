// Package grpcparser provides utilities for parsing gRPC messages.
//
// The main function in this package is ParseGRPCMessage, which takes a gRPC message string
// and returns a GRPCInfo struct containing the service name, method name, and any parameters
// present in the message.
//
// The GRPCInfo struct is defined as follows:
//
// type GRPCInfo struct {
//     ServiceName string     // The name of the gRPC service
//     MethodName  string     // The name of the method being called
//     Params      url.Values // Any parameters included in the message
// }
//
// The ParseGRPCMessage function splits the input message on the "/" character to separate the
// service name from the rest of the message. It then splits the remainder of the message on the
// "?" character to separate the method name from the parameters. If parameters are present,
// they are parsed using the url.ParseQuery function from the net/url package.
//
// If the input message is not in the expected format, ParseGRPCMessage will return an error.
// The expected format is "Service/Method", with an optional "?param1=value1&param2=value2"
// at the end for parameters.
package grpcparser

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	ErrInvalidMessageFormat = "invalid gRPC message format: expected 'Service/Method'"
	ErrInvalidParamFormat   = "invalid parameter format: expected 'key=value'"
)

// GRPCInfo represents a parsed gRPC message. It contains the name of the gRPC service,
// the name of the method being called, and any parameters included in the message.
// This struct is used to provide a structured representation of a gRPC message, which
// can be useful for inspecting the message or passing it to other functions.

type GRPCInfo struct {
	ServiceName string
	MethodName  string
	Params      url.Values
}

// ParseGRPCMessage takes a gRPC message string and returns a GRPCInfo struct containing
// the parsed information. The input message is expected to be in the format "Service/Method",
// with an optional "?param1=value1&param2=value2" at the end for parameters. The function
// splits the input message on the "/" character to separate the service name from the rest
// of the message, then splits the remainder of the message on the "?" character to separate
// the method name from the parameters. If parameters are present, they are parsed using the
// url.ParseQuery function from the net/url package. If the input message is not in the
// expected format, ParseGRPCMessage will return an error.
func ParseGRPCMessage(message string) (*GRPCInfo, error) {
	parts := strings.Split(message, "/")

	if len(parts) < 2 {
		return nil, errors.New(ErrInvalidMessageFormat)
	}

	grpcInfo := &GRPCInfo{
		ServiceName: parts[0],
	}

	// Split the method part into the method name and parameters
	methodAndParams := strings.Split(parts[1], "?")
	grpcInfo.MethodName = methodAndParams[0]

	// Extract parameters if present
	if len(methodAndParams) > 1 {
		params, err := url.ParseQuery(methodAndParams[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse parameters: %w", err)
		}
		grpcInfo.Params = params
	}

	return grpcInfo, nil
}
