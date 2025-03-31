package fetchdata

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"mainHashService/app/repo/postgres"
	"mainHashService/pkg/logger"
)

type RepoImpl struct {
	lg logger.Logger
	db *pgxpool.Pool
}

func New(lg *logger.Logger, db *pgxpool.Pool) *RepoImpl {
	return &RepoImpl{
		lg: *lg,
		db: db,
	}
}

var _ postgres.FetchDataRepo = (*RepoImpl)(nil)

func (r *RepoImpl) GetHashFromQuery(ctx context.Context, query string) ([]postgres.UserData, error) {
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batch []postgres.UserData
	for rows.Next() {
		var data postgres.UserData
		err := rows.Scan(
			&data.ID,
			&data.Name,
			&data.Surname,
			&data.Email,
			&data.UserHash.Hash.PhoneNumber,
			&data.UserHash.Hash.Salt,
			&data.UserHash.Domain,
			&data.CreateAt,
		)
		if err != nil {
			return nil, err
		}

		batch = append(batch, data)

	}

	return batch, nil
}

func (r *RepoImpl) GetHashFromFile(ctx context.Context, reader io.ReadCloser) ([]postgres.UserData, error) {
	defer reader.Close()

	// Читаем файл построчно
	scanner := bufio.NewScanner(reader)
	var users []postgres.UserData

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue // Пропускаем пустые строки
		}

		var tempUser postgres.HashedData

		if err := json.Unmarshal(line, &tempUser); err != nil {
			return nil, err
		}

		// Преобразуем во внутреннюю структуру
		user := postgres.UserData{
			ID:      tempUser.UserID,
			Name:    tempUser.UserName,
			Surname: tempUser.Surname,
			Email:   tempUser.Email,
			UserHash: postgres.UserHash{
				Hash: postgres.Hash{
					UserID:      tempUser.UserID,
					PhoneNumber: tempUser.HashedPhone,
					Salt:        tempUser.Salt,
				},
				Domain: tempUser.DomainNumber,
			},
			CreateAt: tempUser.CreatedAt,
		}

		users = append(users, user)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Вспомогательный метод для тестирования.
func (r *RepoImpl) GetHashedData(ctx context.Context, query string) ([]postgres.HashedData, error) {
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batch []postgres.HashedData
	for rows.Next() {
		var data postgres.HashedData
		err := rows.Scan(
			&data.UserID,
			&data.UserName,
			&data.Surname,
			&data.Email,
			&data.HashedPhone,
			&data.Salt,
			&data.DomainNumber,
			&data.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		batch = append(batch, data)

	}

	return batch, nil
}
