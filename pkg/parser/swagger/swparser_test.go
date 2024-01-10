package swparser

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseSwaggerFile(t *testing.T) {
	t.Run("Valid Swagger File", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "valid_swagger.json", `{"swagger": "2.0", "info": {"version": "1.0", "title": "Test API"}, "paths": {"/test": {"get": {"summary": "Test API"}}}}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertNoError(t, err)
	})

	t.Run("Invalid Swagger File", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "invalid_swagger.json", `{"swagger": "2.0", "info": {"version": "1.0", "title": ""}, "paths": {"/test": {"get": {"summary": "Test API"}}}}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)

		expectedErrorMessage := "Invalid Swagger file format: Title is required"
		assertEqual(t, err.Error(), expectedErrorMessage)
	})

	t.Run("Swagger 1.0 File", func(t *testing.T) {
		tmpSwaggerFile := createTempFile(t, "swagger1.json", `{"swaggerVersion": "1.0", "apiVersion": "1.0", "basePath": "/", "apis": [{"path": "/test", "description": "Test API"}]}`)
		defer removeTempFile(t, tmpSwaggerFile)

		_, err := ParseSwaggerFile(tmpSwaggerFile)
		assertError(t, err)

		expectedErrorMessage := "Unsupported Swagger version"
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
}

// Helper functions...

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
