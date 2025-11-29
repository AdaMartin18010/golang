package oauth2

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// PostgresTokenStore PostgreSQL 令牌存储
type PostgresTokenStore struct {
	db *sql.DB
}

// NewPostgresTokenStore 创建 PostgreSQL 令牌存储
func NewPostgresTokenStore(db *sql.DB) (*PostgresTokenStore, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	store := &PostgresTokenStore{db: db}

	// 初始化表结构
	if err := store.initSchema(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema 初始化数据库表结构
func (s *PostgresTokenStore) initSchema(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS oauth2_tokens (
			access_token TEXT PRIMARY KEY,
			refresh_token TEXT UNIQUE,
			token_type TEXT NOT NULL,
			expires_in BIGINT NOT NULL,
			scope TEXT,
			client_id TEXT NOT NULL,
			user_id TEXT,
			created_at TIMESTAMP NOT NULL,
			expires_at TIMESTAMP NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_client_id ON oauth2_tokens(client_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON oauth2_tokens(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_refresh_token ON oauth2_tokens(refresh_token)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_expires_at ON oauth2_tokens(expires_at)`,
		`CREATE TABLE IF NOT EXISTS oauth2_clients (
			id TEXT PRIMARY KEY,
			secret TEXT NOT NULL,
			redirect_uris JSONB,
			grant_types JSONB,
			scopes JSONB,
			created_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS oauth2_auth_codes (
			code TEXT PRIMARY KEY,
			client_id TEXT NOT NULL,
			redirect_uri TEXT NOT NULL,
			scope TEXT,
			user_id TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			expires_at TIMESTAMP NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_auth_codes_client_id ON oauth2_auth_codes(client_id)`,
		`CREATE INDEX IF NOT EXISTS idx_auth_codes_expires_at ON oauth2_auth_codes(expires_at)`,
	}

	for _, query := range queries {
		if _, err := s.db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

// Save 保存令牌
func (s *PostgresTokenStore) Save(ctx context.Context, token *Token) error {
	query := `INSERT INTO oauth2_tokens 
		(access_token, refresh_token, token_type, expires_in, scope, client_id, user_id, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (access_token) DO UPDATE SET
			refresh_token = EXCLUDED.refresh_token,
			expires_at = EXCLUDED.expires_at`

	_, err := s.db.ExecContext(ctx, query,
		token.AccessToken,
		token.RefreshToken,
		token.TokenType,
		token.ExpiresIn,
		token.Scope,
		token.ClientID,
		token.UserID,
		token.CreatedAt,
		token.ExpiresAt,
	)

	return err
}

// Get 获取令牌
func (s *PostgresTokenStore) Get(ctx context.Context, accessToken string) (*Token, error) {
	query := `SELECT access_token, refresh_token, token_type, expires_in, scope, 
		client_id, user_id, created_at, expires_at
		FROM oauth2_tokens WHERE access_token = $1`

	var token Token
	err := s.db.QueryRowContext(ctx, query, accessToken).Scan(
		&token.AccessToken,
		&token.RefreshToken,
		&token.TokenType,
		&token.ExpiresIn,
		&token.Scope,
		&token.ClientID,
		&token.UserID,
		&token.CreatedAt,
		&token.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Delete 删除令牌
func (s *PostgresTokenStore) Delete(ctx context.Context, accessToken string) error {
	query := `DELETE FROM oauth2_tokens WHERE access_token = $1`
	_, err := s.db.ExecContext(ctx, query, accessToken)
	return err
}

// DeleteByRefreshToken 通过刷新令牌删除
func (s *PostgresTokenStore) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	query := `DELETE FROM oauth2_tokens WHERE refresh_token = $1`
	_, err := s.db.ExecContext(ctx, query, refreshToken)
	return err
}

// CleanupExpired 清理过期令牌
func (s *PostgresTokenStore) CleanupExpired(ctx context.Context) error {
	query := `DELETE FROM oauth2_tokens WHERE expires_at < NOW()`
	_, err := s.db.ExecContext(ctx, query)
	return err
}

// PostgresClientStore PostgreSQL 客户端存储
type PostgresClientStore struct {
	db *sql.DB
}

// NewPostgresClientStore 创建 PostgreSQL 客户端存储
func NewPostgresClientStore(db *sql.DB) (*PostgresClientStore, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	store := &PostgresClientStore{db: db}

	// 初始化表结构（如果还没有创建）
	if err := store.initSchema(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema 初始化数据库表结构
func (s *PostgresClientStore) initSchema(ctx context.Context) error {
	// 表结构已在 PostgresTokenStore 中创建，这里只需要确保存在
	return nil
}

// Get 获取客户端
func (s *PostgresClientStore) Get(ctx context.Context, clientID string) (*Client, error) {
	query := `SELECT id, secret, redirect_uris, grant_types, scopes, created_at
		FROM oauth2_clients WHERE id = $1`

	var client Client
	var redirectURIsJSON, grantTypesJSON, scopesJSON []byte

	err := s.db.QueryRowContext(ctx, query, clientID).Scan(
		&client.ID,
		&client.Secret,
		&redirectURIsJSON,
		&grantTypesJSON,
		&scopesJSON,
		&client.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrClientNotFound
	}
	if err != nil {
		return nil, err
	}

	// 解析 JSON 字段
	if err := json.Unmarshal(redirectURIsJSON, &client.RedirectURIs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal redirect_uris: %w", err)
	}

	var grantTypesStr []string
	if err := json.Unmarshal(grantTypesJSON, &grantTypesStr); err != nil {
		return nil, fmt.Errorf("failed to unmarshal grant_types: %w", err)
	}
	client.GrantTypes = make([]GrantType, len(grantTypesStr))
	for i, gt := range grantTypesStr {
		client.GrantTypes[i] = GrantType(gt)
	}

	if err := json.Unmarshal(scopesJSON, &client.Scopes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scopes: %w", err)
	}

	return &client, nil
}

// ValidateSecret 验证客户端密钥
func (s *PostgresClientStore) ValidateSecret(ctx context.Context, clientID, secret string) error {
	query := `SELECT secret FROM oauth2_clients WHERE id = $1`
	var storedSecret string

	err := s.db.QueryRowContext(ctx, query, clientID).Scan(&storedSecret)
	if err == sql.ErrNoRows {
		return ErrClientNotFound
	}
	if err != nil {
		return err
	}

	if storedSecret != secret {
		return ErrInvalidClientSecret
	}

	return nil
}

// Save 保存客户端
func (s *PostgresClientStore) Save(ctx context.Context, client *Client) error {
	redirectURIsJSON, _ := json.Marshal(client.RedirectURIs)
	
	grantTypesStr := make([]string, len(client.GrantTypes))
	for i, gt := range client.GrantTypes {
		grantTypesStr[i] = string(gt)
	}
	grantTypesJSON, _ := json.Marshal(grantTypesStr)
	
	scopesJSON, _ := json.Marshal(client.Scopes)

	query := `INSERT INTO oauth2_clients (id, secret, redirect_uris, grant_types, scopes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			secret = EXCLUDED.secret,
			redirect_uris = EXCLUDED.redirect_uris,
			grant_types = EXCLUDED.grant_types,
			scopes = EXCLUDED.scopes`

	_, err := s.db.ExecContext(ctx, query,
		client.ID,
		client.Secret,
		redirectURIsJSON,
		grantTypesJSON,
		scopesJSON,
		client.CreatedAt,
	)

	return err
}

// PostgresCodeStore PostgreSQL 授权码存储
type PostgresCodeStore struct {
	db *sql.DB
}

// NewPostgresCodeStore 创建 PostgreSQL 授权码存储
func NewPostgresCodeStore(db *sql.DB) (*PostgresCodeStore, error) {
	if db == nil {
		return nil, errors.New("database connection is required")
	}

	store := &PostgresCodeStore{db: db}

	// 初始化表结构（如果还没有创建）
	if err := store.initSchema(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema 初始化数据库表结构
func (s *PostgresCodeStore) initSchema(ctx context.Context) error {
	// 表结构已在 PostgresTokenStore 中创建，这里只需要确保存在
	return nil
}

// Save 保存授权码
func (s *PostgresCodeStore) Save(ctx context.Context, code *AuthCode) error {
	query := `INSERT INTO oauth2_auth_codes 
		(code, client_id, redirect_uri, scope, user_id, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (code) DO UPDATE SET
			expires_at = EXCLUDED.expires_at`

	_, err := s.db.ExecContext(ctx, query,
		code.Code,
		code.ClientID,
		code.RedirectURI,
		code.Scope,
		code.UserID,
		code.CreatedAt,
		code.ExpiresAt,
	)

	return err
}

// Get 获取授权码
func (s *PostgresCodeStore) Get(ctx context.Context, code string) (*AuthCode, error) {
	query := `SELECT code, client_id, redirect_uri, scope, user_id, created_at, expires_at
		FROM oauth2_auth_codes WHERE code = $1`

	var authCode AuthCode
	err := s.db.QueryRowContext(ctx, query, code).Scan(
		&authCode.Code,
		&authCode.ClientID,
		&authCode.RedirectURI,
		&authCode.Scope,
		&authCode.UserID,
		&authCode.CreatedAt,
		&authCode.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrCodeNotFound
	}
	if err != nil {
		return nil, err
	}

	return &authCode, nil
}

// Delete 删除授权码
func (s *PostgresCodeStore) Delete(ctx context.Context, code string) error {
	query := `DELETE FROM oauth2_auth_codes WHERE code = $1`
	_, err := s.db.ExecContext(ctx, query, code)
	return err
}

// CleanupExpired 清理过期授权码
func (s *PostgresCodeStore) CleanupExpired(ctx context.Context) error {
	query := `DELETE FROM oauth2_auth_codes WHERE expires_at < NOW()`
	_, err := s.db.ExecContext(ctx, query)
	return err
}
