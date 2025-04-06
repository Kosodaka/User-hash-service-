package zipper

import (
	"fmt"
	"github.com/alexmullins/zip"
	"io"
	"mainHashService/pkg/logger"
	"os"
	"path/filepath"
	"time"
)

type Zipper struct {
	SourcePath string
	lg         *logger.Logger
	Password   string
}

func NewZipper(lg *logger.Logger, password string) *Zipper {
	return &Zipper{
		lg:       lg,
		Password: password,
	}
}

// архивирует файлы в директории и ставит пароль на архив.
func (z *Zipper) Zipper(sourceDir string, password string) (string, error) {
	fileName := fmt.Sprintf("unhashed-%s.zip", time.Now().Format("2006-01-02-15-04-05"))
	archivePath := filepath.Join(sourceDir, fileName)

	archiveFile, err := os.OpenFile(archivePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		z.lg.Logger.Error().Msgf("failed to open archive file: %v", err)
		return "", err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(filePath) == ".zip" {
			return nil
		}

		if filePath == archivePath {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			z.lg.Logger.Error().Msgf("failed to get relative path: %v", err)
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			z.lg.Logger.Error().Msgf("failed to create file header: %v", err)
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		if password != "" {
			header.SetPassword(password)
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			z.lg.Logger.Error().Msgf("failed to create writer: %v", err)
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			z.lg.Logger.Error().Msgf("failed to open file: %v", err)
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		os.Remove(archivePath)
		z.lg.Logger.Error().Msgf("failed to walk source directory: %v", err)
		return "", err
	}

	err = zipWriter.Close()
	if err != nil {
		os.Remove(archivePath)
		z.lg.Logger.Error().Msgf("failed to close writer: %v", err)
		return "", err
	}

	z.lg.Logger.Info().Msg("Archive created successfully")
	return fileName, nil
}
