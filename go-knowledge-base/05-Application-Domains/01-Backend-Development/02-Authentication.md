# 认证与授权

> **分类**: 成熟应用领域

---

## JWT 认证

```go
import "github.com/golang-jwt/jwt/v5"

func GenerateToken(userID string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    })

    return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })
}
```

---

## Gin 中间件

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
            return
        }

        token, err := ParseToken(tokenString)
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}
```

---

## OAuth2

```go
import "golang.org/x/oauth2"

var oauthConfig = &oauth2.Config{
    ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
    RedirectURL:  "http://localhost:8080/callback",
    Scopes:       []string{"openid", "email", "profile"},
    Endpoint:     google.Endpoint,
}
```
