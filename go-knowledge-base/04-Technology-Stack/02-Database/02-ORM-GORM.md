# TS-DB-002: GORM - Go ORM Architecture and Patterns

> **维度**: Technology Stack > Database
> **级别**: S (16+ KB)
> **标签**: #gorm #orm #golang #database #sql
> **权威来源**:
>
> - [GORM Documentation](https://gorm.io/) - Official docs
> - [GORM Source Code](https://github.com/go-gorm/gorm) - GitHub
> - [GORM Migrations](https://gorm.io/docs/migration.html) - Schema migrations

---

## 1. GORM Architecture Overview

### 1.1 Core Components

```
┌─────────────────────────────────────────────────────────────────┐
│                       GORM Architecture                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                     Application Layer                      │  │
│  │  db.Create(&user)  db.First(&user)  db.Model(&user).Update│  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │                      Session Layer                         │  │
│  │  - Method Chain Builder                                    │  │
│  │  - Scope Functions                                         │  │
│  │  - Hook Execution                                          │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │                     Statement Layer                        │  │
│  │  - SQL Generation                                          │  │
│  │  - Clause Building                                         │  │
│  │  - Query Building                                          │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │                     Callbacks Layer                        │  │
│  │  ┌──────────┬──────────┬──────────┬──────────┬──────────┐  │  │
│  │  │ Create   │ Query    │ Update   │ Delete   │ Row/Raw  │  │  │
│  │  │──────────│──────────│──────────│──────────│──────────│  │  │
│  │  │Before    │Before    │Before    │Before    │Before    │  │  │
│  │  │Create    │Query     │Update    │Delete    │Execute   │  │  │
│  │  │          │          │          │          │          │  │  │
│  │  │After     │After     │After     │After     │After     │  │  │
│  │  │Create    │Query     │Update    │Delete    │Execute   │  │  │
│  │  └──────────┴──────────┴──────────┴──────────┴──────────┘  │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │                   Dialect Layer                            │  │
│  │  - MySQL, PostgreSQL, SQLite, SQL Server, ClickHouse      │  │
│  │  - Driver-specific SQL generation                          │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │                   Connection Pool                          │  │
│  │  - database/sql pool integration                           │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Session-Based Design

GORM uses a session-based fluent API where each call returns a new session:

```go
// Each method returns a *gorm.DB (new session)
// Sessions are immutable - chain operations safely

db.Where("age > ?", 18).           // Returns new session with WHERE clause
  Where("status = ?", "active").   // Returns new session with AND clause
  Order("created_at desc").        // Returns new session with ORDER BY
  Find(&users)                     // Executes query

// Original db is unchanged
db.Find(&allUsers) // No WHERE clause
```

---

## 2. Model Definition

### 2.1 Struct Tags and Conventions

```go
type User struct {
    // ID is primary key by convention
    // GORM uses uint/int with auto-increment by default
    ID uint `gorm:"primaryKey"`

    // UUID primary key
    ID uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

    // Timestamps - auto-populated
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete

    // Column name override
    Name string `gorm:"column:user_name"`

    // Size constraint
    Email string `gorm:"size:255;uniqueIndex"`

    // NOT NULL constraint
    Password string `gorm:"not null"`

    // Default value
    Role string `gorm:"default:'user'"`

    // Index
    Age int `gorm:"index:idx_age_status"`
    Status string `gorm:"index:idx_age_status"`

    // Composite unique index
    FirstName string `gorm:"uniqueIndex:idx_name"`
    LastName  string `gorm:"uniqueIndex:idx_name"`

    // Ignore field
    TempData string `gorm:"-"`

    // Embedded struct
    Address

    // Foreign key
    CompanyID uint
    Company   Company `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

    // Many-to-many
    Languages []Language `gorm:"many2many:user_languages;"`

    // Serialization
    Settings datatypes.JSON `gorm:"type:jsonb"`
    Metadata datatypes.JSON `gorm:"serializer:json"`
}

type Address struct {
    Street string
    City   string
    Country string
}

// Table name override
func (User) TableName() string {
    return "sys_users"
}

// Dynamic table name
func (u User) TableName() string {
    if u.Role == "admin" {
        return "admin_users"
    }
    return "users"
}
```

### 2.2 Advanced Field Types

```go
type AdvancedUser struct {
    ID uint

    // JSON field
    Config datatypes.JSON

    // Array (PostgreSQL)
    Tags pq.StringArray `gorm:"type:text[]"`

    // Enum
    Status string `gorm:"type:enum('active','inactive','banned')"`

    // Custom type
    EncryptedData EncryptedString

    // Scanner/Valuer interface
    CustomField CustomType
}

// Custom type implementing Scanner/Valuer
type EncryptedString struct {
    Plaintext string
}

func (e EncryptedString) Value() (driver.Value, error) {
    if e.Plaintext == "" {
        return nil, nil
    }
    return encrypt(e.Plaintext), nil
}

func (e *EncryptedString) Scan(value interface{}) error {
    if value == nil {
        return nil
    }
    bytes, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    e.Plaintext = decrypt(bytes)
    return nil
}
```

---

## 3. CRUD Operations

### 3.1 Create

```go
// Single create
user := User{Name: "John", Email: "john@example.com"}
result := db.Create(&user)
// user.ID now populated
// result.Error, result.RowsAffected

// Batch create
users := []User{
    {Name: "John", Email: "john@example.com"},
    {Name: "Jane", Email: "jane@example.com"},
}
db.Create(&users)

// Create with selected fields
db.Select("Name", "Email").Create(&user)

// Create or update (upsert)
db.Clauses(clause.OnConflict{
    Columns:   []clause.Column{{Name: "email"}},
    DoUpdates: clause.AssignmentColumns([]string{"name", "updated_at"}),
}).Create(&user)

// Create in batches (memory efficient for large datasets)
db.CreateInBatches(users, 100)
```

### 3.2 Read Operations

```go
// First record ordered by primary key
var user User
db.First(&user) // SELECT * FROM users ORDER BY id LIMIT 1

// Last record
db.Last(&user)

// Find by primary key
db.First(&user, 10)              // SELECT * FROM users WHERE id = 10
db.First(&user, "10")            // Same
db.Find(&users, []int{1, 2, 3})  // SELECT * WHERE id IN (1,2,3)

// Find with conditions
db.Where("name = ?", "john").First(&user)
db.Where(&User{Name: "john", Age: 20}).First(&user)
db.Where(map[string]interface{}{"name": "john", "age": 20}).Find(&users)

// Inline condition (shorthand)
db.First(&user, "name = ?", "john")
db.Find(&users, "age > ? AND status = ?", 18, "active")

// Struct condition (only non-zero fields)
db.Where(&User{Name: "john", Age: 0}).Find(&users) // Age ignored (zero value)
db.Where(&User{Name: "john"}).Find(&users)

// Map condition (zero values included)
db.Where(map[string]interface{}{"name": "john", "age": 0}).Find(&users)

// NOT conditions
db.Not("name = ?", "john").Find(&users)
db.Not(&User{Name: "john"}).Find(&users)

// Or conditions
db.Where("role = ?", "admin").Or("role = ?", "moderator").Find(&users)

// Order, Limit, Offset
db.Order("age desc, name").Limit(10).Offset(20).Find(&users)

// Distinct
db.Distinct("name", "age").Find(&users)

// Select specific fields
db.Select("name", "email").Find(&users)
db.Select("coalesce(age,?)", 18).Find(&users)

// Scan to struct
var result struct {
    Name  string
    Total int64
}
db.Model(&User{}).Select("name, sum(age) as total").Group("name").Scan(&result)

// Pluck single column
var ages []int64
db.Model(&User{}).Pluck("age", &ages)

// Rows
type Result struct {
    Name string
    Age  int
}
rows, err := db.Model(&User{}).Select("name, age").Rows()
defer rows.Close()
for rows.Next() {
    var r Result
    db.ScanRows(rows, &r)
    // process r
}
```

### 3.3 Update Operations

```go
// Update single field
db.Model(&user).Update("name", "new name")

// Update multiple fields
db.Model(&user).Updates(User{Name: "new", Age: 18})        // Non-zero fields only
db.Model(&user).Updates(map[string]interface{}{"name": "new", "age": 0}) // All fields

// Update with conditions (batch update)
db.Model(&User{}).Where("active = ?", true).Update("age", gorm.Expr("age + ?", 1))

// Update without callbacks
db.Model(&user).UpdateColumn("name", "new")
db.Model(&user).UpdateColumns(User{Name: "new", Age: 18})

// Expression update
db.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
db.Model(&product).Updates(map[string]interface{}{
    "price": gorm.Expr("price * ? + ?", 2, 100),
    "quantity": gorm.Expr("quantity - ?", 1),
})

// Check if field has changed
if db.Model(&user).Update("name", "new").RowsAffected == 0 {
    // No rows updated
}
```

### 3.4 Delete Operations

```go
// Soft delete (requires DeletedAt field)
db.Delete(&user) // Sets DeletedAt to current time

// Find including soft deleted
var users []User
db.Unscoped().Where("age = 20").Find(&users)

// Hard delete
db.Unscoped().Delete(&user)

// Batch delete
db.Where("email LIKE ?", "%@test.com%").Delete(&User{})

// Delete with primary key
db.Delete(&User{}, 10)
db.Delete(&User{}, []int{1, 2, 3})
```

---

## 4. Associations

### 4.1 Relationship Types

```go
type User struct {
    gorm.Model

    // Has One
    CreditCard CreditCard

    // Has Many
    Emails []Email

    // Belongs To
    CompanyID uint
    Company   Company

    // Many to Many
    Languages []Language `gorm:"many2many:user_languages;"`
}

type CreditCard struct {
    gorm.Model
    Number string
    UserID uint // Foreign key
}

type Email struct {
    gorm.Model
    Email  string
    UserID uint // Foreign key
}

type Company struct {
    gorm.Model
    Name string
}

type Language struct {
    gorm.Model
    Name string
}
```

### 4.2 Preloading

```go
// Eager loading
db.Preload("CreditCard").Preload("Emails").Find(&users)

// Nested preloading
db.Preload("Orders.Items").Find(&users)

// Preload with conditions
db.Preload("Emails", "active = ?", true).Find(&users)
db.Preload("Emails", func(db *gorm.DB) *gorm.DB {
    return db.Order("email DESC")
}).Find(&users)

// Joins (eager loading with JOIN)
db.Joins("Company").Find(&users)
db.Joins("Company", db.Where(&Company{Active: true})).Find(&users)
```

### 4.3 Association Operations

```go
// Append associations
db.Model(&user).Association("Languages").Append(&Language{Name: "Go"})

// Replace associations
db.Model(&user).Association("Languages").Replace(&lang1, &lang2)

// Delete associations
db.Model(&user).Association("Languages").Delete(&lang1)

// Clear associations
db.Model(&user).Association("Languages").Clear()

// Count associations
count := db.Model(&user).Association("Languages").Count()

// Find associations
var languages []Language
db.Model(&user).Association("Languages").Find(&languages)
```

---

## 5. Hooks and Callbacks

### 5.1 Available Hooks

```go
// Creating hooks
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.UUID = uuid.New()
    return
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
    // Send welcome email
    return
}

// Updating hooks
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
    // Validation
    if u.Age < 0 {
        return errors.New("age cannot be negative")
    }
    return
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
    // Clear cache
    return
}

// Deleting hooks
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
    // Check if user has active orders
    if u.HasActiveOrders {
        return errors.New("cannot delete user with active orders")
    }
    return
}

// Querying hooks
func (u *User) AfterFind(tx *gorm.DB) (err error) {
    // Decrypt sensitive data
    return
}
```

### 5.2 Registering Custom Callbacks

```go
// Register global callback
db.Callback().Create().Before("gorm:create").Register("my_callback", func(db *gorm.DB) {
    // Custom logic
})

// Remove callback
db.Callback().Create().Remove("gorm:create")

// Replace callback
db.Callback().Create().Replace("gorm:create", func(db *gorm.DB) {
    // New implementation
})
```

---

## 6. Transactions

### 6.1 Manual Transactions

```go
// Begin transaction
tx := db.Begin()

// Do operations
tx.Create(&user)
tx.Create(&order)

// Commit or rollback
if err := tx.Error; err != nil {
    tx.Rollback()
    return err
}
return tx.Commit().Error
```

### 6.2 Automatic Transactions

```go
// Transaction block
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err // Rollback
    }
    if err := tx.Create(&order).Error; err != nil {
        return err // Rollback
    }
    return nil // Commit
})

// Nested transactions (savepoints)
err := db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user)

    tx.Transaction(func(tx2 *gorm.DB) error {
        tx2.Create(&address) // Savepoint
        return errors.New("rollback this") // Rollback to savepoint
    })

    tx.Create(&order)
    return nil
})
```

### 6.3 Transaction Options

```go
import "gorm.io/gorm/clause"

// Read-only transaction
err := db.Clauses(db.Read()).Transaction(func(tx *gorm.DB) error {
    tx.First(&user, 1)
    return nil
})

// Specify isolation level
err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Transaction(func(tx *gorm.DB) error {
    tx.First(&user, 1)
    return nil
})
```

---

## 7. Performance Optimization

### 7.1 Connection Pooling

```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// Get underlying sql.DB for pool configuration
sqlDB, err := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 7.2 Prepared Statements

```go
// Enabled by default
// To disable:
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    PrepareStmt: false,
})
```

### 7.3 Caching

```go
// Enable query cache (database level)
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    QueryFields: true,
})

// Application-level cache
var user User
cacheKey := fmt.Sprintf("user:%d", id)
if err := cache.Get(cacheKey, &user); err == nil {
    return user, nil
}

if err := db.First(&user, id).Error; err != nil {
    return User{}, err
}

cache.Set(cacheKey, user, time.Hour)
return user, nil
```

### 7.4 Batch Operations

```go
// Efficient batch insert
db.CreateInBatches(users, 100)

// Batch update with SQL expression
db.Model(&User{}).Where("active = ?", true).
    UpdateColumn("login_count", gorm.Expr("login_count + ?", 1))
```

---

## 8. Migration

### 8.1 Auto Migration

```go
// Migrate schema
db.AutoMigrate(&User{})
db.AutoMigrate(&User{}, &Product{}, &Order{})

// Migrate specific tables
db.Migrator().CreateTable(&User{})
db.Migrator().AddColumn(&User{}, "Age")
db.Migrator().DropColumn(&User{}, "Age")
db.Migrator().CreateIndex(&User{}, "idx_name")
db.Migrator().DropIndex(&User{}, "idx_name")

// Check if table exists
hasTable := db.Migrator().HasTable(&User{})
```

### 8.2 Versioned Migrations

```go
import "gorm.io/gormigrate"

m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
    {
        ID: "20230101",
        Migrate: func(tx *gorm.DB) error {
            return tx.AutoMigrate(&User{})
        },
        Rollback: func(tx *gorm.DB) error {
            return tx.Migrator().DropTable("users")
        },
    },
})

if err := m.Migrate(); err != nil {
    log.Fatalf("Migration failed: %v", err)
}
```

---

## 9. Checklist

```
GORM Best Practices:
□ Define proper struct tags
□ Use appropriate field types
□ Implement hooks for business logic
□ Use transactions for data consistency
□ Preload associations to avoid N+1
□ Configure connection pooling
□ Use batch operations for large datasets
□ Add database indexes
□ Enable query logging in development
□ Handle errors properly
```
