package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

func main() {
	inputFile := flag.String("input", "", "Input image file path")
	outputDir := flag.String("output", "output", "Output directory")
	numRows := flag.Int("rows", 1, "Number of rows")
	numCols := flag.Int("cols", 1, "Number of columns")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Input file is required")
		flag.Usage()
		return
	}

	startTime := time.Now()

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		return
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image file:", err)
		return
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	fmt.Printf("Image dimensions: width=%d, height=%d\n", width, height)

	cellWidth := width / *numCols
	cellHeight := height / *numRows

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	for row := 0; row < *numRows; row++ {
		for col := 0; col < *numCols; col++ {
			x0 := col * cellWidth
			y0 := row * cellHeight
			x1 := x0 + cellWidth
			y1 := y0 + cellHeight
			if x1 > width {
				x1 = width
			}
			if y1 > height {
				y1 = height
			}

			subImg := img.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(image.Rect(x0, y0, x1, y1))

			outputFile := filepath.Join(*outputDir, fmt.Sprintf("R%dC%d.%s", row+1, col+1, format))
			outFile, err := os.Create(outputFile)
			if err != nil {
				fmt.Println("Error creating output file:", err)
				return
			}
			defer outFile.Close()

			switch format {
			case "jpeg":
				err = jpeg.Encode(outFile, subImg, nil)
			case "png":
				err = png.Encode(outFile, subImg)
			default:
				fmt.Println("Unsupported image format:", format)
				return
			}
			if err != nil {
				fmt.Println("Error saving image file:", err)
				return
			}
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Image split into grid cells successfully. Time taken: %s\n", elapsedTime)
}
