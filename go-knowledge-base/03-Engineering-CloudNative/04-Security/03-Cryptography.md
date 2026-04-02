# 加密库

> **分类**: 工程与云原生

---

## 标准库

```go
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
)
```

---

## 哈希

```go
h := sha256.New()
h.Write([]byte("hello"))
sum := h.Sum(nil)
```

---

## AES 加密

```go
key := make([]byte, 32)
rand.Read(key)

block, _ := aes.NewCipher(key)
```

---

## bcrypt

```go
import "golang.org/x/crypto/bcrypt"

hash, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
bcrypt.CompareHashAndPassword(hash, password)
```
