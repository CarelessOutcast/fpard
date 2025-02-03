package handlers

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/carelessoutcast/fpard/cmd/backend/services"
)

func ExtractPDF(w http.ResponseWriter, r *http.Request) {
	// Check what file they're asking to extract
	fileName := r.PathValue("id") // or manually parse

	// TODO: Make sure filename has valid characters

	// Check for the file
	if _, err := os.Stat(fileName); err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
	}

	// Extract the PDF
	textContent, err := services.GetTextFromPDF(fileName)
	if err != nil {
		log.Printf("File '%v' content parsed failed", fileName)
		http.Error(w, "Error extracting content", http.StatusInternalServerError)
	}
	log.Printf("File '%v' content parsed successfully", fileName)

	// Write the text to the response
	// TODO: How do I send a stream of characters rather than a full out request
	fmt.Fprintf(w, "Content: %v", textContent)
}
