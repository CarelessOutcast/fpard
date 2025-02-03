package services

import (
	"bytes"
	"os"
	"testing"
)

func TestGetTextFromPDF_Mock(t *testing.T) {
	const mockBinPath = "D:/Coding/projects/fpard/cmd/backend/services/extraction/dist/extraction.exe"
	if _, err := os.Stat(mockBinPath); os.IsNotExist(err) {
		t.Fatal(err)
	}

	bin := &Executable{}
	bin.SetBinaryPath(mockBinPath)

	fileName := "test.pdf"
	actual, err := bin.GetTextFromPDF(fileName)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("This is a small test document! \n")
	if !bytes.Equal(expected, actual) {
		t.Fatalf("Extected: %s; Actual: %s;", expected, actual)
	}

	t.Logf("Extected text: %s", expected)
	t.Logf("Actual text: %s", actual)

}
