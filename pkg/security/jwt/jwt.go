package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager JWT 令牌管理器
type TokenManager struct {
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
	issuer           string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	signingMethod    jwt.SigningMethod
}

// Config JWT 配置
type Config struct {
	PrivateKeyPath  string        // RSA 私钥路径
	PublicKeyPath   string        // RSA 公钥路径
	Issuer          string        // 令牌签发者
	AccessTokenTTL  time.Duration // 访问令牌有效期 (推荐: 15分钟)
	RefreshTokenTTL time.Duration // 刷新令牌有效期 (推荐: 7天)
	SigningMethod   string        // 签名算法 (RS256/ES256)
}

// Claims JWT 标准声明
type Claims struct {
	jwt.RegisteredClaims
	UserID   string   `json:"uid"`
	Username string   `json:"username,omitempty"`
	Email    string   `json:"email,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	Scope    []string `json:"scope,omitempty"`
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"` // 秒
	IssuedAt     time.Time `json:"issued_at"`
}

// NewTokenManager 创建令牌管理器
func NewTokenManager(cfg Config) (*TokenManager, error) {
	// 设置默认值
	if cfg.AccessTokenTTL == 0 {
		cfg.AccessTokenTTL = 15 * time.Minute
	}
	if cfg.RefreshTokenTTL == 0 {
		cfg.RefreshTokenTTL = 7 * 24 * time.Hour
	}
	if cfg.Issuer == "" {
		cfg.Issuer = "golang-service"
	}
	if cfg.SigningMethod == "" {
		cfg.SigningMethod = "RS256"
	}

	tm := &TokenManager{
		issuer:           cfg.Issuer,
		accessTokenTTL:   cfg.AccessTokenTTL,
		refreshTokenTTL:  cfg.RefreshTokenTTL,
		signingMethod:    jwt.GetSigningMethod(cfg.SigningMethod),
	}

	// 加载密钥
	if cfg.PrivateKeyPath != "" && cfg.PublicKeyPath != "" {
		if err := tm.LoadKeysFromFile(cfg.PrivateKeyPath, cfg.PublicKeyPath); err != nil {
			return nil, fmt.Errorf("failed to load keys: %w", err)
		}
	} else {
		// 生成临时密钥（仅用于开发）
		if err := tm.GenerateKeys(2048); err != nil {
			return nil, fmt.Errorf("failed to generate keys: %w", err)
		}
	}

	return tm, nil
}

// GenerateAccessToken 生成访问令牌
func (tm *TokenManager) GenerateAccessToken(userID, username, email string, roles []string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:   userID,
		Username: username,
		Email:    email,
		Roles:    roles,
	}

	token := jwt.NewWithClaims(tm.signingMethod, claims)
	return token.SignedString(tm.privateKey)
}

// GenerateRefreshToken 生成刷新令牌
func (tm *TokenManager) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.refreshTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(tm.signingMethod, claims)
	return token.SignedString(tm.privateKey)
}

// GenerateTokenPair 生成令牌对
func (tm *TokenManager) GenerateTokenPair(userID, username, email string, roles []string) (*TokenPair, error) {
	accessToken, err := tm.GenerateAccessToken(userID, username, email, roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := tm.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(tm.accessTokenTTL.Seconds()),
		IssuedAt:     time.Now(),
	}, nil
}

// ValidateToken 验证令牌
func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if token.Method != tm.signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	return claims, nil
}

// RefreshAccessToken 刷新访问令牌
func (tm *TokenManager) RefreshAccessToken(refreshToken string) (*TokenPair, error) {
	// 验证刷新令牌
	claims, err := tm.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 生成新的令牌对
	return tm.GenerateTokenPair(claims.UserID, claims.Username, claims.Email, claims.Roles)
}

// LoadKeysFromFile 从文件加载密钥
func (tm *TokenManager) LoadKeysFromFile(privateKeyPath, publicKeyPath string) error {
	// 加载私钥
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read private key: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyData)
	if privateKeyBlock == nil {
		return errors.New("failed to decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		// 尝试 PKCS8 格式
		key, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return errors.New("private key is not RSA")
		}
	}
	tm.privateKey = privateKey

	// 加载公钥
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyData)
	if publicKeyBlock == nil {
		return errors.New("failed to decode public key PEM")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return errors.New("public key is not RSA")
	}
	tm.publicKey = publicKey

	return nil
}

// GenerateKeys 生成密钥对（仅用于开发/测试）
func (tm *TokenManager) GenerateKeys(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	tm.privateKey = privateKey
	tm.publicKey = &privateKey.PublicKey
	return nil
}

// ExportKeys 导出密钥到 PEM 格式
func (tm *TokenManager) ExportKeys() (privateKeyPEM, publicKeyPEM []byte, err error) {
	if tm.privateKey == nil || tm.publicKey == nil {
		return nil, nil, errors.New("keys not initialized")
	}

	// 导出私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(tm.privateKey)
	privateKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// 导出公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(tm.publicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	publicKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return privateKeyPEM, publicKeyPEM, nil
}
