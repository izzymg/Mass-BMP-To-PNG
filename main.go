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
)

// makeProcessFunc creates a function that attempts to convert any BMP files it finds
// into PNGs and write them into outputDirectory
func makeProcessFunc(inputDirectory string, outputDirectory string, silent bool) func(fileInfo os.FileInfo) error {
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
			fmt.Printf("Processing %s\n", fileInfo.Name())
		}

		// Open file
		inputFile, err := os.Open(filepath.Join(inputDirectory, fileInfo.Name()))
		if err != nil {
			return err
		}
		defer inputFile.Close()

		// Decode BMP into an image
		bmpImage, err := bmp.Decode(inputFile)
		if err != nil {
			return err
		}

		// Transform path "input/img.bmp" into "output/img.png"
		bmpLastDot := strings.LastIndex(fileInfo.Name(), ".")
		bmpOutputFilepath := filepath.Join(outputDirectory, fmt.Sprint(fileInfo.Name()[:bmpLastDot], ".png"))

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
	inputDirFlag := flag.String("input", ".", "Path to process BMP files in")
	outputDirFlag := flag.String("output", ".", "Path to write JPEG files out")
	silent := flag.Bool("silent", false, "Don't print anything to stdout")
	flag.Parse()

	// Trim
	inputDir := trimPath(*inputDirFlag)
	outputDir := trimPath(*outputDirFlag)

	processFile := makeProcessFunc(inputDir, outputDir, *silent)

	// Read files in inputDir and run processFile on each

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		err := processFile(file)
		if err != nil {
			panic(err)
		}
	}
}
