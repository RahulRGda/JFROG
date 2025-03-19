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

# Add data file(csg file) to data folder.
<filename.csv>

## Running the Application

To run the application, provide the path to the CSV file as a command-line argument:

```bash
./readCSV data/<filename.csv>
```

# Finding Url contents stored in output folder.

A file with unquie filename(generated using url) is created and url content is stored in txt format.

# Finding Status of each url and its corresponding file name.

A file(csv) is generated in result folder, Which contains the status of each url and its corresponding filename/error is stored in result file.


