package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Corray333/keep_it/internal/domains/user/entities"
	"github.com/Masterminds/squirrel"
)

func (s *UserRepository) GetUser(ctx context.Context, searchUser *entities.User) (*entities.User, error) {
	// Инициализируем SQL builder
	qb := squirrel.Select().From("users").PlaceholderFormat(squirrel.Dollar)

	// Определяем, какие поля будем возвращать
	var columns []string
	where := squirrel.And{} // Изменяем на And срез для условий

	// Если указан ID - ищем только по нему>and возвращаем все поля
	if searchUser.ID != 0 {
		where = append(where, squirrel.Eq{"user_id": searchUser.ID})
		columns = []string{"user_id", "username", "tg_id", "email", "avatar", "password", "ref_code"}
	} else {
		// Для остальных случаев возвращаем только id, username и avatar
		columns = []string{"user_id", "username", "avatar"}

		// Если указан TelegramID - ищем только по нему
		if searchUser.TelegramID != 0 {
			where = append(where, squirrel.Eq{"tg_id": searchUser.TelegramID})
		} else if searchUser.Username != "" || searchUser.Email != "" {
			// Ищем по username или email с OR условием
			orConditions := squirrel.Or{}
			if searchUser.Username != "" {
				orConditions = append(orConditions, squirrel.Eq{"username": searchUser.Username})
			}
			if searchUser.Email != "" {
				orConditions = append(orConditions, squirrel.Eq{"email": searchUser.Email})
			}
			where = append(where, orConditions)
		}
	}

	// Если нет условий для поиска - возвращаем ошибку
	if len(where) == 0 {
		return nil, fmt.Errorf("at least one search field must be provided")
	}

	// Собираем запрос
	query := qb.Columns(columns...).Where(where)

	// Преобразуем в SQL
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос
	user := &entities.User{}
	err = s.DB.GetContext(ctx, user, sqlStr, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return user, nil
}
