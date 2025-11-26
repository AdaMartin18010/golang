package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken 无效的令牌
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken 过期的令牌
	ErrExpiredToken = errors.New("expired token")
	// ErrInvalidSigningMethod 无效的签名方法
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

// Config JWT配置
type Config struct {
	SecretKey       string        // 密钥（HMAC）
	PrivateKey      *rsa.PrivateKey // 私钥（RSA）
	PublicKey       *rsa.PublicKey  // 公钥（RSA）
	SigningMethod   string        // 签名方法 (HS256, RS256)
	AccessTokenTTL  time.Duration // Access Token 过期时间
	RefreshTokenTTL time.Duration // Refresh Token 过期时间
	Issuer          string        // 签发者
	Audience        string        // 受众
}

// Claims JWT Claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Email    string   `json:"email,omitempty"`
	jwt.RegisteredClaims
}

// JWT JWT管理器
type JWT struct {
	config Config
}

// NewJWT 创建JWT管理器
func NewJWT(config Config) (*JWT, error) {
	if config.SigningMethod == "" {
		config.SigningMethod = "HS256"
	}
	if config.AccessTokenTTL == 0 {
		config.AccessTokenTTL = 15 * time.Minute
	}
	if config.RefreshTokenTTL == 0 {
		config.RefreshTokenTTL = 7 * 24 * time.Hour
	}

	return &JWT{config: config}, nil
}

// GenerateAccessToken 生成Access Token
func (j *JWT) GenerateAccessToken(userID, username string, roles []string, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.config.Issuer,
			Audience:  []string{j.config.Audience},
		},
	}

	return j.signToken(claims)
}

// GenerateRefreshToken 生成Refresh Token
func (j *JWT) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.config.Issuer,
			Audience:  []string{j.config.Audience},
		},
	}

	return j.signToken(claims)
}

// ValidateToken 验证Token
func (j *JWT) ValidateToken(tokenString string) (*Claims, error) {
	token, err := j.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// 检查是否过期
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

// RefreshToken 刷新Token
func (j *JWT) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// 生成新的Access Token和Refresh Token
	accessToken, err := j.GenerateAccessToken(claims.UserID, claims.Username, claims.Roles, claims.Email)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := j.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// signToken 签名Token
func (j *JWT) signToken(claims Claims) (string, error) {
	var token *jwt.Token

	switch j.config.SigningMethod {
	case "HS256", "HS384", "HS512":
		if j.config.SecretKey == "" {
			return "", errors.New("secret key is required for HMAC signing")
		}
		token = jwt.NewWithClaims(jwt.GetSigningMethod(j.config.SigningMethod), claims)
		return token.SignedString([]byte(j.config.SecretKey))

	case "RS256", "RS384", "RS512":
		if j.config.PrivateKey == nil {
			return "", errors.New("private key is required for RSA signing")
		}
		token = jwt.NewWithClaims(jwt.GetSigningMethod(j.config.SigningMethod), claims)
		return token.SignedString(j.config.PrivateKey)

	default:
		return "", ErrInvalidSigningMethod
	}
}

// parseToken 解析Token
func (j *JWT) parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			if j.config.SecretKey == "" {
				return nil, errors.New("secret key is required for HMAC verification")
			}
			return []byte(j.config.SecretKey), nil
		}

		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			if j.config.PublicKey == nil {
				return nil, errors.New("public key is required for RSA verification")
			}
			return j.config.PublicKey, nil
		}

		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	})
}

// ExtractUserID 从Token中提取UserID
func (j *JWT) ExtractUserID(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// ExtractRoles 从Token中提取Roles
func (j *JWT) ExtractRoles(tokenString string) ([]string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	return claims.Roles, nil
}
