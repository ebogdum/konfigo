package pipeline

import (
	"konfigo/internal/logger"
	"konfigo/internal/parser"
	"konfigo/internal/reader"
	"runtime"
	"sync"
)

// OptimizedFileProcessor provides optimized parallel file processing
type OptimizedFileProcessor struct {
	numWorkers int
}

// NewOptimizedFileProcessor creates a new optimized file processor
func NewOptimizedFileProcessor() *OptimizedFileProcessor {
	return &OptimizedFileProcessor{
		numWorkers: runtime.NumCPU(),
	}
}

// ProcessFiles processes files in parallel with optimized memory usage.
// It preserves the original source index from each entry for ordered merging.
func (ofp *OptimizedFileProcessor) ProcessFiles(entries []sourceEntry, formatOverride string) []parseResult {
	if len(entries) == 0 {
		return nil
	}

	// Limit workers for memory efficiency
	maxWorkers := ofp.numWorkers
	if len(entries) < maxWorkers {
		maxWorkers = len(entries)
	}

	jobs := make(chan sourceEntry, len(entries))
	results := make(chan parseResult, len(entries))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for entry := range jobs {
				result := ofp.processFile(entry.FilePath, formatOverride)
				result.Index = entry.Index
				results <- result
			}
		}()
	}

	// Send jobs
	for _, entry := range entries {
		jobs <- entry
	}
	close(jobs)

	// Wait for completion
	wg.Wait()
	close(results)

	// Collect results
	parsedResults := make([]parseResult, 0, len(entries))
	for res := range results {
		parsedResults = append(parsedResults, res)
	}

	return parsedResults
}

// processFile processes a single file
func (ofp *OptimizedFileProcessor) processFile(path string, formatOverride string) parseResult {
	content, err := reader.ReadFile(path)
	if err != nil {
		logger.Debug("Failed to read file %s: %v", path, err)
		return parseResult{FilePath: path, Err: err}
	}

	data, err := parser.Parse(path, content, formatOverride)
	return parseResult{FilePath: path, Data: data, Err: err}
}

