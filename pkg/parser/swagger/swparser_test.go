package swparser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestParseSwaggerFile(t *testing.T) {
	t.Run("Valid Swagger File", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "valid_swagger.json", `{"swagger": "2.0", "info": {"version": "1.0", "title": "Test API"}, "paths": {"/test": {"get": {"summary": "Test API"}}}}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertNoError(t, err)
	})

	t.Run("Invalid_Swagger_File", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "invalid_swagger.json", `{"swagger": "2.0", "info": {"version": "1.0", "title": ""}, "paths": {"/test": {"get": {"summary": "Test API"}}}}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)

		expectedErrorMessage := "invalid Swagger 2.0 file format: Missing required fields"
		assertEqual(t, err.Error(), expectedErrorMessage)
	})

	t.Run("Unsupported_Swagger_Version", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "swagger1.json", `{"swaggerVersion": "", "apiVersion": "1.0", "basePath": "/", "apis": [{"path": "/test", "description": "Test API"}]}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)

		expectedErrorMessage := "unsupported Swagger version: "
		assertEqual(t, err.Error(), expectedErrorMessage)
	})

	t.Run("Swagger 3.0 File", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "swagger3.json", `{"openapi": "3.0.0", "info": {"version": "1.0", "title": "Test API"}, "paths": {"/test": {"get": {"summary": "Test API"}}}}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertNoError(t, err)
	})

	t.Run("Extra Fields", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "extra_fields.json", `{"swagger": "2.0", "info": {"version": "1.0", "title": "Test API", "extra": "field"}, "paths": {"/test": {"get": {"summary": "Test API"}}}}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertNoError(t, err)
	})

	t.Run("Unsupported Swagger Version", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "unsupported_swagger.json", `{"swagger": "4.0", "info": {"version": "1.0", "title": "Test API"}}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Non-Existent File", func(t *testing.T) {
		_, err := ParseSwaggerFile("non_existent.json")
		assertError(t, err)
	})

	t.Run("Invalid JSON File", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "invalid_json.json", `{"swagger": "2.0", "info": {"version": "1.0", "title": "Test API"`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "extra_fields.json", `{"swagger": "2.0", "info": {"version": "1.0", "title": "Test API", "extra": "field"}, "paths": {}}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Different Operations", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "operations.json", `{"swagger": "2.0", "info": {"version": "1.0", "title": "Test API"}, "paths": {"/test": {"get": {}, "post": {}, "put": {}, "delete": {}}}}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertNoError(t, err)
	})

	t.Run("Empty File", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "empty.json", ``)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Valid JSON but Invalid Swagger File", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "valid_json_invalid_swagger.json", `{"foo": "bar"}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Valid Swagger 1.0 File", func(t *testing.T) {
		swaggerContent := `{
			"swaggerVersion": "1.0",
			"apiVersion": "1.0.0",
			"basePath": "http://example.com",
			"apis": [{"path": "/test"}]
		}`
		tmpSwaggerFile := createTempFile(t, "valid_swagger_1_0.json", swaggerContent)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertNoError(t, err)
	})

	t.Run("Invalid Swagger Version", func(t *testing.T) {
		swaggerContent := `{
			"swagger": "4.0",
			"info": {
				"title": "Test API",
				"version": "1.0.0"
			},
			"paths": {}
		}`
		tmpSwaggerFile := createTempFile(t, "invalid_swagger_version.json", swaggerContent)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Missing Fields in Swagger 3.0.0 File", func(t *testing.T) {
		swaggerContent := `{
			"openapi": "3.0.0",
			"info": {
				"title": "Test API"
			},
			"paths": {}
		}`
		tmpSwaggerFile := createTempFile(t, "missing_fields_swagger_3_0_0.json", swaggerContent)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Missing Fields in Swagger 1.0 File", func(t *testing.T) {
		swaggerContent := `{
			"swaggerVersion": "1.0",
			"basePath": "http://example.com",
			"apis": [{"path": "/test"}]
		}`
		tmpSwaggerFile := createTempFile(t, "missing_fields_swagger_1_0.json", swaggerContent)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Invalid JSON in Swagger 3.0.0 File", func(t *testing.T) {
		swaggerContent := `{
			"openapi": "3.0.0",
			"info": {
				"title": "Test API",
				"version": "1.0.0"
			},
			"paths": {}
		`
		tmpSwaggerFile := createTempFile(t, "invalid_json_swagger_3_0_0.json", swaggerContent)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Invalid JSON in Swagger 1.0 File", func(t *testing.T) {
		swaggerContent := `{
			"swaggerVersion": "1.0",
			"apiVersion": "1.0.0",
			"basePath": "http://example.com",
			"apis": [{"path": "/test"}]
		`
		tmpSwaggerFile := createTempFile(t, "invalid_json_swagger_1_0.json", swaggerContent)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("File Reading Error", func(t *testing.T) {
		tmpSwaggerFile := "/path/to/non/existent/file.json"

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)
	})

	t.Run("Different Operations", func(t *testing.T) {
		swaggerContent := `{
			"swagger": "2.0",
			"info": {
				"title": "Test API",
				"version": "1.0.0"
			},
			"paths": {
				"/test": {
					"get": {},
					"post": {},
					"put": {},
					"delete": {}
				}
			}
		}`
		tmpSwaggerFile := createTempFile(t, "different_operations.json", swaggerContent)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertNoError(t, err)
	})
}

// Helper functions...
func createTempFile(t *testing.T, fileName, content string) string {
	filePath := filepath.Join(".", fileName)

	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	return filePath
}

func removeTempFile(t *testing.T, filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		t.Errorf("Failed to remove temporary file: %v", err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("Expected an error but got none")
	}
}

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Fatalf("Expected %v but got %v", want, got)
	}
}
