package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/coreycole/datastarui/serviceworker"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// FaviconConfig defines the configuration for generating a favicon
type FaviconConfig struct {
	Width   int
	Height  int
	Name    string
	Purpose string
}

// getFaviconConfigs returns all favicon configurations from serviceworker
func getFaviconConfigs() []FaviconConfig {
	icons := serviceworker.DefaultIcons()
	configs := []FaviconConfig{}

	for _, icon := range icons {
		if icon.Type == "image/png" {
			sizes := strings.Split(icon.Sizes, "x")
			if len(sizes) == 2 {
				width, _ := strconv.Atoi(sizes[0])
				height, _ := strconv.Atoi(sizes[1])

				// Extract filename from src
				filename := filepath.Base(icon.Src)

				configs = append(configs, FaviconConfig{
					Width:   width,
					Height:  height,
					Name:    filename,
					Purpose: icon.Purpose,
				})
			}
		}
	}

	return configs
}

func main() {
	// Parse command line flags
	svgFile := flag.String("svg", "static/favicons/logo.svg", "Path to the SVG file")
	outputDir := flag.String("out", "static/favicons", "Output directory for favicons")
	flag.Parse()

	favicons := getFaviconConfigs()

	// Validate input file exists
	if _, err := os.Stat(*svgFile); os.IsNotExist(err) {
		fmt.Printf("Error: SVG file does not exist: %s\n", *svgFile)
		os.Exit(1)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Read SVG file
	svgData, err := os.ReadFile(*svgFile)
	if err != nil {
		fmt.Printf("Error reading SVG file: %v\n", err)
		os.Exit(1)
	}

	for _, favicon := range favicons {
		if err := generateFavicon(svgData, favicon, *outputDir); err != nil {
			fmt.Printf("Error generating %s: %v\n", favicon.Name, err)
			continue
		}
		fmt.Printf("Generated: %s (%dx%d)\n", favicon.Name, favicon.Width, favicon.Height)
	}

	fmt.Println("All favicons generated successfully!")
}

func generateFavicon(svgData []byte, config FaviconConfig, outputDir string) error {
	// Parse SVG
	icon, err := oksvg.ReadIconStream(strings.NewReader(string(svgData)), oksvg.IgnoreErrorMode)
	if err != nil {
		return fmt.Errorf("failed to parse SVG: %w", err)
	}

	// Set viewbox to icon dimensions
	icon.SetTarget(0, 0, float64(config.Width), float64(config.Height))

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))

	// Create scanner and rasterizer
	scanner := rasterx.NewScannerGV(config.Width, config.Height, img, img.Bounds())
	raster := rasterx.NewDasher(config.Width, config.Height, scanner)

	// Draw SVG to image
	icon.Draw(raster, 1.0)

	// Create output file
	outputPath := filepath.Join(outputDir, config.Name)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Encode as PNG
	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}
