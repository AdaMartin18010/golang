package oauth2

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// OIDCProvider OpenID Connect 提供者
type OIDCProvider struct {
	server      *Server
	issuer      string
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
	userStore   UserStore
}

// UserStore 用户存储接口
type UserStore interface {
	GetUser(ctx context.Context, userID string) (*UserInfo, error)
}

// UserInfo 用户信息
type UserInfo struct {
	Subject           string `json:"sub"`
	Name              string `json:"name,omitempty"`
	GivenName         string `json:"given_name,omitempty"`
	FamilyName        string `json:"family_name,omitempty"`
	MiddleName        string `json:"middle_name,omitempty"`
	Nickname          string `json:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Profile           string `json:"profile,omitempty"`
	Picture           string `json:"picture,omitempty"`
	Website           string `json:"website,omitempty"`
	Email             string `json:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
	Gender            string `json:"gender,omitempty"`
	Birthdate         string `json:"birthdate,omitempty"`
	Zoneinfo          string `json:"zoneinfo,omitempty"`
	Locale            string `json:"locale,omitempty"`
	PhoneNumber       string `json:"phone_number,omitempty"`
	PhoneVerified     bool   `json:"phone_number_verified,omitempty"`
	Address           string `json:"address,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`
}

// IDTokenClaims ID Token Claims
type IDTokenClaims struct {
	Issuer              string   `json:"iss"`
	Subject             string   `json:"sub"`
	Audience            []string `json:"aud"`
	ExpirationTime      int64    `json:"exp"`
	IssuedAt            int64    `json:"iat"`
	AuthTime            int64    `json:"auth_time,omitempty"`
	Nonce               string   `json:"nonce,omitempty"`
	AccessTokenHash     string   `json:"at_hash,omitempty"`
	CodeHash            string   `json:"c_hash,omitempty"`
	AuthenticationClass string   `json:"acr,omitempty"`
	AuthenticationMethod string `json:"amr,omitempty"`
	AuthorizedParty     string   `json:"azp,omitempty"`
	jwt.RegisteredClaims
}

// NewOIDCProvider 创建 OIDC 提供者
func NewOIDCProvider(server *Server, issuer string, privateKey *rsa.PrivateKey, userStore UserStore) (*OIDCProvider, error) {
	if issuer == "" {
		return nil, errors.New("issuer is required")
	}
	if privateKey == nil {
		return nil, errors.New("private key is required")
	}
	if userStore == nil {
		return nil, errors.New("user store is required")
	}

	return &OIDCProvider{
		server:     server,
		issuer:     issuer,
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		userStore:  userStore,
	}, nil
}

// GenerateIDToken 生成 ID Token
func (p *OIDCProvider) GenerateIDToken(ctx context.Context, userID, clientID, nonce string, accessToken string) (string, error) {
	// 获取用户信息
	userInfo, err := p.userStore.GetUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	now := time.Now()
	claims := &IDTokenClaims{
		Issuer:         p.issuer,
		Subject:        userInfo.Subject,
		Audience:       []string{clientID},
		ExpirationTime: now.Add(3600 * time.Second).Unix(), // 1 小时
		IssuedAt:       now.Unix(),
		AuthTime:       now.Unix(),
		Nonce:          nonce,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  p.issuer,
			Subject: userInfo.Subject,
			Audience: jwt.ClaimStrings{clientID},
			ExpiresAt: jwt.NewNumericDate(now.Add(3600 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// 如果提供了访问令牌，计算访问令牌哈希
	if accessToken != "" {
		claims.AccessTokenHash = p.hashToken(accessToken)
	}

	// 创建 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "default" // Key ID

	// 签名
	tokenString, err := token.SignedString(p.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateIDToken 验证 ID Token
func (p *OIDCProvider) ValidateIDToken(ctx context.Context, tokenString, clientID, nonce string) (*IDTokenClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &IDTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*IDTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 验证发行者
	if claims.Issuer != p.issuer {
		return nil, errors.New("invalid issuer")
	}

	// 验证受众
	audienceValid := false
	for _, aud := range claims.Audience {
		if aud == clientID {
			audienceValid = true
			break
		}
	}
	if !audienceValid {
		return nil, errors.New("invalid audience")
	}

	// 验证过期时间
	if time.Now().Unix() > claims.ExpirationTime {
		return nil, errors.New("token expired")
	}

	// 验证 nonce（如果提供）
	if nonce != "" && claims.Nonce != nonce {
		return nil, errors.New("invalid nonce")
	}

	return claims, nil
}

// GetUserInfo 获取用户信息
func (p *OIDCProvider) GetUserInfo(ctx context.Context, accessToken string) (*UserInfo, error) {
	// 验证访问令牌
	token, err := p.server.ValidateToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	// 获取用户信息
	userInfo, err := p.userStore.GetUser(ctx, token.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo, nil
}

// GetDiscoveryDocument 获取 Discovery 文档
func (p *OIDCProvider) GetDiscoveryDocument() *DiscoveryDocument {
	return &DiscoveryDocument{
		Issuer:                            p.issuer,
		AuthorizationEndpoint:            fmt.Sprintf("%s/oauth/authorize", p.issuer),
		TokenEndpoint:                    fmt.Sprintf("%s/oauth/token", p.issuer),
		UserInfoEndpoint:                 fmt.Sprintf("%s/userinfo", p.issuer),
		JWKSUri:                          fmt.Sprintf("%s/.well-known/jwks.json", p.issuer),
		RegistrationEndpoint:            fmt.Sprintf("%s/register", p.issuer),
		ScopesSupported:                  []string{"openid", "profile", "email", "address", "phone", "offline_access"},
		ResponseTypesSupported:           []string{"code", "id_token", "token", "id_token token"},
		ResponseModesSupported:           []string{"query", "fragment", "form_post"},
		GrantTypesSupported:             []string{"authorization_code", "client_credentials", "refresh_token"},
		ACRValuesSupported:               []string{},
		SubjectTypesSupported:           []string{"public"},
		IDTokenSigningAlgValuesSupported: []string{"RS256"},
		IDTokenEncryptionAlgValuesSupported: []string{},
		IDTokenEncryptionEncValuesSupported: []string{},
		UserInfoSigningAlgValuesSupported:   []string{"RS256"},
		UserInfoEncryptionAlgValuesSupported: []string{},
		UserInfoEncryptionEncValuesSupported: []string{},
		RequestObjectSigningAlgValuesSupported: []string{},
		RequestObjectEncryptionAlgValuesSupported: []string{},
		RequestObjectEncryptionEncValuesSupported: []string{},
		TokenEndpointAuthMethodsSupported:        []string{"client_secret_basic", "client_secret_post"},
		TokenEndpointAuthSigningAlgValuesSupported: []string{},
		DisplayValuesSupported:                    []string{},
		ClaimTypesSupported:                       []string{"normal"},
		ClaimsSupported: []string{
			"sub", "iss", "aud", "exp", "iat", "auth_time",
			"nonce", "acr", "amr", "azp", "at_hash", "c_hash",
			"name", "given_name", "family_name", "middle_name",
			"nickname", "preferred_username", "profile", "picture",
			"website", "email", "email_verified", "gender",
			"birthdate", "zoneinfo", "locale", "phone_number",
			"phone_number_verified", "address", "updated_at",
		},
		ServiceDocumentation: fmt.Sprintf("%s/docs", p.issuer),
		ClaimsLocalesSupported: []string{},
		UILocalesSupported:    []string{},
		ClaimsParameterSupported: false,
		RequestParameterSupported: false,
		RequestURIParameterSupported: false,
		RequireRequestURIRegistration: false,
		OPPolicyURI:                    fmt.Sprintf("%s/policy", p.issuer),
		OPTosURI:                       fmt.Sprintf("%s/tos", p.issuer),
	}
}

// GetJWKS 获取 JWKS (JSON Web Key Set)
func (p *OIDCProvider) GetJWKS() *JWKS {
	// 从公钥生成 JWK
	jwk := &JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: "default",
		Alg: "RS256",
	}

	// 提取 RSA 公钥的模数和指数
	n := p.publicKey.N
	e := p.publicKey.E

	// 编码为 Base64URL
	nBytes := n.Bytes()
	eBytes := make([]byte, 4)
	eBytes[0] = byte(e >> 24)
	eBytes[1] = byte(e >> 16)
	eBytes[2] = byte(e >> 8)
	eBytes[3] = byte(e)

	jwk.N = base64.RawURLEncoding.EncodeToString(nBytes)
	jwk.E = base64.RawURLEncoding.EncodeToString(eBytes)

	return &JWKS{
		Keys: []*JWK{jwk},
	}
}

// hashToken 计算令牌哈希（用于 at_hash 和 c_hash）
func (p *OIDCProvider) hashToken(token string) string {
	// 使用 SHA-256 哈希令牌的前半部分
	// 这是 OIDC 规范的要求
	hash := sha256.Sum256([]byte(token))
	halfHash := hash[:len(hash)/2]
	return base64.RawURLEncoding.EncodeToString(halfHash)
}

// DiscoveryDocument OIDC Discovery 文档
type DiscoveryDocument struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	UserInfoEndpoint                 string   `json:"userinfo_endpoint"`
	JWKSUri                          string   `json:"jwks_uri"`
	RegistrationEndpoint             string   `json:"registration_endpoint,omitempty"`
	ScopesSupported                  []string `json:"scopes_supported"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	ResponseModesSupported           []string `json:"response_modes_supported"`
	GrantTypesSupported             []string `json:"grant_types_supported"`
	ACRValuesSupported               []string `json:"acr_values_supported"`
	SubjectTypesSupported           []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	IDTokenEncryptionAlgValuesSupported []string `json:"id_token_encryption_alg_values_supported"`
	IDTokenEncryptionEncValuesSupported []string `json:"id_token_encryption_enc_values_supported"`
	UserInfoSigningAlgValuesSupported   []string `json:"userinfo_signing_alg_values_supported"`
	UserInfoEncryptionAlgValuesSupported []string `json:"userinfo_encryption_alg_values_supported"`
	UserInfoEncryptionEncValuesSupported []string `json:"userinfo_encryption_enc_values_supported"`
	RequestObjectSigningAlgValuesSupported []string `json:"request_object_signing_alg_values_supported"`
	RequestObjectEncryptionAlgValuesSupported []string `json:"request_object_encryption_alg_values_supported"`
	RequestObjectEncryptionEncValuesSupported []string `json:"request_object_encryption_enc_values_supported"`
	TokenEndpointAuthMethodsSupported        []string `json:"token_endpoint_auth_methods_supported"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported"`
	DisplayValuesSupported                    []string `json:"display_values_supported"`
	ClaimTypesSupported                       []string `json:"claim_types_supported"`
	ClaimsSupported                           []string `json:"claims_supported"`
	ServiceDocumentation                      string   `json:"service_documentation,omitempty"`
	ClaimsLocalesSupported                    []string `json:"claims_locales_supported"`
	UILocalesSupported                        []string `json:"ui_locales_supported"`
	ClaimsParameterSupported                  bool     `json:"claims_parameter_supported"`
	RequestParameterSupported                  bool     `json:"request_parameter_supported"`
	RequestURIParameterSupported              bool     `json:"request_uri_parameter_supported"`
	RequireRequestURIRegistration             bool     `json:"require_request_uri_registration"`
	OPPolicyURI                               string   `json:"op_policy_uri,omitempty"`
	OPTosURI                                  string   `json:"op_tos_uri,omitempty"`
}

// JWKS JSON Web Key Set
type JWKS struct {
	Keys []*JWK `json:"keys"`
}

// JWK JSON Web Key
type JWK struct {
	Kty string `json:"kty"` // Key Type
	Use string `json:"use"` // Public Key Use
	Kid string `json:"kid"` // Key ID
	Alg string `json:"alg"` // Algorithm
	N   string `json:"n"`   // Modulus (Base64URL)
	E   string `json:"e"`   // Exponent (Base64URL)
}

// MemoryUserStore 内存用户存储（用于测试）
type MemoryUserStore struct {
	users map[string]*UserInfo
	mu    sync.RWMutex
}

// NewMemoryUserStore 创建内存用户存储
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users: make(map[string]*UserInfo),
	}
}

// GetUser 获取用户信息
func (s *MemoryUserStore) GetUser(ctx context.Context, userID string) (*UserInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[userID]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", userID)
	}

	return user, nil
}

// Save 保存用户信息
func (s *MemoryUserStore) Save(ctx context.Context, user *UserInfo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[user.Subject] = user
	return nil
}
