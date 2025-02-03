package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func SavePDF(file io.Reader, fileName string) error {
	// TODO: Create a unique filename

	// Path to the file storage
	filePath := filepath.Join("storage", "pdfs", fileName)

	// Check if path exists
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		log.Printf("error creating directory: %v", err)
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Create the file
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("error creating file: %v", err)
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()
	// Fill the contents into the file
	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Printf("error saving file: %v", err)
		return fmt.Errorf("error saving file: %v", err)
	}

	return nil
}
