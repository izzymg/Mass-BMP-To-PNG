package main

import (
	"flag"
	"fmt"
	"golang.org/x/image/bmp"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// makeProcessFunc creates a function that attempts to convert any BMP files it finds
// into PNGs and write them into outputDir
func makeProcessFunc(inputDir string, outputDir string, silent bool, clean bool) func(fileInfo os.FileInfo) error {
	return func(fileInfo os.FileInfo) error {

		// Skip directories
		if fileInfo.IsDir() == true {
			return nil
		}

		// Skip files not containing bmp extension
		if filepath.Ext(fileInfo.Name()) != ".bmp" {
			return nil
		}

		if silent == false {
			fmt.Printf("Processing \"%s\"\n", fileInfo.Name())
		}

		inputFilepath := filepath.Join(inputDir, fileInfo.Name())

		// Open file
		inputFile, err := os.Open(inputFilepath)
		if err != nil {
			return err
		}
		// Defer removal of file for execution after it's closed
		if clean == true {
			defer os.Remove(inputFilepath)
		}
		defer inputFile.Close()

		// Decode BMP into an image
		bmpImage, err := bmp.Decode(inputFile)
		if err != nil {
			return err
		}

		// Transform path "input/img.bmp" into "output/img.png"
		bmpLastDot := strings.LastIndex(fileInfo.Name(), ".")
		bmpOutputFilepath := filepath.Join(outputDir, fmt.Sprint(fileInfo.Name()[:bmpLastDot], ".png"))

		// Create output file
		outputFile, err := os.Create(bmpOutputFilepath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		// Encode BMP into PNG file
		err = png.Encode(outputFile, bmpImage)
		if err != nil {
			return err
		}
		return nil
	}
}

// trimPath cleans p and removes excess whitespace and quote characters
func trimPath(p string) string {
	s := filepath.Clean(strings.TrimSpace(p))
	if len(s) < 2 {
		return s
	}
	firstChar := s[0]
	lastChar := s[len(s)-1]
	if (firstChar == '"' && lastChar == '"') || (firstChar == '\'' && lastChar == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}

func main() {

	// Parse flags
	silent := flag.Bool("silent", false, "Don't print anything to stdout")
	clean := flag.Bool("clean", false, "Delete BMPs after processing")
	concurrency := flag.Int("c", 5, "Number of concurrent operations")
	inputDirFlag := flag.String("input", ".", "Path to process BMP files in")
	outputDirFlag := flag.String("output", ".", "Path to write JPEG files out")
	flag.Parse()

	if *concurrency < 1 {
		*concurrency = 1
	}

	// Timing
	start := time.Now()

	// Trim
	inputDir := trimPath(*inputDirFlag)
	outputDir := trimPath(*outputDirFlag)

	processFile := makeProcessFunc(inputDir, outputDir, *silent, *clean)

	// Try to read directory

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}

	// Setup waitgroup, and semaphore to limit concurrency
	var wg sync.WaitGroup
	wg.Add(len(files))
	var sem = make(chan int, *concurrency)

	// processFile for each file in inputDir
	for _, file := range files {
		sem <- 1
		go func(f os.FileInfo) {
			defer wg.Done()
			err := processFile(f)
			if err != nil {
				panic(err)
			}
			<-sem
		}(file)
	}

	wg.Wait()

	if *silent == false {
		fmt.Printf("Processed %d files in %.3fs\n", len(files), time.Since(start).Seconds())
	}
}
