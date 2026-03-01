package specifications

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/golang/internal/domain/user"
)

// TestActiveUserSpecification 测试活跃用户规约
func TestActiveUserSpecification(t *testing.T) {
	spec := ActiveUserSpecification{}

	tests := []struct {
		name     string
		user     *user.User
		expected bool
	}{
		{
			name:     "active user - always returns true (placeholder)",
			user:     user.NewUser("test@example.com", "Test User"),
			expected: true,
		},
		{
			name: "any user - placeholder implementation",
			user: &user.User{
				ID:    "test-id",
				Email: "inactive@example.com",
				Name:  "Inactive User",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := spec.IsSatisfiedBy(tt.user)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestEmailSpecification 测试邮箱规约
func TestEmailSpecification(t *testing.T) {
	tests := []struct {
		name     string
		spec     EmailSpecification
		user     *user.User
		expected bool
	}{
		{
			name:     "email matches exactly",
			spec:     EmailSpecification{Email: "test@example.com"},
			user:     &user.User{Email: "test@example.com"},
			expected: true,
		},
		{
			name:     "email does not match",
			spec:     EmailSpecification{Email: "test@example.com"},
			user:     &user.User{Email: "other@example.com"},
			expected: false,
		},
		{
			name:     "case insensitive match - lowercase",
			spec:     EmailSpecification{Email: "TEST@EXAMPLE.COM"},
			user:     &user.User{Email: "test@example.com"},
			expected: true,
		},
		{
			name:     "case insensitive match - mixed case",
			spec:     EmailSpecification{Email: "Test@Example.com"},
			user:     &user.User{Email: "test@example.com"},
			expected: true,
		},
		{
			name:     "case insensitive match - user uppercase",
			spec:     EmailSpecification{Email: "test@example.com"},
			user:     &user.User{Email: "TEST@EXAMPLE.COM"},
			expected: true,
		},
		{
			name:     "empty email in spec",
			spec:     EmailSpecification{Email: ""},
			user:     &user.User{Email: ""},
			expected: true,
		},
		{
			name:     "empty user email with non-empty spec",
			spec:     EmailSpecification{Email: "test@example.com"},
			user:     &user.User{Email: ""},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.spec.IsSatisfiedBy(tt.user)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCreatedAfterSpecification 测试创建时间规约
func TestCreatedAfterSpecification(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		spec     CreatedAfterSpecification
		user     *user.User
		expected bool
	}{
		{
			name:     "user created after specified time",
			spec:     CreatedAfterSpecification{After: baseTime},
			user:     &user.User{CreatedAt: baseTime.Add(1 * time.Hour)},
			expected: true,
		},
		{
			name:     "user created exactly at specified time",
			spec:     CreatedAfterSpecification{After: baseTime},
			user:     &user.User{CreatedAt: baseTime},
			expected: false,
		},
		{
			name:     "user created before specified time",
			spec:     CreatedAfterSpecification{After: baseTime},
			user:     &user.User{CreatedAt: baseTime.Add(-1 * time.Hour)},
			expected: false,
		},
		{
			name:     "user created 1 second after",
			spec:     CreatedAfterSpecification{After: baseTime},
			user:     &user.User{CreatedAt: baseTime.Add(1 * time.Second)},
			expected: true,
		},
		{
			name:     "user created 1 second before",
			spec:     CreatedAfterSpecification{After: baseTime},
			user:     &user.User{CreatedAt: baseTime.Add(-1 * time.Second)},
			expected: false,
		},
		{
			name:     "very old user with recent cutoff",
			spec:     CreatedAfterSpecification{After: baseTime},
			user:     &user.User{CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			expected: false,
		},
		{
			name:     "recent user with old cutoff",
			spec:     CreatedAfterSpecification{After: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			user:     &user.User{CreatedAt: baseTime},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.spec.IsSatisfiedBy(tt.user)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestEmailDomainSpecification 测试邮箱域名规约
func TestEmailDomainSpecification(t *testing.T) {
	tests := []struct {
		name     string
		spec     EmailDomainSpecification
		user     *user.User
		expected bool
	}{
		{
			name:     "domain matches exactly",
			spec:     EmailDomainSpecification{Domain: "example.com"},
			user:     &user.User{Email: "user@example.com"},
			expected: true,
		},
		{
			name:     "domain does not match",
			spec:     EmailDomainSpecification{Domain: "example.com"},
			user:     &user.User{Email: "user@other.com"},
			expected: false,
		},
		{
			name:     "case insensitive - spec uppercase",
			spec:     EmailDomainSpecification{Domain: "EXAMPLE.COM"},
			user:     &user.User{Email: "user@example.com"},
			expected: true,
		},
		{
			name:     "case insensitive - user uppercase",
			spec:     EmailDomainSpecification{Domain: "example.com"},
			user:     &user.User{Email: "user@EXAMPLE.COM"},
			expected: true,
		},
		{
			name:     "case insensitive - both mixed",
			spec:     EmailDomainSpecification{Domain: "Example.Com"},
			user:     &user.User{Email: "user@example.com"},
			expected: true,
		},
		{
			name:     "subdomain does not match parent",
			spec:     EmailDomainSpecification{Domain: "example.com"},
			user:     &user.User{Email: "user@mail.example.com"},
			expected: false,
		},
		{
			name:     "subdomain matches exactly",
			spec:     EmailDomainSpecification{Domain: "mail.example.com"},
			user:     &user.User{Email: "user@mail.example.com"},
			expected: true,
		},
		{
			name:     "invalid email - no at symbol",
			spec:     EmailDomainSpecification{Domain: "example.com"},
			user:     &user.User{Email: "invalid-email"},
			expected: false,
		},
		{
			name:     "invalid email - multiple at symbols",
			spec:     EmailDomainSpecification{Domain: "example.com"},
			user:     &user.User{Email: "user@@example.com"},
			expected: false,
		},
		{
			name:     "invalid email - empty string",
			spec:     EmailDomainSpecification{Domain: "example.com"},
			user:     &user.User{Email: ""},
			expected: false,
		},
		{
			name:     "email with plus sign",
			spec:     EmailDomainSpecification{Domain: "example.com"},
			user:     &user.User{Email: "user+tag@example.com"},
			expected: true,
		},
		{
			name:     "empty domain in spec",
			spec:     EmailDomainSpecification{Domain: ""},
			user:     &user.User{Email: "user@example.com"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.spec.IsSatisfiedBy(tt.user)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAndSpec 测试 And 规约组合
func TestAndSpec(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		left     interface{ IsSatisfiedBy(*user.User) bool }
		right    interface{ IsSatisfiedBy(*user.User) bool }
		user     *user.User
		expected bool
	}{
		{
			name:     "both specs satisfied",
			left:     EmailDomainSpecification{Domain: "example.com"},
			right:    CreatedAfterSpecification{After: baseTime},
			user:     &user.User{Email: "user@example.com", CreatedAt: baseTime.Add(1 * time.Hour)},
			expected: true,
		},
		{
			name:     "left spec not satisfied",
			left:     EmailDomainSpecification{Domain: "example.com"},
			right:    CreatedAfterSpecification{After: baseTime},
			user:     &user.User{Email: "user@other.com", CreatedAt: baseTime.Add(1 * time.Hour)},
			expected: false,
		},
		{
			name:     "right spec not satisfied",
			left:     EmailDomainSpecification{Domain: "example.com"},
			right:    CreatedAfterSpecification{After: baseTime},
			user:     &user.User{Email: "user@example.com", CreatedAt: baseTime.Add(-1 * time.Hour)},
			expected: false,
		},
		{
			name:     "neither spec satisfied",
			left:     EmailDomainSpecification{Domain: "example.com"},
			right:    CreatedAfterSpecification{After: baseTime},
			user:     &user.User{Email: "user@other.com", CreatedAt: baseTime.Add(-1 * time.Hour)},
			expected: false,
		},
		{
			name:     "email specs combined",
			left:     EmailDomainSpecification{Domain: "example.com"},
			right:    EmailSpecification{Email: "specific@example.com"},
			user:     &user.User{Email: "specific@example.com"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			andSpec := &AndSpec{left: tt.left, right: tt.right}
			result := andSpec.IsSatisfiedBy(tt.user)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestActiveUserWithEmailDomain 测试组合规约函数
func TestActiveUserWithEmailDomain(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		user     *user.User
		expected bool
	}{
		{
			name:     "active user with matching domain",
			domain:   "example.com",
			user:     &user.User{Email: "user@example.com"},
			expected: true,
		},
		{
			name:     "active user with non-matching domain",
			domain:   "example.com",
			user:     &user.User{Email: "user@other.com"},
			expected: false,
		},
		{
			name:     "case insensitive domain match",
			domain:   "EXAMPLE.COM",
			user:     &user.User{Email: "user@example.com"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := ActiveUserWithEmailDomain(tt.domain)
			result := spec.IsSatisfiedBy(tt.user)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSpecification_ComplexComposition 测试复杂规约组合场景
func TestSpecification_ComplexComposition(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// 创建测试用户
	activeRecentUser := &user.User{
		Email:     "active@example.com",
		CreatedAt: baseTime.Add(24 * time.Hour),
	}
	activeOldUser := &user.User{
		Email:     "active@example.com",
		CreatedAt: baseTime.Add(-24 * time.Hour),
	}
	otherDomainUser := &user.User{
		Email:     "user@other.com",
		CreatedAt: baseTime.Add(24 * time.Hour),
	}

	tests := []struct {
		name     string
		spec     *AndSpec
		user     *user.User
		expected bool
	}{
		{
			name:     "active + example.com + recent",
			spec:     &AndSpec{left: ActiveUserSpecification{}, right: EmailDomainSpecification{Domain: "example.com"}},
			user:     activeRecentUser,
			expected: true,
		},
		{
			name:     "active + other.com domain",
			spec:     &AndSpec{left: ActiveUserSpecification{}, right: EmailDomainSpecification{Domain: "example.com"}},
			user:     otherDomainUser,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.spec.IsSatisfiedBy(tt.user)
			assert.Equal(t, tt.expected, result)
		})
	}

	// 测试嵌套的 AndSpec
	nestedSpec := &AndSpec{
		left: &AndSpec{
			left:  ActiveUserSpecification{},
			right: EmailDomainSpecification{Domain: "example.com"},
		},
		right: CreatedAfterSpecification{After: baseTime},
	}

	assert.True(t, nestedSpec.IsSatisfiedBy(activeRecentUser), "Nested spec should match active recent user with correct domain")
	assert.False(t, nestedSpec.IsSatisfiedBy(activeOldUser), "Nested spec should not match old user")
	assert.False(t, nestedSpec.IsSatisfiedBy(otherDomainUser), "Nested spec should not match user with different domain")
}

// TestSpecifications_WithRealUsers 使用真实 User 实体测试规约
func TestSpecifications_WithRealUsers(t *testing.T) {
	// 创建真实用户
	user1 := user.NewUser("user1@example.com", "User One")
	user2 := user.NewUser("user2@other.com", "User Two")
	user3 := user.NewUser("user3@example.com", "User Three")

	// 等待一点时间，然后创建一个新时间 cutoff
	time.Sleep(10 * time.Millisecond)
	cutoff := time.Now()
	time.Sleep(10 * time.Millisecond)

	user4 := user.NewUser("user4@example.com", "User Four")

	// 测试邮箱域名规约
	domainSpec := EmailDomainSpecification{Domain: "example.com"}
	assert.True(t, domainSpec.IsSatisfiedBy(user1))
	assert.False(t, domainSpec.IsSatisfiedBy(user2))
	assert.True(t, domainSpec.IsSatisfiedBy(user3))
	assert.True(t, domainSpec.IsSatisfiedBy(user4))

	// 测试创建时间规约
	timeSpec := CreatedAfterSpecification{After: cutoff}
	assert.False(t, timeSpec.IsSatisfiedBy(user1))
	assert.False(t, timeSpec.IsSatisfiedBy(user2))
	assert.False(t, timeSpec.IsSatisfiedBy(user3))
	assert.True(t, timeSpec.IsSatisfiedBy(user4))

	// 测试组合规约
	combined := &AndSpec{
		left:  EmailDomainSpecification{Domain: "example.com"},
		right: CreatedAfterSpecification{After: cutoff},
	}
	assert.False(t, combined.IsSatisfiedBy(user1), "User1: wrong time")
	assert.False(t, combined.IsSatisfiedBy(user2), "User2: wrong domain and time")
	assert.False(t, combined.IsSatisfiedBy(user3), "User3: wrong time")
	assert.True(t, combined.IsSatisfiedBy(user4), "User4: matches domain and time")
}

// BenchmarkEmailSpecification 邮箱规约性能测试
func BenchmarkEmailSpecification(b *testing.B) {
	spec := EmailSpecification{Email: "test@example.com"}
	u := &user.User{Email: "test@example.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spec.IsSatisfiedBy(u)
	}
}

// BenchmarkEmailDomainSpecification 邮箱域名规约性能测试
func BenchmarkEmailDomainSpecification(b *testing.B) {
	spec := EmailDomainSpecification{Domain: "example.com"}
	u := &user.User{Email: "user@example.com"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spec.IsSatisfiedBy(u)
	}
}

// BenchmarkAndSpec And 规约性能测试
func BenchmarkAndSpec(b *testing.B) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	andSpec := &AndSpec{
		left:  EmailDomainSpecification{Domain: "example.com"},
		right: CreatedAfterSpecification{After: baseTime},
	}
	u := &user.User{
		Email:     "user@example.com",
		CreatedAt: baseTime.Add(1 * time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		andSpec.IsSatisfiedBy(u)
	}
}
