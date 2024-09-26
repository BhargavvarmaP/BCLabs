package utils

import (
	"bytes"           // Importing bytes package for byte manipulation
	"crypto/rand"     // Importing crypto/rand package for generating random bytes
	"encoding/hex"    // Importing encoding/hex package for encoding byte slices to hexadecimal strings
	"fmt"             // Importing fmt package for formatted I/O operations
	"io/ioutil"       // Importing ioutil package for reading file content
	"math"            // Importing math package for mathematical operations
	"mime/multipart"  // Importing mime/multipart package for handling multipart file uploads
)

// generateFileID generates a random unique identifier for a file.
// Returns the generated file ID as a hexadecimal string or an error if the generation fails.
func generateFileID() (string, error) {
	bytes := make([]byte, 16) // Create a byte slice of 16 bytes
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error generating file ID: %w", err) // Return error if reading random bytes fails
	}
	return hex.EncodeToString(bytes), nil // Return the generated ID as a hexadecimal string
}

// SplitFile splits a file into chunks of a specified size for storage.
// Parameters:
// - file: The multipart file header containing the file to be split.
// Returns:
// - fileID: A unique identifier for the file.
// - chunks: A slice of byte slices representing the chunks of the file.
// - error: An error if any occurred during the process.
func SplitFile(file *multipart.FileHeader) (string, [][]byte, error) {
	f, err := file.Open() // Open the file for reading
	if err != nil {
		return "", nil, err // Return error if opening the file fails
	}
	defer f.Close() // Ensure the file is closed after the function execution

	fileID, err := generateFileID() // Generate a unique file ID
	if err != nil {
		return "", nil, err // Return error if ID generation fails
	}

	fileContent, err := ioutil.ReadAll(f) // Read the entire content of the file
	if err != nil {
		return "", nil, err // Return error if reading the file fails
	}

	const chunkSize = 1 * 1024 * 1024 // Define the chunk size as 1 MB
	totalChunks := int(math.Ceil(float64(len(fileContent)) / float64(chunkSize))) // Calculate the total number of chunks
	chunks := make([][]byte, 0, totalChunks) // Pre-allocate a slice for efficiency

	// Split the file content into chunks
	for i := 0; i < totalChunks; i++ {
		start := i * chunkSize // Calculate the starting index of the chunk
		end := (i + 1) * chunkSize // Calculate the ending index of the chunk
		if end > len(fileContent) {
			end = len(fileContent) // Ensure the ending index does not exceed file content length
		}
		chunks = append(chunks, fileContent[start:end]) // Append the chunk to the slice
	}

	return fileID, chunks, nil // Return the generated file ID and the chunks
}

// MergeChunks merges multiple byte slices (chunks) into a single byte slice.
// Parameters:
// - chunks: A slice of byte slices to be merged.
// Returns:
// - A single byte slice containing the merged data or an error if any occurred during the process.
func MergeChunks(chunks [][]byte) ([]byte, error) {
	return bytes.Join(chunks, nil), nil // Join the chunks into a single byte slice and return
}
