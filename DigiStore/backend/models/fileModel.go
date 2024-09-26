package models

import (
	"gorm.io/gorm" // Importing GORM for ORM functionalities
)

// File represents a file stored in the distributed file storage system.
// It includes fields for managing file chunks and their metadata.
type File struct {
	gorm.Model             // Embedding GORM's Model struct for basic fields like ID, CreatedAt, UpdatedAt, and DeletedAt
	FileID       string    `json:"file_id"`       // Unique identifier for the file, used to group its chunks
	ChunkID      int       `json:"chunk_id"`      // Sequential identifier for each chunk of the file
	ChunkContent []byte    `json:"chunk_content"`  // Binary content of the file chunk, stored as a byte slice
}
