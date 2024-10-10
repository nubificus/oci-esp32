package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func extractLayer(layer v1.Layer, destDir string) error {
	layerReader, err := layer.Compressed()
	if err != nil {
		return fmt.Errorf("failed to get layer: %v", err)
	}
	defer layerReader.Close()

	gzipReader, err := gzip.NewReader(layerReader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar: %v", err)
		}

		if header.Typeflag == tar.TypeReg {
			targetPath := filepath.Join(destDir, header.Name)
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to copy file: %v", err)
			}
			fmt.Printf("Extracted: %s\n", targetPath)
		}
	}
	return nil
}

func printImageMetadata(img v1.Image) {
	digest, err := img.Digest()
	if err != nil {
		log.Fatalf("Error fetching image digest: %v", err)
	}
	fmt.Printf("Image Digest: %s\n", digest)

	configFile, err := img.ConfigFile()
	if err != nil {
		log.Fatalf("Error fetching image config file: %v", err)
	}
	fmt.Printf("Architecture: %s\n", configFile.Architecture)
	fmt.Printf("OS: %s\n", configFile.OS)
	fmt.Printf("Created: %s\n", configFile.Created.Format(time.RFC3339))

	if len(configFile.Config.Env) > 0 {
		fmt.Println("Environment Variables:")
		for _, env := range configFile.Config.Env {
			fmt.Printf("  %s\n", env)
		}
	} else {
		fmt.Println("No environment variables set.")
	}

	layers, err := img.Layers()
	if err != nil {
		log.Fatalf("Error fetching image layers: %v", err)
	}
	for i, layer := range layers {
		size, err := layer.Size()
		if err != nil {
			log.Fatalf("Error fetching layer size: %v", err)
		}
		diffID, err := layer.DiffID()
		if err != nil {
			log.Fatalf("Error fetching layer DiffID: %v", err)
		}
		fmt.Printf("Layer %d - Size: %d bytes, DiffID: %s\n", i+1, size, diffID)
	}
}

func main() {
	imageRef := flag.String("image", "alpine:latest", "Image reference to fetch")
	destDir := flag.String("dest", "./extracted_files", "Destination directory for extracted files")
	imgArch := flag.String("arch", "amd64", "Image Architecture to fetch")
	imgOs := flag.String("os", "linux", "Image OS to fetch")

	flag.Parse()

	ref, err := name.ParseReference(*imageRef)
	if err != nil {
		log.Fatalf("Error parsing image reference: %v", err)
	}

	img, err := remote.Image(ref, remote.WithPlatform(v1.Platform {
		OS:           *imgOs,
		Architecture: *imgArch,
	}))
	if err != nil {
		log.Fatalf("Error fetching image: %v", err)
	}

	printImageMetadata(img)

	layers, err := img.Layers()
	if err != nil {
		log.Fatalf("Error getting image layers: %v", err)
	}

	for i, layer := range layers {
		fmt.Printf("Processing Layer %d...\n", i+1)
		err := extractLayer(layer, *destDir)
		if err != nil {
			log.Fatalf("Error extracting layer: %v", err)
		}
	}
	fmt.Println("Extraction complete!")
}
