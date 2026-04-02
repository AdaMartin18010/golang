package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"user-service/internal/domain"
)

func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func Migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		phone VARCHAR(20),
		email_verified BOOLEAN DEFAULT FALSE,
		avatar_url VARCHAR(500),
		last_login_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
	`

	_, err := db.Exec(schema)
	return err
}

type userRepository struct {
	db    *sql.DB
	cache Cache
}

func NewUserRepository(db *sql.DB, cache Cache) domain.UserRepository {
	return &userRepository{
		db:    db,
		cache: cache,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.AvatarURL,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	// Try cache first
	if r.cache != nil {
		if user, err := r.cache.GetUser(id.String()); err == nil {
			return user, nil
		}
	}

	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, 
		       email_verified, avatar_url, last_login_at, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName,
		&user.LastName, &user.Phone, &user.EmailVerified, &user.AvatarURL,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	// Cache the result
	if r.cache != nil {
		r.cache.SetUser(id.String(), user, 5*time.Minute)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, 
		       email_verified, avatar_url, last_login_at, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName,
		&user.LastName, &user.Phone, &user.EmailVerified, &user.AvatarURL,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET first_name = $1, last_name = $2, phone = $3, avatar_url = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`
	result, err := r.db.Exec(query, user.FirstName, user.LastName, user.Phone, user.AvatarURL, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrUserNotFound
	}

	// Invalidate cache
	if r.cache != nil {
		r.cache.DeleteUser(user.ID.String())
	}

	return nil
}

func (r *userRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrUserNotFound
	}

	// Invalidate cache
	if r.cache != nil {
		r.cache.DeleteUser(id.String())
	}

	return nil
}

func (r *userRepository) Exists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}
