package swparser

import (
	"testing"
	"io/ioutil"
	"os"
)

func TestParseSwaggerFile(t *testing.T) {
	t.Run("Valid Swagger File", func (t *testing.T) {
		// Create a temporary valid Swagger/ OpenAPI file for testing
		tmpSwaggerFile := createTempFile(t, "valid_swagger.json", `{"swagger": "2.0", "info": {"version": "1.0", "title": "Test API"}}`)
		
		// Test parsing the Swagger file
		_, err := ParseSwaggerFile(tmpSwaggerFile)
		if err != nil {
			t.Errorf("Failed to parse valid Swagger file: %v", err)
		}

		// Clean up temporary file
		removeTempFile(t, tmpSwaggerFile)
	})

	t.Run("Invalid Swagger File", func(t *testing.T) {
		// Create a temporary Swagger file with invalid content
		tmpSwaggerFile := createTempFile(t, "invalid_swagger.json", `{"swagger": "2.0", "info": {"version": "1.0"}}`)

		// Test parsing the invalid Swagger file
		_, err := ParseSwaggerFile(tmpSwaggerFile)
		if err == nil {
			t.Error("Expected an error for invalid Swagger file, but got nil")
		} else {
			expectedErrorMessage := "Invalid Swagger file format: Title is required"
			if err.Error() != expectedErrorMessage {
				t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMessage, err.Error())
			}
		}

		// Clean up temporary file
		removeTempFile(t, tmpSwaggerFile)
	})
}

// Helper funcetion to create a temporary file with supplied content for testing
func createTempFile(t *testing.T, fileName, content string) string {
	tmpFile, err := ioutil.TempFile("", fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()

	err = ioutil.WriteFile(tmpFile.Name(), []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	
	return tmpFile.Name()
}

// Helper function to remove a temporary file created for testing
func removeTempFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		t.Errorf("Failed to remove temporary file: %v", err)
	}
}
