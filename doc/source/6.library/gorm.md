# GORM ORM 框架

```{contents} 目录
:depth: 3
```

## GORM 概述

GORM 是 Go 语言中最流行的 ORM 框架，功能全面且易于使用。

## 安装

```bash
go get gorm.io/gorm
go get gorm.io/driver/mysql    # MySQL
go get gorm.io/driver/postgres # PostgreSQL
go get gorm.io/driver/sqlite   # SQLite
```

## 连接数据库

```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func main() {
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    
    // 配置连接池
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
}
```

## 模型定义

```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:255;not null"`
    Email     string         `gorm:"uniqueIndex;size:255"`
    Age       int            `gorm:"default:18"`
    Active    bool           `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除
}

// 自定义表名
func (User) TableName() string {
    return "users"
}
```

## CRUD 操作

### 创建

```go
// 创建记录
user := User{Name: "Alice", Email: "alice@example.com"}
result := db.Create(&user)
// user.ID 会被自动赋值

// 批量创建
users := []User{{Name: "Bob"}, {Name: "Charlie"}}
db.Create(&users)

// 选择性创建
db.Select("Name", "Email").Create(&user)

// 忽略某些字段
db.Omit("Age").Create(&user)
```

### 查询

```go
// 查询单条记录
var user User
db.First(&user, 1)                    // 主键查询
db.First(&user, "name = ?", "Alice")  // 条件查询

// 查询所有
var users []User
db.Find(&users)

// 条件查询
db.Where("age > ?", 18).Find(&users)
db.Where(&User{Name: "Alice"}).Find(&users)
db.Where(map[string]interface{}{"name": "Alice"}).Find(&users)

// 链式条件
db.Where("name = ?", "Alice").
    Where("age > ?", 18).
    Find(&users)

// 选择字段
db.Select("name", "email").Find(&users)

// 排序和分页
db.Order("created_at desc").
    Limit(10).
    Offset(0).
    Find(&users)
```

### 更新

```go
// 更新单个字段
db.Model(&user).Update("name", "Bob")

// 更新多个字段
db.Model(&user).Updates(User{Name: "Bob", Age: 20})
db.Model(&user).Updates(map[string]interface{}{"name": "Bob", "age": 20})

// 批量更新
db.Model(&User{}).Where("age < ?", 18).Update("active", false)

// 使用 Select 指定更新字段
db.Model(&user).Select("Name").Updates(User{Name: "Bob", Age: 0})

// 使用 Omit 排除字段
db.Model(&user).Omit("Age").Updates(User{Name: "Bob", Age: 0})
```

### 删除

```go
// 软删除（需要有 DeletedAt 字段）
db.Delete(&user)

// 永久删除
db.Unscoped().Delete(&user)

// 条件删除
db.Where("name = ?", "Alice").Delete(&User{})

// 批量删除
db.Delete(&User{}, []int{1, 2, 3})
```

## 关联

### 一对一

```go
type User struct {
    ID      uint
    Name    string
    Profile Profile  // has one
}

type Profile struct {
    ID     uint
    UserID uint
    Bio    string
}

// 查询包含关联
db.Preload("Profile").Find(&users)
```

### 一对多

```go
type User struct {
    ID    uint
    Name  string
    Posts []Post  // has many
}

type Post struct {
    ID     uint
    UserID uint
    Title  string
}

// 预加载
db.Preload("Posts").Find(&users)

// 条件预加载
db.Preload("Posts", "published = ?", true).Find(&users)
```

### 多对多

```go
type User struct {
    ID    uint
    Name  string
    Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    ID    uint
    Name  string
    Users []User `gorm:"many2many:user_roles;"`
}

// 预加载
db.Preload("Roles").Find(&users)
```

## ⚠️ 常见陷阱

### 陷阱 1：零值更新问题

```go
// ❌ 零值不会更新
db.Model(&user).Updates(User{Name: "Bob", Age: 0})
// Age 不会被更新为 0

// ✅ 使用 map 或 Select
db.Model(&user).Updates(map[string]interface{}{"age": 0})
db.Model(&user).Select("Age").Updates(User{Age: 0})
```

### 陷阱 2：N+1 查询

```go
// ❌ N+1 问题
var users []User
db.Find(&users)
for _, user := range users {
    var posts []Post
    db.Where("user_id = ?", user.ID).Find(&posts)  // 每个用户都查一次
}

// ✅ 使用 Preload
db.Preload("Posts").Find(&users)
```

### 陷阱 3：goroutine 中共享 DB

```go
// ❌ DB 实例是安全的，但 Session 不是
db.Model(&user).Update("name", "Alice")  // 在多个 goroutine 中这样用是安全的

// ❌ 但是 Session 不安全
tx := db.Begin()
go func() {
    tx.Create(&user1)  // 不安全！
}()
go func() {
    tx.Create(&user2)  // 不安全！
}()

// ✅ 每个 goroutine 使用独立的 Session
go func() {
    tx := db.Session(&gorm.Session{})
    tx.Create(&user1)
}()
```

### 陷阱 4：忽略错误

```go
// ❌ 忽略错误
db.Create(&user)

// ✅ 检查错误
if err := db.Create(&user).Error; err != nil {
    return err
}

// 或者
result := db.Create(&user)
if result.Error != nil {
    return result.Error
}
fmt.Println("Rows affected:", result.RowsAffected)
```

## 事务

```go
// 自动事务
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err  // 返回错误会自动回滚
    }
    if err := tx.Create(&post).Error; err != nil {
        return err
    }
    return nil  // 返回 nil 会自动提交
})

// 手动事务
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&post).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

## 钩子

```go
type User struct {
    ID   uint
    Name string
    Hash string
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.Hash = generateHash(u.Name)
    return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    // 创建后的操作
    return nil
}

// 可用钩子：
// BeforeSave, AfterSave
// BeforeCreate, AfterCreate
// BeforeUpdate, AfterUpdate
// BeforeDelete, AfterDelete
// AfterFind
```

## 原生 SQL

```go
// 原生查询
var users []User
db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)

// 原生执行
db.Exec("UPDATE users SET age = ? WHERE name = ?", 20, "Alice")
```

## 参考资源

- [GORM 官方文档](https://gorm.io/docs/)
- [GORM GitHub](https://github.com/go-gorm/gorm)
