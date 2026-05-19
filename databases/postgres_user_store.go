package databases

import (
	"database/sql"
	"errors"
	"fmt"
	"go-server/models"
	"go-server/repositories"
	_ "github.com/lib/pq"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(databaseUrl string) (*PostgresUserStore, error) {
	db, err := sql.Open("postgres", databaseUrl)

	if err != nil {
		return nil, err
	}

	store := &PostgresUserStore{db: db}

	if err := store.createSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *PostgresUserStore) Close() error {
	return s.db.Close()
}

func (s *PostgresUserStore) FindAll() ([]models.User, error) {
	rows, err := s.db.Query(`
		SELECT id, name, age
		FROM users
		ORDER BY id
	`)
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			return []models.User{}, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return []models.User{}, err
	}

	return users, nil
}

func (s *PostgresUserStore) FindByID(id int) (models.User, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, name, age
		FROM users
		WHERE id=$1
	`, id).Scan(&user.ID, &user.Name, &user.Age)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repositories.ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (s *PostgresUserStore) Create(user models.User) (models.User, error) {
	err := s.db.QueryRow(`
		INSERT INTO users (name, age)
		VALUES ($1, $2)
		RETURNING id
	`, user.Name, user.Age).Scan(&user.ID)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *PostgresUserStore) DeleteByID(id int) (models.User, error) {
	var user models.User
	err := s.db.QueryRow(`
		DELETE FROM users
		WHERE id=$1
		RETURNING id, name, age
	`, id).Scan(&user.ID, &user.Name, &user.Age)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, repositories.ErrUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (s *PostgresUserStore) createSchema() error {
	_, err := s.db.Exec(`
 		CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			age INTEGER NOT NULL CHECK (age > 0),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)

	if err != nil {
		return fmt.Errorf("Create user table: %w", err)
	}

	return nil
}
