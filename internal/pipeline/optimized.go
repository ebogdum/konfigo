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
	pool       sync.Pool
}

// NewOptimizedFileProcessor creates a new optimized file processor
func NewOptimizedFileProcessor() *OptimizedFileProcessor {
	return &OptimizedFileProcessor{
		numWorkers: runtime.NumCPU(),
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 64*1024) // 64KB initial capacity
			},
		},
	}
}

// ProcessFiles processes files in parallel with optimized memory usage
func (ofp *OptimizedFileProcessor) ProcessFiles(files []string, formatOverride string) []parseResult {
	if len(files) == 0 {
		return nil
	}

	// Limit workers for memory efficiency
	maxWorkers := ofp.numWorkers
	if len(files) < maxWorkers {
		maxWorkers = len(files)
	}

	jobs := make(chan string, len(files))
	results := make(chan parseResult, len(files))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				result := ofp.processFile(path, formatOverride)
				results <- result
			}
		}()
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for completion
	wg.Wait()
	close(results)

	// Collect results
	var parsedResults []parseResult
	for res := range results {
		parsedResults = append(parsedResults, res)
	}

	return parsedResults
}

// processFile processes a single file with memory optimization
func (ofp *OptimizedFileProcessor) processFile(path string, formatOverride string) parseResult {
	// Get buffer from pool
	bufferInterface := ofp.pool.Get()
	buffer := bufferInterface.([]byte)
	defer func() {
		// Reset buffer and return to pool
		buffer = buffer[:0]
		ofp.pool.Put(buffer)
	}()

	content, err := reader.ReadFile(path)
	if err != nil {
		logger.Debug("Failed to read file %s: %v", path, err)
		return parseResult{FilePath: path, Err: err}
	}

	data, err := parser.Parse(path, content, formatOverride)
	return parseResult{FilePath: path, Data: data, Err: err}
}

// BatchProcessor provides optimized batch processing with memory management
type BatchProcessor struct {
	maxBatchSize    int
	memoryThreshold int64 // in bytes
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor() *BatchProcessor {
	return &BatchProcessor{
		maxBatchSize:    runtime.NumCPU() * 2,
		memoryThreshold: 100 * 1024 * 1024, // 100MB
	}
}

// ProcessInBatches processes items in memory-efficient batches
func (bp *BatchProcessor) ProcessInBatches(items []interface{}, processor func(item interface{}) error) error {
	batchSize := bp.maxBatchSize
	
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		
		batch := items[i:end]
		
		// Process batch
		for _, item := range batch {
			if err := processor(item); err != nil {
				return err
			}
		}
		
		// Force garbage collection between batches for memory management
		if i%bp.maxBatchSize == 0 {
			runtime.GC()
		}
	}
	
	return nil
}
