package writer

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	filePrefix = "unhashed_"
)

type FileWriter struct {
	CurrentFile *os.File
	CurrentSize int64
	Dir         string
	FileIndex   int
}

func NewFileWriter() *FileWriter {
	return &FileWriter{
		CurrentFile: nil,
		CurrentSize: 0,
		Dir:         "",
		FileIndex:   0,
	}
}
func (fw *FileWriter) CreateNewFile() error {
	if fw.CurrentFile != nil {
		fw.CurrentFile.Close()
	}

	fileName := fmt.Sprintf("%s%04d.json", filePrefix, fw.FileIndex)
	filePath := filepath.Join(fw.Dir, fileName)

	file, err := os.Create(filePath)
	if err != nil {

		return fmt.Errorf("failed to create file %s: %w", fileName, err)
	}

	fw.CurrentFile = file
	fw.CurrentSize = 0
	fw.FileIndex++
	return nil
}
