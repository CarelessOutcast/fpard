package services

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// TODO: Clean this up with edge cases
func writeWavFile(filename string, pcmData []byte, sampleRate int) error {
	const DEFAULT_SAMPLE_RATE int = 22050
	const DEFAULT_FILE_NAME string = "test_output.wav"

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create Wav header
	var (
		channels      = uint16(1)  // Mono
		bitPerSample  = uint16(16) // 16-bit
		audioFormat   = uint16(1)  // PCM
		byteRate      = uint32(sampleRate * int(channels) * int(bitPerSample) / 8)
		blockAlign    = uint16(channels * bitPerSample / 8)
		dataChunkSize = uint32(len(pcmData))
		riffChunkSize = dataChunkSize + 36
		fmtChunkSize  = uint32(16)
	)

	// Write the RIFF header
	file.Write([]byte("RIFF"))
	binary.Write(file, binary.LittleEndian, riffChunkSize)
	file.Write([]byte("WAVE"))

	// Write fmt chunk
	file.Write([]byte("fmt "))
	binary.Write(file, binary.LittleEndian, fmtChunkSize)
	binary.Write(file, binary.LittleEndian, audioFormat)
	binary.Write(file, binary.LittleEndian, channels)
	binary.Write(file, binary.LittleEndian, uint32(sampleRate))
	binary.Write(file, binary.LittleEndian, byteRate)
	binary.Write(file, binary.LittleEndian, blockAlign)
	binary.Write(file, binary.LittleEndian, bitPerSample)

	// Write data chunk
	file.Write([]byte("data"))
	binary.Write(file, binary.LittleEndian, dataChunkSize)
	file.Write(pcmData)

	return nil
}

func int16ToBytes(int16buffer []int16) []byte {
	byteBuffer := new(bytes.Buffer)
	for _, num := range int16buffer {
		err := binary.Write(byteBuffer, binary.LittleEndian, num)
		if err != nil {
			return nil
		}
	}
	return byteBuffer.Bytes()
}

// TODO: Generalize this because header actually ends at 'data'
func removeHeader(wavData []byte) ([]byte, error) {
	return wavData[44:], nil
}

// TODO: Simplify this test
func TestTextToSpeech(t *testing.T) {
	// Loading file
	t.Log("Loading file")
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Couldn't get current working directory: %v", err)
	}
	fpath := filepath.Join(cwd, "tts", "lib")
	fTestText := filepath.Join(fpath, "test.txt")
	fTestWav := filepath.Join(fpath, "test.wav")

	// ---- Get test text ----
	t.Logf("Loading file: %s", fTestText)

	data, err := os.ReadFile(fTestText)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	testText := string(data)
	t.Logf("Extracted : %s", testText)

	// ---- Get expected wav output ----
	t.Logf("Loading file: %s", fTestWav)
	data, err = os.ReadFile(fTestWav)
	expected, err := removeHeader(data)
	if err != nil {
		t.Fatalf("Couldn't remove header: %v", err)
	}

	// ---- Run the test ----
	t.Log("Starting the test")

	// Start Espeak
	espeak, err := NewESpeak()
	if err != nil {
		t.Fatalf("Unable to start espeak file: %v", err)
	}
	err = espeak.SetVoice("en")
	if err != nil {
		t.Fatalf("Unable to add espeak voice: %v", err)
	}
	serviceOutput, err := espeak.Speak(testText)
	if err != nil {
		t.Fatalf("Unable to synthesize voice: %v", err)
	}
	actual := int16ToBytes(serviceOutput)

	if !bytes.Equal(actual, expected) {
		t.Fatalf("Actual: %d; Expected: %d;", len(actual), len(expected))
	}
}
