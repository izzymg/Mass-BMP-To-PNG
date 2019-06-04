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
func makeProcessFunc(outputDir string, silent bool, clean bool) func(fileInfo os.FileInfo) error {
	return func(fileInfo os.FileInfo) error {

		// Skip directories
		if fileInfo.IsDir() == true {
			return nil
		}

		filename := fileInfo.Name()

		// Skip files not containing bmp extension
		if filepath.Ext(filename) != ".bmp" {
			return nil
		}

		if silent == false {
			fmt.Printf("Processing \"%s\"\n", filename)
		}

		// Open file
		openedFile, err := os.Open(filename)
		if err != nil {
			return err
		}

		// Defer removal of file for execution after it's closed
		if clean == true {
			defer os.Remove(filename)
		}
		defer openedFile.Close()

		// Decode BMP into an image.Image
		bmpImage, err := bmp.Decode(openedFile)
		if err != nil {
			return err
		}

		// Transform path "input/img.bmp" into "output/img.png"
		fileLastExt := strings.LastIndex(filename, ".")
		outputFp := filepath.Join(outputDir, fmt.Sprint(filename[:fileLastExt], ".png"))

		// Create output file
		output, err := os.Create(outputFp)
		if err != nil {
			return err
		}
		defer output.Close()

		// Encode BMP into PNG file
		err = png.Encode(output, bmpImage)
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

	// Return string without its first and last characters if it's enclosed in "" or ''
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

	processFile := makeProcessFunc(outputDir, *silent, *clean)

	// Try to read directory
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}

	// Change to input directory
	err = os.Chdir(inputDir)
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
