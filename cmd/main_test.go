package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_run(t *testing.T) {
	tests := []struct {
		name       string
		cfg        config
		targetFile string
	}{
		{
			name: "single thread test",
			cfg: config{
				InputFile:  "../test/test1.in",
				OutputFile: "../test/test1.myout",
				WindowSize: 10,
				NumWorkers: 1,
			},
			targetFile: "../test/test1.out",
		},
		{
			name: "multi thread test",
			cfg: config{
				InputFile:  "../test/test1.in",
				OutputFile: "../test/test1.myout",
				WindowSize: 10,
				NumWorkers: 4,
			},
			targetFile: "../test/test1.out",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run(tt.cfg)
			assert.True(t, fileCompare(tt.cfg.OutputFile, tt.targetFile))
		})
	}
}

func fileCompare(file1, file2 string) bool {
	sf, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}

	df, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}

	sscan := bufio.NewScanner(sf)
	dscan := bufio.NewScanner(df)

	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return false
		}
	}

	return true
}
