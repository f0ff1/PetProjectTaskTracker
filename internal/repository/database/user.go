package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"TaskTracker/internal/model"
)

func (s *DataBaseRepo) GetOrCreateUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*model.User, error) {

	log.Printf("🔍 Поиск пользователя с telegram_id=%d", telegramID)

	user, err := s.GetUserByTelegramID(ctx, telegramID)
	if err == nil {

		log.Printf("✅ Пользователь найден: ID=%d", user.ID)
		_ = s.UpdateUserActivity(ctx, telegramID)
		return user, nil
	}

	log.Printf("⚠️ Пользователь не найден: %v, создаю нового", err)

	query := `
        INSERT INTO users (telegram_id, username, first_name, last_name, is_admin, created_at, last_active)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, telegram_id, username, first_name, last_name, is_admin, created_at, last_active
    `

	isAdmin := telegramID == 1977074293

	var userModel model.User
	now := time.Now()

	err = s.dbPool.QueryRow(ctx, query,
		telegramID, username, firstName, lastName, isAdmin, now, now,
	).Scan(
		&userModel.ID,
		&userModel.TelegramID,
		&userModel.Username,
		&userModel.FirstName,
		&userModel.LastName,
		&userModel.IsAdmin,
		&userModel.CreatedAt,
		&userModel.LastActive,
	)

	if err != nil {
		return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	log.Printf("✅ Создан новый пользователь: ID=%d, Telegram=%d, Admin=%v",
		userModel.ID, userModel.TelegramID, userModel.IsAdmin)
	return &userModel, nil
}

// GetUserByTelegramID получает пользователя по Telegram ID
func (s *DataBaseRepo) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	query := `SELECT * FROM users WHERE telegram_id = $1`

	var user model.User
	err := s.dbPool.QueryRow(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.LastActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь не найден")
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return &user, nil
}

// UpdateUserActivity обновляет время последней активности
func (s *DataBaseRepo) UpdateUserActivity(ctx context.Context, telegramID int64) error {
	query := `UPDATE users SET last_active = CURRENT_TIMESTAMP WHERE telegram_id = $1`
	_, err := s.dbPool.Exec(ctx, query, telegramID)
	return err
}

// GetAllUsers получает всех пользователей (только для админа)
func (s *DataBaseRepo) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	query := `
        SELECT id, telegram_id, username, first_name, last_name, is_admin, created_at, last_active
        FROM users
        ORDER BY id
    `

	rows, err := s.dbPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователей: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.TelegramID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.LastActive,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// GetUserByID получает пользователя по ID
func (s *DataBaseRepo) GetUserByID(ctx context.Context, userID int) (*model.User, error) {
	query := `
        SELECT id, telegram_id, username, first_name, last_name, is_admin, created_at, last_active
        FROM users
        WHERE id = $1
    `

	var user model.User
	err := s.dbPool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.LastActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с ID %d не найден", userID)
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return &user, nil
}
