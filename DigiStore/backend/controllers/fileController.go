package controllers

import (
  "fmt"
  "Distributed_File_Storage/models" // Importing models package for database interactions
  "Distributed_File_Storage/utils"  // Importing utils package for file handling utilities
  "github.com/gin-gonic/gin"        // Importing Gin framework for HTTP handling
  "net/http"                        // Importing net/http package for HTTP status codes
  "sync"                            // Importing sync package for goroutine synchronization
)

// UploadFile handles the file upload process.
func UploadFile(c *gin.Context) {
  // Retrieve the file from the form input
  file, err := c.FormFile("file")
  if err != nil {
    // Return a bad request status if no file is received
    c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
    return
  }

  // Split the file into chunks for processing
  fileID, chunks, err := utils.SplitFile(file)
  if err != nil {
    // Return an internal server error if file splitting fails
    c.JSON(http.StatusInternalServerError, gin.H{"error": "File splitting failed"})
    return
  }

  // WaitGroup to synchronize goroutines
  var wg sync.WaitGroup
  // Channel for error communication between goroutines
  var errors = make(chan error, len(chunks)) 

  // Loop through each chunk and upload it concurrently
  for idx, chunk := range chunks {
    wg.Add(1) // Increment the WaitGroup counter
    go func(i int, content []byte) {
      defer wg.Done() // Decrement the counter when the goroutine completes
      // Create a record for the chunk in the database
      chunkRecord := models.File{FileID: fileID, ChunkID: i, ChunkContent: content}
      if err := models.DB.Create(&chunkRecord).Error; err != nil {
        // Send error to the channel if database operation fails
        errors <- err 
        return
      }
    }(idx, chunk) // Pass the current index and chunk content to the goroutine
  }

  wg.Wait() // Wait for all goroutines to finish

  // Check for any errors from the goroutines
  select {
  case err := <-errors:
    // If an error occurred, log it and return an error response
    fmt.Println("Error from goroutine:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file: " + err.Error()})
    return
  default:
    // No errors in the channel, proceed
  }

  // Return a success response with the file ID
  c.JSON(http.StatusOK, gin.H{"file_id": fileID})
}

// GetFiles retrieves the list of uploaded files from the database.
func GetFiles(c *gin.Context) {
  var files []models.File
  // Query the database for all files
  if err := models.DB.Find(&files).Error; err != nil {
      // Return an internal server error if retrieval fails
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving files"})
      return
  }
  // Return a success response with the list of files
  c.JSON(http.StatusOK, gin.H{"files": files})
}

// DownloadFile handles the file download process by merging its chunks.
func DownloadFile(c *gin.Context) {
  fileID := c.Param("id") // Retrieve the file ID from the URL parameter

  // Sanitize and validate fileID (implementation details not shown)

  var fileChunks []models.File
  // Query the database for the chunks of the specified file, ordered by chunk ID
  if err := models.DB.Where("file_id = ?", fileID).Order("chunk_id asc").Find(&fileChunks).Error; err != nil {
    // Return an internal server error if retrieval fails
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving file chunks"})
    return
  }

  // Create a slice to hold the content of the chunks
  chunks := make([][]byte, len(fileChunks))
  for i, chunk := range fileChunks {
    chunks[i] = chunk.ChunkContent // Populate the slice with chunk content
  }

  // Merge the chunks into a single file
  mergedFile, err := utils.MergeChunks(chunks)
  if err != nil {
    // Return an internal server error if merging fails
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error merging chunks"})
    return
  }

  // Set the Content-Disposition header for file download
  c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileID))
  // Send the merged file as a response
  c.Data(http.StatusOK, "application/octet-stream", mergedFile)
}
