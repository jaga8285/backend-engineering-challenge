# Event parser

Command line tool for parsing events stored as json objects in an input file

The parser is split between 4 main stages:
* File reading (internal/ioProcess/inputProcess.go)
* Average per minute calculation (internal/calc/averageMinute)
* Moving average calculation (internal/calc/averageMinute)
* File writing (internal/ioProcess/outputProcess.go)

All stages run concurrently, meaning a stage can start before the previous one finishes

## Build
``` go build -o event_cli cmd/main.go```

## Usage
``` ./event_cli --window_size 10 --input_file test/test1.in```
Flags:
* ```--window_size``` specifies the size of the moving average window (mandatory)
* ```--input_file | -i``` specifies the name of the input file (mandatory)
* ```--output_file | -o``` specifies the name of the output file (omit this flag to print output to stdout)
* ```--num_threads``` specifies the number of threads used for average calculation (default:4)

## Testing

Integration tests:
``` go test event_cli/cmd```

Unit tests:
``` go test ./...```
