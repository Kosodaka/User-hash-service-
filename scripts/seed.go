package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/lib/pq"
	"math/rand"
	"time"
)

// Скрипт для заполнения бд

const (
	host       = "localhost"
	port       = 5432
	user       = "postgres"
	password   = "pass"
	dbname     = "user_data"
	hmacSecret = "hash_111" // Замените на ваш реальный HMAC секрет
	domain     = 1
)

func main() {
	// Подключение к БД
	connStr := "host=db port=5432 user=postgres password=pass dbname=user_data sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Генерация 2000 пользователей
	for i := 4001; i < 2000000; i++ {
		// Генерация данных
		userName := fmt.Sprintf("User%d", i)
		surname := fmt.Sprintf("Surname%d", rand.Intn(100))
		email := fmt.Sprintf("user%d@example.com", i)
		phone := fmt.Sprintf("+79%09d", rand.Intn(1000000000)) // Российский номер формата +79XXXXXXXXX
		salt := rand.Intn(10000)
		domainNumber := 1
		createdAt := time.Now().Add(-time.Duration(rand.Intn(365)) * 24 * time.Hour)

		// Хэширование телефона
		hashedPhone := hashPhoneNumber(phone, salt, domainNumber)

		// Вставка в БД
		_, err := db.Exec(`
			INSERT INTO users (user_name, surname, email, hashed_phone, salt, domain_number, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			userName, surname, email, hashedPhone, salt, domainNumber, createdAt)

		if err != nil {
			fmt.Printf("Ошибка вставки пользователя %d: %v\n", i, err)
			continue
		}
	}

	fmt.Println("Успешно добавлено 2000 пользователей")
}

func generateHMAC(phoneNumber, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(phoneNumber))
	return hex.EncodeToString(h.Sum(nil))
}

func prepareDataForHashing(data string, secret string) string {
	hmacHash := generateHMAC(data, secret)
	return fmt.Sprintf("%s:%s", data, hmacHash)
}

func hashPhoneNumber(phone string, salt int, domain int) string {
	prepared := prepareDataForHashing(phone, hmacSecret)
	hashBytes := []byte(prepared)

	hashedBytes := make([]byte, len(hashBytes))
	for i, b := range hashBytes {
		hashedBytes[i] = b ^ byte(salt) ^ byte(domain)
	}

	return hex.EncodeToString(hashedBytes)
}
