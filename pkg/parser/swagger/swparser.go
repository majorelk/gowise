// Package SWParser provides functionality for working with Swagger/OpenAPI files
package swparser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// SwaggerInfo represents the relevant information from a Swagger file.
// Version is the Swagger version of the file.
// Info is a struct that contains additional information about the API.
type SwaggerInfo struct {
	Version string `json:"swagger"`
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

type Operation struct {
	Summary     string `json:"summary"`
	Description string `json:"description"`
	// Add other fields as needed
}

// InfoField represents the 'info' field in a Swagger file.
// Version is the version of the API.
// Title is the title of the API.
type InfoField struct {
	Version string `json:"version"`
	Title   string `json:"title"`
}

// Swagger 1 (Spec 1.2) struct
type SwaggerInfo1 struct {
	SwaggerVersion string `json:"swaggerVersion"`
	ApiVersion     string `json:"apiVersion"`
	BasePath       string `json:"basePath"`
	Apis           []struct {
		Path        string `json:"path"`
		Description string `json:"description"`
	} `json:"apis"`
	Info struct {
		Title             string `json:"title"`
		Description       string `json:"description"`
		TermsOfServiceUrl string `json:"termsOfServiceUrl"`
		Contact           string `json:"contact"`
		License           string `json:"license"`
		LicenseUrl        string `json:"licenseUrl"`
	} `json:"info"`
}

// SwaggerInfoInterface is an interface that includes the methods needed for working with SwaggerInfo structs.
type SwaggerInfoInterface interface {
	GetVersion() string
}

func (si SwaggerInfo) GetVersion() string {
	return si.Version
}

func (si1 SwaggerInfo1) GetVersion() string {
	return si1.SwaggerVersion
}

// ParseSwaggerFile parses a Swagger/OpenAPI file.
//
// filePath is the path to the Swagger file.
//
// The function reads the content of the specified file and unmarshals it into a SwaggerInfo struct.
// It then validates the Swagger version and the required fields.
//
// If the Swagger version is "1.0", additional methods for handling version 1.0 can be implemented. (WIP)
// For versions "2.0" and "3.0", the function checks if the Title field is not empty.
// Additional validation methods for these versions can be added as needed.
//
// The function returns a pointer to a SwaggerInfo struct and an error if there are any issues with parsing or validation.
// Possible errors include:
// - An error occurred while reading the file
// - An error occurred while unmarshalling the JSON content
// - The Swagger version is unsupported
// - The Title field is empty for Swagger version "2.0" or "3.0"
func ParseSwaggerFile(filePath string) (SwaggerInfoInterface, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON content into the SwaggerInfo struct
	var swaggerInfo SwaggerInfo
	if err := json.Unmarshal(content, &swaggerInfo); err != nil {
		return nil, err
	}

	// Validate that the required fields are present
	if swaggerInfo.Version == "1.0" {
		var swaggerInfo1 SwaggerInfo1
		if err := json.Unmarshal(content, &swaggerInfo1); err != nil {
			return nil, err
		}
		// Validate the Swagger 1.0 fields
		if swaggerInfo1.SwaggerVersion == "" || swaggerInfo1.ApiVersion == "" || swaggerInfo1.BasePath == "" || len(swaggerInfo1.Apis) == 0 {
			return nil, errors.New("invalid Swagger 1.0 file format: Missing required fields")
		}
		return &swaggerInfo1, nil
	} else {
		switch swaggerInfo.Version {
		case "2.0", "3.0":
			// Validate the Swagger 2.0 and 3.0 fields
			if swaggerInfo.Info.Title == "" {
				return nil, errors.New("Invalid Swagger file format: Title is required")
			}
			if swaggerInfo.Info.Version == "" || len(swaggerInfo.Paths) == 0 {
				return nil, errors.New("invalid Swagger 2.0/3.0 file format: Missing required fields")
			}
			return &swaggerInfo, nil
		default:
			return nil, errors.New("Unsupported Swagger version")
		}
	}
}
