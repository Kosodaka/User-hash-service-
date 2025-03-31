package butcher

import "mainHashService/app/repo/postgres"

// Разбивает данные на батчи по 200 элементов.
func BatchUsers(users []postgres.UserData, batchSize int) [][]postgres.UserData {
	if batchSize <= 0 {
		panic("batch size must be positive")
	}

	var batches [][]postgres.UserData
	total := len(users)

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batches = append(batches, users[i:end])
	}

	return batches
}
