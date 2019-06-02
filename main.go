package main

import (
	"flag"
	"fmt"
	"golang.org/x/image/bmp"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// makeWalkFunc creates a WalkFunc that attempts to convert any BMP files it finds
// into PNGs and write them into outputDirectory
func makeWalkFunc(outputDirectory string) filepath.WalkFunc {
	return func(inputFilepath string, inputFileInfo os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// Skip directories
		if inputFileInfo.IsDir() == true {
			return nil
		}

		// Skip files not containing bmp extension
		if filepath.Ext(inputFileInfo.Name()) != ".bmp" {
			return nil
		}

		// Open file
		fmt.Printf("Opening %s\n", inputFilepath)
		inputFile, err := os.Open(inputFilepath)
		if err != nil {
			return err
		}
		defer inputFile.Close()

		// Decode BMP into an image
		fmt.Println("Decoding")
		bmpImage, err := bmp.Decode(inputFile)
		if err != nil {
			return err
		}

		// Transform path "input/img.bmp" into "output/img.png"
		bmpBasename := filepath.Base(inputFilepath)
		bmpLastDot := strings.LastIndex(bmpBasename, ".")
		bmpOutputFilepath := filepath.Join(outputDirectory, fmt.Sprint(bmpBasename[:bmpLastDot], ".png"))

		// Create output file
		fmt.Printf("Creating %s\n", bmpOutputFilepath)
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

// trimInputPath cleans p and removes excess whitespace and quote characters
func trimInputPath(p string) string {
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
	inputDirectory := flag.String("input", ".", "Path to process BMP files in")
	outputDirectory := flag.String("output", ".", "Path to write JPEG files out")
	flag.Parse()

	walkFunc := makeWalkFunc(trimInputPath(*inputDirectory))

	err := filepath.Walk(trimInputPath(*outputDirectory), walkFunc)
	if err != nil {
		panic(err)
	}
}
