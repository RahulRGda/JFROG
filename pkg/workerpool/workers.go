package workerpool

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/RahulRGda/jfrog/pkg/logging"
	"github.com/RahulRGda/jfrog/pkg/urlprocessor"
)

// Worker function that processes tasks

func StartWorkers(URLS []string) {
	numWorkers := len(URLS)
	if numWorkers > 50 {
		numWorkers = 50
	}

	numTasks := len(URLS)

	tasks := make(chan string, numTasks)
	results := make(chan urlprocessor.TaskResult, numTasks) // Use numTasks for buffer

	var wg sync.WaitGroup

	// Start workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go urlprocessor.Worker(i, tasks, results, &wg)
	}

	// Send tasks to the channel
	for _, url := range URLS {
		tasks <- url
	}
	close(tasks) // Close task channel to signal no more tasks

	// Wait for all workers to finish
	wg.Wait()

	close(results) // Close results channel after all workers finish

	var allResults []urlprocessor.TaskResult
	for result := range results {
		allResults = append(allResults, result)
	}

	if err := writeResultsToCSV(allResults); err != nil {
		log.Fatalf("Error writing CSV: %v", err)
	}
}
func writeResultsToCSV(results []urlprocessor.TaskResult) error {
	currentDir, err := getCurrentFileDir()
	if err != nil {
		logging.Error("Error getting current dir", err)
	}
	file, err := os.Create(filepath.Dir(filepath.Dir(currentDir)) + "/result/output" + time.Now().Format("02-01-2006") + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"URL", "Status", "Filename", "Error"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, result := range results {
		row := []string{result.URL, result.Status, result.Filename, result.Error}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func getCurrentFileDir() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("could not get caller info")
	}
	return filepath.Dir(filename), nil
}
