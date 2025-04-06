package unhasher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mainHashService/internal/entity"
	"mainHashService/internal/usecase"

	"mainHashService/internal/utills/butcher"
	utills "mainHashService/internal/utills/mapper"
	"mainHashService/internal/utills/writer"
	"mainHashService/internal/utills/zipper"
	"mainHashService/pkg/logger"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	numWorkers  = 10
	maxFileSize = 1 << 30 // 1 ГБ в байтах
	butchSize   = 200
)

type UnhasherUCImpl struct {
	FetchRepo      FetchRepo
	S3Repo         S3Repo
	unhashEndpoint string
	unhashClient   http.Client
	FileWriter     writer.FileWriter
	FileZipper     zipper.Zipper
	lg             logger.Logger
}

func New(lg *logger.Logger, unhshEndpoint string, fetchRepo FetchRepo, s3Repo S3Repo, fileWriter writer.FileWriter, fileZipper zipper.Zipper) *UnhasherUCImpl {
	return &UnhasherUCImpl{
		lg:             *lg,
		unhashEndpoint: unhshEndpoint,
		FetchRepo:      fetchRepo,
		S3Repo:         s3Repo,
		FileWriter:     fileWriter,
		FileZipper:     fileZipper,
		unhashClient: http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

var _ usecase.UnhasherUC = (*UnhasherUCImpl)(nil)

// UnhashWorker - пул воркеров для расшифровки данных.
func (uc *UnhasherUCImpl) unhashWorker(ctx context.Context, wg *sync.WaitGroup, dataChan <-chan []UserData, resultChan chan<- entity.UnhashedData, errChan chan<- error) {
	defer wg.Done()

	for batch := range dataChan {
		select {
		case <-ctx.Done():
			return
		default:
			unhashed, err := uc.getUnhashedData(ctx, batch)
			if err != nil {
				uc.lg.Logger.Error().Msgf("failed to get unhashed data: %v", err)
				errChan <- ErrFailedUnhashing
				continue
			}
			resultChan <- unhashed
		}
	}
}

// UnhashFromQuery - достает данных из бд по звпросу от пользователя,
// делит данные на батчи по 200 элементов и через пул воркеров отправляет данные на расшифровку,
// полученные данные записывает в файл, затем полученный файл архивирует паролит и кладет в минио.
// Возвращает название бакета и url для скачивания.
func (uc *UnhasherUCImpl) UnhashFromQuery(ctx context.Context, fields []string, filters []entity.QueryStmt) (string, string, error) {

	query, args, err := uc.FetchRepo.QueryBuilder(fields, filters)
	if err != nil {
		return "", "", err
	}

	if ok, err := uc.processingQuery(query); !ok {
		return "", "", err
	}

	tempDir, err := uc.createTempDir()
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to create temp dir: %v", err)
		return "", "", err
	}
	defer os.RemoveAll(tempDir)

	fetchedData, err := uc.FetchRepo.GetHashFromQuery(ctx, query, args)
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to get hash from query: %v", err)
		return "", "", err
	}

	dataChan := make(chan []UserData, numWorkers)
	resultChan := make(chan entity.UnhashedData, numWorkers)
	errChan := make(chan error, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go uc.unhashWorker(ctx, &wg, dataChan, resultChan, errChan)
	}

	go func() {
		batches := butcher.BatchUsers(fetchedData, butchSize)
		if len(batches) == 0 {
			uc.lg.Logger.Error().Msgf("no data fetched")
			return
		}
		for _, batch := range batches {
			select {
			case dataChan <- batch:
			case <-ctx.Done():
				return
			}
		}
		close(dataChan)
	}()

	var unhashedBatches []entity.UnhashedData
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		unhashedBatches = append(unhashedBatches, result)
	}

	select {
	case err := <-errChan:
		return "", "", err
	default:
	}

	if len(unhashedBatches) == 0 {
		uc.lg.Logger.Error().Msgf("no data processed")
		return "", "", err
	}

	err = uc.runWriter(unhashedBatches, fetchedData, tempDir)
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to run writer: %v", err)
		return "", "", err
	}

	zipName, err := uc.FileZipper.Zipper(tempDir, uc.FileZipper.Password)
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to zip: %v", err)
		return "", "", err
	}

	filePath := strings.Join([]string{tempDir, zipName}, "/")
	url, bucketName, err := uc.S3Repo.UploadObject(ctx, filePath, zipName)
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to upload object: %v", err)
		return "", "", err
	}

	return url, bucketName, nil
}

// runWriter - Основной writer.
// Запускает writer который мапит структуры entity.UnhashedData в структуры UserData к финальному виду.
// Затем записывает их в файл.
func (uc *UnhasherUCImpl) runWriter(unhashedBatches []entity.UnhashedData, userData []UserData, tempDir string) error {
	writer := &uc.FileWriter
	writer.Dir = tempDir

	dataChan := make(chan ResultStruct, numWorkers)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		currentSize := 0
		var encoder *json.Encoder

		if err := writer.CreateNewFile(); err != nil {
			uc.lg.Logger.Error().Msgf("failed to create new file: %v", err)
			errChan <- ErrCreateInitialFile
			return
		}
		encoder = json.NewEncoder(writer.CurrentFile)

		for result := range dataChan {
			recordSize := len(result.Name) + len(result.Surname) + len(result.Email) + len(result.ClearPhone) + 20

			if currentSize+recordSize > maxFileSize {
				if err := writer.CreateNewFile(); err != nil {
					uc.lg.Logger.Error().Msgf("failed to create new file: %v", err)
					errChan <- ErrCreateNewFile
					return
				}
				encoder = json.NewEncoder(writer.CurrentFile)
				currentSize = 0
			}

			if err := encoder.Encode(result); err != nil {
				uc.lg.Logger.Error().Msgf("failed to encode: %v", err)
				errChan <- ErrWriteDataToFile
				return
			}

			currentSize += recordSize
		}

		if writer.CurrentFile != nil {
			if err := writer.CurrentFile.Close(); err != nil {
				uc.lg.Logger.Error().Msgf("failed to close file: %v", err)
				errChan <- err
			}
		}
	}()

	var allMappedData []utills.Mapper
	for _, data := range unhashedBatches {
		allMappedData = append(allMappedData, utills.MapperUnhash(data, userData)...)
	}

	go func() {
		for _, data := range allMappedData {
			select {
			case dataChan <- ResultStruct{
				UserID:     data.UserID,
				Name:       data.Name,
				Surname:    data.Surname,
				Email:      data.Email,
				ClearPhone: data.ClearNumber,
			}:
			case <-errChan:
				close(dataChan)
				return
			}
		}
		close(dataChan)
	}()

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

// Вспомогательный метод для валидации запроса.
func (uc *UnhasherUCImpl) processingQuery(query string) (bool, error) {
	// Разрешаем только SELECT запросы
	normalized := strings.ToUpper(strings.TrimSpace(query))
	if !strings.HasPrefix(normalized, "SELECT ") {
		return false, ErrOnlySelectAllowed
	}

	// Запрещаем опасные ключевые слова, даже в комментариях
	dangerousKeywords := []string{
		"INSERT", "UPDATE", "DELETE", "DROP", "ALTER",
		"TRUNCATE", "CREATE", "EXEC", "UNION", "--", "/*",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(normalized, keyword) {
			return false, fmt.Errorf("%w: dangerous keyword '%s' found", ErrQueryWithInjection, keyword)
		}
	}

	return true, nil
}

// Вспомогательный метод для расхеширования данных.
// Отправляет запрос на расхеширование и возвращает расхешированные данные.
func (uc *UnhasherUCImpl) getUnhashedData(ctx context.Context, data []UserData) (entity.UnhashedData, error) {
	uc.lg.Logger.Info().Msgf("len of hash: %d", len(data))
	var hashSalt []Hash
	for _, d := range data {
		hashSalt = append(hashSalt, Hash{
			UserID:      d.ID,
			PhoneNumber: d.Phone,
			Salt:        d.Salt,
		})
	}
	return uc.UnhashData(ctx, &Unhashdata{HashSalt: hashSalt, Domain: data[0].Domain})
}

// Вспомогательный метод для создания временной директории.
func (uc *UnhasherUCImpl) createTempDir() (string, error) {
	path := fmt.Sprintf("temp_%d", time.Now().UnixNano())
	tempDir, err := os.MkdirTemp("", path)
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to create temp dir: %v", err)
		return "", ErrCreeateTempDir
	}
	return tempDir, nil
}

// UnhashFromFile - метод для расхеширования данных из файла.
// Скачивает файл с minio читает его в структуру и батчами по 200 отправляет на расшифровку
// в сервис расшифровки.
func (uc *UnhasherUCImpl) UnhashFromFile(ctx context.Context, bucket, objectName string) (string, string, error) {
	objectReader, err := uc.S3Repo.DownloadObject(ctx, bucket, objectName)
	if err != nil {
		return "", "", err
	}

	fetchedData, err := uc.FetchRepo.GetHashFromFile(ctx, objectReader)
	if err != nil {
		return "", "", err
	}
	tempDir, err := uc.createTempDir()
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(tempDir)

	dataChan := make(chan []UserData, numWorkers)
	resultChan := make(chan entity.UnhashedData, numWorkers)
	errChan := make(chan error, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go uc.unhashWorker(ctx, &wg, dataChan, resultChan, errChan)
	}

	go func() {
		batches := butcher.BatchUsers(fetchedData, butchSize)
		if len(batches) == 0 {
			uc.lg.Logger.Error().Msgf("no data fetched")
			return
		}
		for _, batch := range batches {
			select {
			case dataChan <- batch:
			case <-ctx.Done():
				return
			}
		}
		close(dataChan)
	}()

	// Собираем результаты
	var unhashedBatches []entity.UnhashedData
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		unhashedBatches = append(unhashedBatches, result)
	}

	select {
	case err := <-errChan:
		return "", "", err
	default:
	}

	if len(unhashedBatches) == 0 {
		return "", "", fmt.Errorf("no data processed")
	}

	err = uc.runWriter(unhashedBatches, fetchedData, tempDir)
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to run writer: %v", err)
		return "", "", err
	}

	zipName, err := uc.FileZipper.Zipper(tempDir, uc.FileZipper.Password)
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to zip: %v", err)
		return "", "", err
	}

	filePath := strings.Join([]string{tempDir, zipName}, "/")
	url, bucketName, err := uc.S3Repo.UploadObject(ctx, filePath, zipName)
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to upload object: %v", err)
		return "", "", err
	}

	return url, bucketName, nil
}

// UplooadFile - Метод для загрузки файл на Minio
func (uc *UnhasherUCImpl) UplooadFile(ctx context.Context, reader io.ReadCloser, objectName string, size int64) (string, string, error) {
	return uc.S3Repo.UploadObjectFromFile(ctx, reader, objectName, size)
}

// Метод для получения захешированных данных.
// Необходим для тестирования.
func (uc *UnhasherUCImpl) GetHashedFile(ctx context.Context, query string) error {
	tempDir, err := uc.createTempDir()
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to create temp dir: %v", err)
		return err
	}
	data, err := uc.FetchRepo.GetHashedData(ctx, query)
	if err != nil {
		return err
	}

	err = uc.runWriterSimple(data, tempDir)
	if err != nil {
		uc.lg.Logger.Error().Msgf("failed to run writer: %v", err)
		return err
	}

	return nil
}

// Writer для того, чтобы записать захешированные данные в файл.
// нужен исключительно для тестирования.
func (uc *UnhasherUCImpl) runWriterSimple(userData []HashedData, tempDir string) error {
	// Создаем файл во временной директории
	filePath := filepath.Join(tempDir, "data.json")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	// Записываем данные в файл
	for _, data := range userData {
		// Формируем структуру для записи
		dbUser := struct {
			UserID       int64     `json:"user_id"`
			UserName     string    `json:"user_name"`
			Surname      string    `json:"surname"`
			Email        string    `json:"email"`
			HashedPhone  string    `json:"hashed_phone"`
			Salt         int64     `json:"salt"`
			DomainNumber int64     `json:"domain_number"`
			CreatedAt    time.Time `json:"created_at"`
		}{
			UserID:       data.UserID,
			UserName:     data.UserName,
			Surname:      data.Surname,
			Email:        data.Email,
			HashedPhone:  data.HashedPhone,
			Salt:         data.Salt,
			DomainNumber: data.DomainNumber,
			CreatedAt:    data.CreatedAt,
		}

		// Кодируем и записываем в файл
		if err := encoder.Encode(dbUser); err != nil {
			return fmt.Errorf("failed to write user data: %w", err)
		}
	}

	return nil
}
