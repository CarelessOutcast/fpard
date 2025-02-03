package services

import (
	_ "embed"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// Grab the binary data
// go:embed extraction/dist/extraction.exe
var extractTextBin []byte

type Executable struct {
	binaryPath string
}

func NewExecutable() (*Executable, error) {

	log.Printf("Embedded binary size: %d bytes", len(extractTextBin))

	if len(extractTextBin) == 0 {
		return nil, errors.New("embedded binary is empty")
	}

	log.Printf("")

	// Create a stage for binary
	tmpFile, err := os.CreateTemp("", "extractText-*.exe")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close() // because OS will keep it open (non-executable)

	// Dump binary into stage
	if _, err := tmpFile.Write(extractTextBin); err != nil {
		return nil, err
	}
	log.Printf("")

	// Set permissions to execute
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		return nil, err
	}
	log.Print("temp file created")

	return &Executable{binaryPath: tmpFile.Name()}, nil
}

// TODO: Make this more abstract GetTextFrom() onyl certain file extensions
// are supported.

func (e *Executable) GetTextFromPDF(fileName string) ([]byte, error) {
	if e.binaryPath == "" {
		return nil, errors.New("Binary path is not set")
	}

	// Find the file from the storage solution
	filePath := filepath.Join("storage", "pdfs", fileName)
	log.Printf("filepath: %v", filePath)

	// Hand off to embedded python binary
	cmd := exec.Command(e.binaryPath, "../../"+filePath)

	// Use the python Binary to extract the output of the pdf
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// output is serialized as [,\d,\d,...,\d] so deserialize
	// TODO: Extract the serialization logic between python and golang
	outputStr := strings.TrimSpace(string(output))
	outputStr = strings.Trim(outputStr, "[]")
	byteStrings := strings.Split(outputStr, ",")

	var result []byte
	for _, bs := range byteStrings {
		b, err := strconv.Atoi(bs)
		if err != nil {
			return nil, err
		}
		result = append(result, byte(b))
	}

	// return the text
	return result, err
}

func (e *Executable) SetBinaryPath(path string) {
	e.binaryPath = path
}
