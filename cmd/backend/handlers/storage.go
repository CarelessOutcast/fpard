package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/carelessoutcast/fpard/cmd/backend/services"
)

const uploadDir = "./storage/pdfs"

func UploadPDF(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the pdf (security)
	err := r.ParseMultipartForm(32 << 20) // 32 mb
	if err != nil {
		http.Error(w, "Unable to parse the form", http.StatusBadRequest)
		return
	}

	file, metaData, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
	}
	defer file.Close()

	log.Printf("Upload file: %s, Size: %d bytes\n", metaData.Filename, metaData.Size)

	// create a unique filename
	fileName := generateFileName(metaData.Filename)

	// call Storage service to save the file
	err = services.SavePDF(file, fileName)

	log.Printf("File uploaded successfully: %s", fileName)
	fmt.Fprintf(w, "File uploaded successfully: %s", fileName)
}

func generateFileName(fName string) string {
	return fmt.Sprintf("%d_%s", time.Now().Unix(), fName)
}
