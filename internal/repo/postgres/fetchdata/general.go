package fetchdata

import (
	"bufio"
	"context"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"mainHashService/internal/entity"
	"mainHashService/internal/repo/postgres"
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

func (r *RepoImpl) QueryBuilder(fields []string, filters []entity.QueryStmt) (string, []interface{}, error) {
	builder := sq.Select(fields...).From("users").PlaceholderFormat(sq.Dollar)

	for _, cond := range filters {
		builder = builder.Where(sq.Expr(cond.Clause, cond.Value))
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		r.lg.Logger.Error().Msgf("build SQL failed: %v", err)
		return "", nil, err
	}

	return sql, args, nil
}

func (r *RepoImpl) GetHashFromQuery(ctx context.Context, query string, args []interface{}) ([]postgres.UserData, error) {
	r.lg.Logger.Debug().Msgf("query: %s, args: %v", query, args)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var batch []postgres.UserData
	for rows.Next() {
		var user postgres.UserData
		user, err := pgx.RowToStructByNameLax[postgres.UserData](rows)
		if err != nil {
			r.lg.Logger.Error().Msgf("failed to get unhashed data: %v", err)
			return nil, err
		}
		batch = append(batch, user)

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
			ID:       tempUser.UserID,
			Name:     tempUser.UserName,
			Surname:  tempUser.Surname,
			Email:    tempUser.Email,
			Phone:    tempUser.HashedPhone,
			Salt:     tempUser.Salt,
			Domain:   tempUser.DomainNumber,
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
