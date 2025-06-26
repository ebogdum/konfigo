package pipeline

import (
	"konfigo/internal/schema"
)

// processSingle handles single processing mode with schema
func (p *Pipeline) processSingle(baseConfig map[string]interface{}, loadedSchema *schema.Schema, varsFromFileGlobal map[string]interface{}, envVarsForSchema map[string]string) (map[string]interface{}, error) {
	varsToProcess := varsFromFileGlobal
	if varsToProcess == nil {
		varsToProcess = make(map[string]interface{}) // Ensure not nil for schema.Process
	}
	
	processedConfig, err := schema.Process(baseConfig, loadedSchema, varsToProcess, envVarsForSchema)
	if err != nil {
		return nil, err
	}
	
	return processedConfig, nil
}
