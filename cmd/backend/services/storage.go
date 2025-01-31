package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func SavePDF(file io.Reader, fileName string) error {
	filePath := filepath.Join("storage", "pdfs", fileName)

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		log.Printf("error creating directory: %v", err)
		return fmt.Errorf("error creating directory: %v", err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("error creating file: %v", err)
		return fmt.Errorf("error creating file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		log.Printf("error saving file: %v", err)
		return fmt.Errorf("error saving file: %v", err)
	}

	return nil
}
