package impl

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"mainHashService/app/entity"
	"mainHashService/app/repo"
	"mainHashService/pkg/logger"
	"net/http"
	"strings"
	"time"
)

type CheckerRepoImpl struct {
	lg             logger.Logger
	unhashEndpoint string
	// Ключ для определения HMAC
	SecretKey    string
	unhashClient http.Client
}

func New(lg *logger.Logger, unhashEndpoint string, secretKey string) *CheckerRepoImpl {
	return &CheckerRepoImpl{
		lg:             *lg,
		unhashEndpoint: unhashEndpoint,
		SecretKey:      secretKey,
		unhashClient: http.Client{
			Timeout: time.Second * 60,
		},
	}
}

// CheckHash - отправляет запрос в другой сервис для того, чтобы проверить,
// правильно ли пользователь вытащил данные из бд.
func (r *CheckerRepoImpl) CheckHash(ctx context.Context, hash entity.Checker) (bool, error) {
	url := "http://" + r.unhashEndpoint + "/unhash"

	data, err := json.Marshal(hash)
	if err != nil {
		r.lg.Error(err.Error())
		return false, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(data)))
	if err != nil {
		r.lg.Error(err.Error())
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.unhashClient.Do(req)
	if err != nil {
		r.lg.Error(err.Error())
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.lg.Error(err.Error())
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		r.lg.Logger.Error().Msgf("status: %d, url: %s,  body: %s", resp.StatusCode, url, string(body))
		resp.Body.Close()

		return false, repo.ErrRepoInternal
	}
	var unhashedData VerifyHash
	err = json.Unmarshal(body, &unhashedData)
	if err != nil {
		r.lg.Error(err.Error())
		return false, err
	}
	flag := verifyUnhashedData(unhashedData, r.SecretKey)
	return flag, nil
}

// Вспомогательная функция для генерации HMAC.
func generateHMAC(phoneNumber, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(phoneNumber))
	return hex.EncodeToString(h.Sum(nil))
}

// Вспомогательная функция для проверки расшифрованных данных.
// После расшифровки номер приходит в виде - номер:HMAС.
func verifyUnhashedData(unhashedData VerifyHash, secret string) bool {
	var flag bool
	var phoneNumber, hmacOriginal, hmacCalculated string
	for _, data := range unhashedData.Hash {
		parts := strings.Split(data.PhoneNumber, ":")
		if len(parts) != 2 {
			return false
		}

		phoneNumber = parts[0]
		hmacOriginal = parts[1]

		hmacCalculated = generateHMAC(phoneNumber, secret)
		flag = hmacCalculated == hmacOriginal
	}

	return flag
}
