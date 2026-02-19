package models

import "time"

type File struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FileName  string    `json:"fileName"`
	FilePath  string    `json:"filePath"`
	FileSize  int64     `json:"fileSize"`
	MimeType  string    `json:"mimeType"`
	CreatedAt time.Time `json:"createdAt"`
}
