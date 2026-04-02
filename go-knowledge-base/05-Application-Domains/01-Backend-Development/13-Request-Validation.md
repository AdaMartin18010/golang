# 请求验证 (Request Validation)

> **分类**: 成熟应用领域  
> **标签**: #validation #api #security

---

## 基本验证

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"gte=0,lte=150"`
    Password string `json:"password" validate:"required,min=8,containsany=!@#$%"`
    Phone    string `json:"phone" validate:"e164"`
    Website  string `json:"website" validate:"omitempty,url"`
}

func ValidateRequest(req interface{}) error {
    return validate.Struct(req)
}
```

---

## 自定义验证器

```go
// 注册自定义验证器
func init() {
    validate.RegisterValidation("strongpassword", strongPasswordValidator)
    validate.RegisterValidation("notcommon", notCommonPasswordValidator)
}

func strongPasswordValidator(fl validator.FieldLevel) bool {
    password := fl.Field().String()
    
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)
    
    return hasUpper && hasLower && hasNumber && hasSpecial
}

// 使用
type RegisterRequest struct {
    Password string `validate:"required,min=12,strongpassword"`
}
```

---

## 结构化验证错误

```go
type ValidationError struct {
    Field   string `json:"field"`
    Rule    string `json:"rule"`
    Message string `json:"message"`
}

func FormatValidationErrors(err error) []ValidationError {
    var errors []ValidationError
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            errors = append(errors, ValidationError{
                Field:   e.Field(),
                Rule:    e.Tag(),
                Message: formatErrorMessage(e),
            })
        }
    }
    
    return errors
}

func formatErrorMessage(e validator.FieldError) string {
    switch e.Tag() {
    case "required":
        return fmt.Sprintf("%s is required", e.Field())
    case "min":
        return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
    case "max":
        return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
    case "email":
        return fmt.Sprintf("%s must be a valid email address", e.Field())
    default:
        return fmt.Sprintf("%s is invalid", e.Field())
    }
}
```

---

## Gin 集成

```go
func ValidationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 已绑定到请求对象
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(400, gin.H{
                "error": "validation failed",
                "details": FormatValidationErrors(err),
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// 使用
r.POST("/users", ValidationMiddleware(), createUserHandler)
```

---

## 跨字段验证

```go
type PasswordResetRequest struct {
    Password        string `json:"password" validate:"required"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type DateRangeRequest struct {
    StartDate time.Time `json:"start_date" validate:"required"`
    EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}
```

---

## 安全验证

```go
type SecurityHeaders struct {
    UserAgent    string `header:"User-Agent" validate:"required,notbot"`
    ContentType  string `header:"Content-Type" validate:"required,oneof=application/json multipart/form-data"`
    XRequestID   string `header:"X-Request-ID" validate:"omitempty,uuid"`
}

// SQL 注入检查
func init() {
    validate.RegisterValidation("nosqlinject", func(fl validator.FieldLevel) bool {
        value := fl.Field().String()
        dangerous := []string{
            ";", "--", "/*", "*/", "xp_", "sp_", "exec", "union", "select", "insert", "update", "delete", "drop",
        }
        lower := strings.ToLower(value)
        for _, d := range dangerous {
            if strings.Contains(lower, d) {
                return false
            }
        }
        return true
    })
}
```
