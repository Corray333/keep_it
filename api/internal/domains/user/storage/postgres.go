package storage

import (
	"fmt"

	"github.com/Corray333/keep_it/internal/domains/user/types"
	"github.com/Corray333/keep_it/pkg/server/auth"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserStorage struct {
	db    *sqlx.DB
	redis *redis.Client
}

// New creates a new storage and tables
func NewStorage(db *sqlx.DB, redis *redis.Client) *UserStorage {
	return &UserStorage{
		db:    db,
		redis: redis,
	}
}

// InsertUser inserts a new user into the database and returns the id
func (s *UserStorage) InsertUser(user types.User, agent string) (int, string, error) {
	passHash, err := auth.Hash(user.Password)
	if err != nil {
		return -1, "", err
	}
	user.Password = passHash

	tx, err := s.db.Beginx()
	if err != nil {
		return -1, "", err
	}

	rows := tx.QueryRow(`
		INSERT INTO users (username, email, tg_username, password, avatar, ref_code) VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id;
	`, user.Username, user.Email, user.TelegramUsername, user.Password, "/images/avatars/default_avatar.png", user.RefCode)

	if err := rows.Scan(&user.ID); err != nil {
		return -1, "", err
	}

	refresh, err := auth.CreateToken(user.ID, auth.RefreshTokenLifeTime)
	if err != nil {
		return -1, "", err
	}

	_, err = tx.Queryx(`
		INSERT INTO user_token (user_id, token) VALUES ($1, $2);
	`, user.ID, refresh)
	if err != nil {
		return -1, "", err
	}

	tx.Commit()

	return user.ID, refresh, nil
}

// LoginUser checks if the user exists and the password is correct
func (s *UserStorage) LoginUser(user types.User, agent string) (int, string, error) {
	password := user.Password

	rows := s.db.QueryRow(`
		SELECT user_id, password FROM users WHERE email = $1;
	`, user.Email)

	if err := rows.Scan(&user.ID, &user.Password); err != nil {
		return -1, "", err
	}
	if !auth.Verify(user.Password, password) {
		return -1, "", fmt.Errorf("invalid password")
	}

	// Auto update refresh token
	refresh, err := auth.CreateToken(user.ID, auth.RefreshTokenLifeTime)
	if err != nil {
		return -1, "", err
	}

	_, err = s.db.Queryx(`
		INSERT INTO user_token (user_id, token) VALUES ($1, $2) ON CONFLICT (user_id, token) DO UPDATE SET token = $3;
	`, user.ID, refresh, refresh)
	if err != nil {
		return -1, "", err
	}

	return user.ID, refresh, nil
}

// RefreshToken checks if the refresh token is valid and returns a new pair of tokens
func (s *UserStorage) RefreshToken(id int, agent, oldRefresh string) (string, string, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return "", "", err
	}

	rows := s.db.QueryRow(`
		SELECT token FROM user_token WHERE user_id = $1 AND token = $2;
	`, id, oldRefresh)

	var refresh string
	if err := rows.Scan(&refresh); err != nil {
		return "", "", err
	}
	if refresh != oldRefresh {
		fmt.Println()
		fmt.Println(refresh, " --- ", oldRefresh)
		fmt.Println()
		return "", "", fmt.Errorf("invalid refresh token")
	}

	newRefresh, err := auth.CreateToken(id, auth.RefreshTokenLifeTime)
	if err != nil {
		return "", "", err
	}

	newAccess, err := auth.CreateToken(id, auth.AccessTokenLifeTime)
	if err != nil {
		return "", "", err
	}

	_, err = tx.Queryx(`
		UPDATE user_token SET token = $1 WHERE user_id = $2 AND token = $3;
	`, newRefresh, id, oldRefresh)
	if err != nil {
		return "", "", err
	}

	tx.Commit()

	return newAccess, newRefresh, nil
}

func (s *UserStorage) SelectUser(id string) (types.User, error) {
	var user types.User
	rows, err := s.db.Queryx(`
		SELECT * FROM users WHERE user_id = $1;
	`, id)
	if err != nil {
		return user, err
	}
	if !rows.Next() {
		return user, fmt.Errorf("user not found")
	}
	if err := rows.StructScan(&user); err != nil {
		return user, err
	}
	user.Password = ""
	return user, nil
}

func (s *UserStorage) UpdateUser(user types.User) error {
	_, err := s.db.Queryx(`
		UPDATE users SET username = $1, email = $2, avatar = $3 WHERE user_id = $4;
	`, user.Username, user.Email, user.Avatar, user.ID)
	return err
}

func (s *UserStorage) CheckUsername(username string) (bool, error) {
	res := false
	err := s.db.QueryRow("SELECT COUNT(*) > 0 AS exists FROM users WHERE username = $1", username).Scan(&res)
	if err != nil {
		return false, err
	}
	return res, nil
}
