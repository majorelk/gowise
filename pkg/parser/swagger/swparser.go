// Package swparser provides functionality to parse and validate Swagger files.
// It supports Swagger versions 1.0, 2.0, and 3.0.
package swparser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// SwaggerInfoInterface is an interface for Swagger information.
// It provides a method to get the Swagger version.
// This interface is used to abstract the details of different Swagger versions.
type SwaggerInfoInterface interface {
	GetVersion() SwaggerVersion
}

type SwaggerVersion string

const (
	Swagger3_0 SwaggerVersion = "3.0"
	Swagger2_0 SwaggerVersion = "2.0"
	Swagger1_0 SwaggerVersion = "1.0"
)

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
func (si SwaggerInfo) GetVersion() SwaggerVersion {
	return SwaggerVersion(si.Version)
}

// GetVersion returns the Swagger version of the SwaggerInfo1.
func (si1 SwaggerInfo1) GetVersion() SwaggerVersion {
	return SwaggerVersion(si1.SwaggerVersion)
}

// GetVersion returns the Swagger version of the SwaggerInfo3.
func (si3 SwaggerInfo3) GetVersion() SwaggerVersion {
	return SwaggerVersion(si3.OpenAPI)
}

func unmarshalAndValidate(content []byte, target interface{}, validate func(interface{}) error) error {
	if err := json.Unmarshal(content, target); err != nil {
		return fmt.Errorf("unable to unmarshal content: %w", err)
	}
	if err := validate(target); err != nil {
		return err
	}
	return nil
}

func ParseSwaggerFile(filePath string) (SwaggerInfoInterface, error) {
	// Reject file paths that are not simple file names
	if filePath != filepath.Base(filePath) {
		return nil, fmt.Errorf("invalid Swagger file path: %s", filePath)
	}

	// Reject file paths that attempt to navigate directories
	if strings.Contains(filePath, "..") || strings.Contains(filePath, "/") {
		return nil, fmt.Errorf("invalid Swagger file path: %s", filePath)
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %s: %w", filePath, err)
	}

	var genericMap map[string]interface{}
	if err := json.Unmarshal(content, &genericMap); err != nil {
		return nil, fmt.Errorf("unable to parse file %s: %w", filePath, err)
	}

	version, ok := genericMap["swagger"].(string)
	if !ok {
		version, ok = genericMap["openapi"].(string)
		if !ok {
			version, ok = genericMap["swaggerVersion"].(string)
			if !ok {
				return nil, errors.New("swagger version is missing")
			}
		}
	}

	switch version {
	case "2.0":
		var swaggerInfo SwaggerInfo
		err := unmarshalAndValidate(content, &swaggerInfo, func(target interface{}) error {
			swaggerInfo, ok := target.(*SwaggerInfo)
			if !ok {
				return errors.New("invalid target type")
			}
			if swaggerInfo.Info.Title == "" || swaggerInfo.Info.Version == "" || len(swaggerInfo.Paths) == 0 {
				return errors.New("invalid Swagger 2.0 file format: Missing required fields")
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return &swaggerInfo, nil

	case "1.0":
		var swaggerInfo1 SwaggerInfo1
		err := unmarshalAndValidate(content, &swaggerInfo1, func(target interface{}) error {
			swaggerInfo1, ok := target.(*SwaggerInfo1)
			if !ok {
				return errors.New("invalid target type")
			}
			if swaggerInfo1.ApiVersion == "" || swaggerInfo1.BasePath == "" || len(swaggerInfo1.Apis) == 0 {
				return errors.New("invalid Swagger 1.0 file format: Missing required fields")
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return &swaggerInfo1, nil

	case "3.0.0":
		var swaggerInfo3 SwaggerInfo3
		err := unmarshalAndValidate(content, &swaggerInfo3, func(target interface{}) error {
			swaggerInfo3, ok := target.(*SwaggerInfo3)
			if !ok {
				return errors.New("invalid target type")
			}
			if swaggerInfo3.Info.Title == "" || swaggerInfo3.Info.Version == "" || len(swaggerInfo3.Paths) == 0 {
				return errors.New("invalid Swagger 3.0 file format: Missing required fields")
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return &swaggerInfo3, nil

	default:
		return nil, fmt.Errorf("unsupported Swagger version: %s", version)
	}
}
