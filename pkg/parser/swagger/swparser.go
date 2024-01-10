// Package swparser provides functionality to parse and validate Swagger files.
// It supports Swagger versions 1.0, 2.0, and 3.0.
package swparser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// SwaggerInfoInterface is an interface for Swagger information.
// It provides a method to get the Swagger version.
type SwaggerInfoInterface interface {
	GetVersion() string
}

// Operation represents a Swagger operation object.
// An operation describes a single API operation on a path.
// Summary is a short summary of the operation.
// Description is a verbose explanation of the operation.
type Operation struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

// SwaggerInfo represents a Swagger 2.0 or 3.0 document.
// Version is the Swagger version (should be "2.0" or "3.0").
// Info provides metadata about the API. The metadata can be used by the clients if needed.
// Paths holds the relative paths to the individual endpoints. The path is appended to the URL from the Server Object in order to construct the full URL.
type SwaggerInfo struct {
	Version string `json:"swagger"` // The Swagger version
	Info    struct {
		Title   string `json:"title"`
		Version string `json:"version"`
	} `json:"info"`
	Paths map[string]struct {
		Get     *Operation `json:"get"`
		Put     *Operation `json:"put"`
		Post    *Operation `json:"post"`
		Delete  *Operation `json:"delete"`
		Options *Operation `json:"options"`
		Head    *Operation `json:"head"`
		Patch   *Operation `json:"patch"`
		Trace   *Operation `json:"trace"`
	} `json:"paths"`
}

// Info represents the information section of a Swagger document.
type Info struct {
	Title             string `json:"title"`
	Description       string `json:"description"`
	TermsOfServiceUrl string `json:"termsOfServiceUrl"`
	Contact           string `json:"contact"`
	License           string `json:"license"`
	LicenseUrl        string `json:"licenseUrl"`
}

// SwaggerInfo1 represents a Swagger 1.0 document.
type SwaggerInfo1 struct {
	SwaggerVersion string `json:"swaggerVersion"`
	ApiVersion     string `json:"apiVersion"`
	BasePath       string `json:"basePath"`
	Apis           []struct {
		Path        string `json:"path"`
		Description string `json:"description"`
	} `json:"apis"`
	Info `json:"info"`
}

// SwaggerInfo3 represents a Swagger 3.0 document.
// It includes fields for the OpenAPI version, API information, servers, paths, components, security requirements, and tags.
type SwaggerInfo3 struct {
	OpenAPI string `json:"openapi"` // The OpenAPI version
	Info    struct {
		Title       string `json:"title"`       // The title of the API
		Description string `json:"description"` // A short description of the API
		Version     string `json:"version"`     // The version of the API
	} `json:"info"`
	Servers []struct {
		URL string `json:"url"` // The URL of the servers where the API is available
	} `json:"servers"`
	Paths      map[string]interface{} `json:"paths"`      // The available paths and operations for the API
	Components map[string]interface{} `json:"components"` // The available components for the API
	Security   []map[string][]string  `json:"security"`   // The security requirements for the API
	Tags       []struct {
		Name string `json:"name"` // The available tags for the API
	} `json:"tags"`
}

// GetVersion returns the Swagger version of the SwaggerInfo.
// The returned version should be "2.0" or "3.0".
func (si SwaggerInfo) GetVersion() string {
	return si.Version
}

// GetVersion returns the Swagger version of the SwaggerInfo1.
// The returned version should be "1.0".
func (si1 SwaggerInfo1) GetVersion() string {
	return si1.SwaggerVersion
}

// GetVersion returns the Swagger version of the SwaggerInfo3.
// The returned version should be "3.0".
func (si3 SwaggerInfo3) GetVersion() string {
	return si3.OpenAPI
}

// ParseSwaggerFile parses a Swagger file located at filePath and returns a SwaggerInfoInterface representing the parsed information.
// The Swagger file should be in JSON format and conform to the Swagger 1.0, 2.0, or 3.0 specification.
// If the file cannot be read, or if the file content cannot be parsed as a Swagger document, an error is returned.
// An error is also returned if the Swagger document is invalid (e.g., it is missing required fields) or if it uses an unsupported Swagger version.
func ParseSwaggerFile(filePath string) (SwaggerInfoInterface, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var swaggerInfo SwaggerInfo
	if err := json.Unmarshal(content, &swaggerInfo); err != nil {
		return nil, err
	}

	if swaggerInfo.Version == "" {
		var swaggerInfo1 SwaggerInfo1
		if err := json.Unmarshal(content, &swaggerInfo1); err != nil {
			return nil, err
		}
		if swaggerInfo1.SwaggerVersion != "" {
			if swaggerInfo1.ApiVersion == "" || swaggerInfo1.BasePath == "" || len(swaggerInfo1.Apis) == 0 {
				return nil, errors.New("invalid Swagger 1.0 file format: Missing required fields")
			}
			return &swaggerInfo1, nil
		}

		var swaggerInfo3 SwaggerInfo3
		if err := json.Unmarshal(content, &swaggerInfo3); err == nil && swaggerInfo3.OpenAPI != "" {
			if swaggerInfo3.Info.Title == "" || swaggerInfo3.Info.Version == "" || len(swaggerInfo3.Paths) == 0 {
				return nil, errors.New("invalid Swagger 3.0 file format: Missing required fields")
			}
			return &swaggerInfo3, nil
		}

		return nil, errors.New("unsupported Swagger version")
	} else {
		switch swaggerInfo.Version {
		case "2.0":
			if swaggerInfo.Info.Title == "" {
				return nil, errors.New("invalid Swagger file format: Title is required")
			}
			if swaggerInfo.Info.Version == "" || len(swaggerInfo.Paths) == 0 {
				return nil, errors.New("invalid Swagger 2.0 file format: Missing required fields")
			}
			return &swaggerInfo, nil
		default:
			return nil, errors.New("unsupported Swagger version")
		}
	}
}
