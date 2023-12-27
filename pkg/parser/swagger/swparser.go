// Package SWParser provides functionality for working with Swagger/OpenAPI files
package swparser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// SwaggerInfo represents the relevant information from a Swagger file.
type SwaggerInfo struct {
	Version	string `json:"swagger"`
	Title 	string `json:"title"`i
}

// ParseSwaggerFile parses a Swagger/ OpenAPI file
//
// It reads the content of the specified file and unmarshals it into a SwaggerInfo struct.
// The function then validates the Swagger version and the required fields.
//
// If the swagger version is "1.0," additional methods for handling version 1.0 can be implemented. (WIP)
// For versions "2.0" and "3.0", the function checks if the Title field is not empty.
// Additional validation methods for these versions can be added as needed.
//
// The function returns a pointer to a SwaggerInfo struct and an error if there are any issues with parsing or validation.
func ParseSwaggerFile(filePath string) (*SwaggerInfo, error) {
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
		// Need methods to handle version 1.0 if needed
	} else {
		switch swaggerInfo.Version {
		case "2.0", "3.0":
			if swaggerInfo.Title =="" {
				return nil, errors.New("Invalid Swagger file format: Title is required")
			}
		// Need more validation methods for versions 2.0 and 3.0
		default: 
			return nil, errors.New("Unsupported Swagger version")
		}
	}

	return &swaggerInfo, nil
}

