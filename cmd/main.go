package main

import (
	"encoding/csv"
	"os"
	"path/filepath"

	"github.com/jfrog/assesment/pkg/logging"
	"github.com/jfrog/assesment/pkg/workerpool"
)

func main() {
	logging.InitLogger()
	args := os.Args

	if len(args) != 2 {
		logging.Info("Usage: readCSV <csvfilename>")
		return
	}

	csvFilename := args[1]

	// Check if the file exists and is a regular file
	fileInfo, err := os.Stat(csvFilename)
	if err != nil {
		if os.IsNotExist(err) {
			logging.Error("Unable to find file", csvFilename, "due to the error:", err)
		} else {
			logging.Error("Error: %v\n", err)
		}
		return
	}

	if !fileInfo.Mode().IsRegular() {
		logging.Error("Error: '%s' is not a regular file.\n", csvFilename)
		return
	}

	// Check for .csv extension
	if filepath.Ext(csvFilename) != ".csv" {
		logging.Error("Error: '%s' is not a .csv file.\n", csvFilename)
		return
	}

	file, err := os.Open(csvFilename)
	if err != nil {
		logging.Error("Error opening file: %v\n", err)
		return
	}
	logging.Info("Reading file %s ....\n", csvFilename)

	defer file.Close()

	reader := csv.NewReader(file)

	urls, err := reader.ReadAll()
	if err != nil {
		logging.Error("Error reading CSV: %v\n", err)
		return
	}

	var urlss []string
	for _, url := range urls[1:] {
		urlss = append(urlss, url...)
	}

	workerpool.StartWorkers(urlss)

}
