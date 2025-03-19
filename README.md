# URL Downloader from CSV

This command-line application downloads content from URLs listed in a CSV file, saves the downloaded content to randomly named text files within the output folder, and stores the URL processing status and output file name in date-wise result files within the results folder.

## Installation

1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/RahulRGda/JFROG.git
    cd jfrog
    ```

2.  **Build the Application:**
    ```bash
    go build -o readCSV cmd/main.go
    ```

# Add data file(csv file) to data folder.
<filename.csv>

## Running the Application

To run the application, provide the path to the CSV file as a command-line argument:

```bash
./readCSV data/<filename.csv>
```

# Output.

1. **Downloaded URL Content:**
    The downloaded content from each URL is saved as a .txt file with a unique, randomly generated filename within the output folder

2. **URL Processing Status and Results:**
    A CSV file is generated in the results folder. This file contains the processing status of each URL, including the corresponding output filename (if successful) or the error message (if the URL download failed). The results are stored in files organized by date.

# Errors Logging.

All errors encountered during the application's execution are logged to a file within the onelog folder.
