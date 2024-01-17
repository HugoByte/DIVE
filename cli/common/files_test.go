package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestDiveFileHandler_ReadFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Write some data to the temporary file
	testData := []byte("test data")
	_, err = tempFile.Write(testData)
	if err != nil {
		t.Fatal(err)
	}

	// Instantiate the diveFileHandler
	df := NewDiveFileHandler()

	// Test ReadFile function
	readData, err := df.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Compare read data with the original data
	if string(readData) != string(testData) {
		t.Errorf("Read data doesn't match the original data. Expected: %s, Got: %s", testData, readData)
	}
}

func TestDiveFileHandler_ReadJson(t *testing.T) {
	// Create a temporary file for testing

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dirPath := filepath.Join(pwd, "output")
	err = os.MkdirAll(dirPath, 0777)
	if err != nil {
		t.Fatal(err)
	}
	filePath := filepath.Join(dirPath, "testfile.json")
	tempFile, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(fmt.Sprintf("%s/%s", pwd, "output"))

	// Write JSON data to the temporary file
	testData := map[string]string{"key": "value"}
	encodedData, err := json.Marshal(testData)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tempFile.Write(encodedData)
	if err != nil {
		t.Fatal(err)
	}

	// Close the file before reading it
	if err := tempFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Instantiate the diveFileHandler

	df := NewDiveFileHandler()

	// Test ReadJson function
	var decodedData map[string]string
	err = df.ReadJson(tempFile.Name(), &decodedData)
	if err != nil {
		t.Fatalf("Error reading JSON file: %v", err)
	}
	fmt.Println(decodedData)

	// Compare decoded data with the original data
	if decodedData["key"] != testData["key"] {
		t.Errorf("Read JSON data doesn't match the original data. Expected: %v, Got: %v", testData, decodedData)
	}
}
