package mapper

import (
	"mainHashService/internal/entity"
	"mainHashService/internal/repo/postgres"
	"strings"
)

type Mapper struct {
	UserID      int64
	Name        string
	Surname     string
	Email       string
	ClearNumber string
}

func MapperUnhash(unhashedData entity.UnhashedData, userData []postgres.UserData) []Mapper {
	userMap := make(map[int64]postgres.UserData)

	// Сбор пользователей
	for _, user := range userData {
		userMap[user.ID] = user
	}
	var result []Mapper
	// Обработка расхешированных данных
	for _, item := range unhashedData.HashSalt {
		clearPhone := strings.Split(item.PhoneNumber, ":")
		if user, ok := userMap[item.UserID]; ok {
			res := Mapper{
				UserID:      user.ID,
				Name:        user.Name,
				Surname:     user.Surname,
				Email:       user.Email,
				ClearNumber: clearPhone[0],
			}
			result = append(result, res)
		}

	}
	return result
}
