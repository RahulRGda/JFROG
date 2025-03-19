package urlprocessor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/RahulRGda/jfrog/pkg/logging"
	"github.com/RahulRGda/jfrog/pkg/readenv"
)

type TaskResult struct {
	URL      string
	Status   string
	Filename string
	Error    string
}

func Worker(id int, tasks <-chan string, results chan<- TaskResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		logging.Info("Worker", id, "processing task", task)
		filename, err := downloadAndSave(task)
		result := TaskResult{URL: task}
		if err != nil {
			result.Status = "failed"
			result.Filename = ""
			result.Error = err.Error()
			logging.Error("Worker", id, "failed to process task", task, ":", err)
		} else {
			result.Status = "success"
			result.Filename = filename
			result.Error = ""
			logging.Info("Worker", id, "successfully processed task", task, "saved to", filename)
		}
		results <- result
	}
	// results channel closed when task channel is closed, and worker is done.
}

func downloadAndSave(url string) (filename string, err error) {
	maxRetries := readenv.GetEnvInt("maxRetries", 3)
	retryDelay := time.Duration(readenv.GetEnvInt64("retryDelay", 3)) * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		filename, err = downloadAndSaveAttempt(url)
		if err == nil {
			return filename, nil // Success, no need to retry
		}

		if attempt < maxRetries {
			if os.IsTimeout(err) || isRetryableHTTPError(err) {
				logging.Info("Retry attempt", attempt+1, "for", url, ":", err)
				time.Sleep(retryDelay)
				retryDelay *= 2 // Exponential backoff
				continue
			}
			logging.Error("Error downloading", url, "(non-retryable):", err)
			return "", err // Non-retryable error, stop trying
		}

		logging.Error("Max retries exceeded for %s: %v\n", url, err)
		return "", err // Max retries reached
	}
	return filename, nil
}

func isRetryableHTTPError(err error) bool {
	if httpErr, ok := err.(*httpError); ok {
		return httpErr.StatusCode >= 500
	}
	return false
}

type httpError struct {
	StatusCode int
	Err        error
}

func (e *httpError) Error() string {
	return fmt.Sprintf("HTTP error: %d, %v", e.StatusCode, e.Err)
}

func downloadAndSaveAttempt(url string) (outputPath string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", &httpError{StatusCode: resp.StatusCode, Err: fmt.Errorf("HTTP status: %s", resp.Status)}
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	outputPath, err = generateRandomFilename(url, ".txt")
	if err != nil {
		logging.Error("Unable to generate output file name for url:", url)
	}

	if err := os.WriteFile(outputPath, content, 0644); err != nil {
		logging.Error("Error writing file:", err)
		return "", err
	}

	logging.Info(fmt.Sprintf("Downloaded from %s and saved to %s\n", url, outputPath))
	return filepath.Base(outputPath), nil
}

func generateUniqueString(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func getCurrentFileDir() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("could not get caller info")
	}
	return filepath.Dir(filename), nil
}

func generateRandomFilename(url, extension string) (string, error) {
	cwd, err := getCurrentFileDir()
	if err != nil {
		logging.Info("Error getting current working directory:", err)
		return "", err
	}

	filename := generateUniqueString(url)
	// Construct the desired file path
	outputPath := filepath.Join(filepath.Dir(filepath.Dir(cwd)), "output", filename+extension)
	return outputPath, nil
}
