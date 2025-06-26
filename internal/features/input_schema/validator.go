package input_schema

import (
	"fmt"
	"konfigo/internal/logger"
	"reflect"
)

// Validate validates a configuration map against an input schema.
func Validate(config map[string]interface{}, ref *Ref) error {
	logger.Log("Validating against input schema: %s", ref.Path)
	
	schemaMap, err := LoadSchemaMap(ref)
	if err != nil {
		return err
	}
	
	return compareStructure(config, schemaMap, "", ref.Strict)
}

// compareStructure recursively compares the structure of data against a schema.
// It validates that required keys exist and types match.
func compareStructure(data, schema map[string]interface{}, path string, strict bool) error {
	for key, schemaVal := range schema {
		currentPath := buildPath(path, key)
		
		dataVal, exists := data[key]
		if !exists {
			return fmt.Errorf("input schema validation: path '%s' missing required key", currentPath)
		}
		
		if err := validateTypeMatch(dataVal, schemaVal, currentPath, strict); err != nil {
			return err
		}
	}
	
	// In strict mode, check for unexpected keys in data
	if strict {
		for key := range data {
			if _, exists := schema[key]; !exists {
				currentPath := buildPath(path, key)
				return fmt.Errorf("input schema validation: path '%s' unexpected key found in strict mode", currentPath)
			}
		}
	}
	
	return nil
}

// validateTypeMatch validates that data value type matches schema value type.
func validateTypeMatch(dataVal, schemaVal interface{}, path string, strict bool) error {
	schemaType := reflect.TypeOf(schemaVal)
	dataType := reflect.TypeOf(dataVal)
	
	// Handle nested maps
	if schemaMap, isSchemaMap := schemaVal.(map[string]interface{}); isSchemaMap {
		if dataMap, isDataMap := dataVal.(map[string]interface{}); isDataMap {
			return compareStructure(dataMap, schemaMap, path, strict)
		}
		return fmt.Errorf("input schema validation: path '%s' type mismatch, expected map, got %v", path, dataType)
	}
	
	// Handle type compatibility
	if !areTypesCompatible(schemaType, dataType) {
		return fmt.Errorf("input schema validation: path '%s' type mismatch, expected %v, got %v", path, schemaType, dataType)
	}
	
	return nil
}

// areTypesCompatible checks if two types are compatible for schema validation.
func areTypesCompatible(schemaType, dataType reflect.Type) bool {
	if schemaType == dataType {
		return true
	}
	
	// Special case: JSON numbers are often float64, but may represent integers
	// Allow any numeric type to match any other numeric type
	if isNumericType(schemaType) && isNumericType(dataType) {
		return true
	}
	
	return false
}

// isNumericType checks if a type is numeric (int variants or float variants)
func isNumericType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		 reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		 reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// buildPath constructs a dot-separated path for error reporting.
func buildPath(basePath, key string) string {
	if basePath == "" {
		return key
	}
	return basePath + "." + key
}
