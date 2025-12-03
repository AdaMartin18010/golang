package specifications

import (
	"strings"
	"time"

	"github.com/yourusername/golang/internal/domain/user"
)

// ActiveUserSpecification 活跃用户规约
type ActiveUserSpecification struct{}

func (s ActiveUserSpecification) IsSatisfiedBy(u *user.User) bool {
	// 假设 User 有 Status 字段
	// 实际实现需要根据 User 实体的实际字段
	return true // 占位实现
}

// EmailSpecification 邮箱规约
type EmailSpecification struct {
	Email string
}

func (s EmailSpecification) IsSatisfiedBy(u *user.User) bool {
	return strings.EqualFold(u.Email, s.Email)
}

// CreatedAfterSpecification 创建时间规约
type CreatedAfterSpecification struct {
	After time.Time
}

func (s CreatedAfterSpecification) IsSatisfiedBy(u *user.User) bool {
	return u.CreatedAt.After(s.After)
}

// EmailDomainSpecification 邮箱域名规约
type EmailDomainSpecification struct {
	Domain string
}

func (s EmailDomainSpecification) IsSatisfiedBy(u *user.User) bool {
	parts := strings.Split(u.Email, "@")
	if len(parts) != 2 {
		return false
	}
	return strings.EqualFold(parts[1], s.Domain)
}

// 组合规约示例

// ActiveUserWithEmailDomain 活跃且邮箱域名匹配的用户
func ActiveUserWithEmailDomain(domain string) *AndSpec {
	return &AndSpec{
		left:  ActiveUserSpecification{},
		right: EmailDomainSpecification{Domain: domain},
	}
}

// AndSpec And 规约
type AndSpec struct {
	left  interface{ IsSatisfiedBy(*user.User) bool }
	right interface{ IsSatisfiedBy(*user.User) bool }
}

func (s *AndSpec) IsSatisfiedBy(u *user.User) bool {
	return s.left.IsSatisfiedBy(u) && s.right.IsSatisfiedBy(u)
}

