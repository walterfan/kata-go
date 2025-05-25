# Quick start

* Initialize project
```shell
go mod init github.com/walterfan/kata-auth

go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v5
go get github.com/casbin/casbin/v2
```

* Run

```shell
go run main.go
```

* Test

```shell

#  get JWT token
export TOKEN=$(curl -s -X POST http://localhost:8080/token \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"pass"}' | jq -r '.token')
echo $TOKEN
#  get 200 response
curl -v -H "Authorization: Bearer $TOKEN" http://localhost:8080/admin
#  get 403 response
curl -v -H "Authorization: Bearer $TOKEN" http://localhost:8080/user
```

* explain

Casbin 是一个强大的开源访问控制库，支持多种编程语言（如 Go、Java、Python 等）。它提供了灵活的权限管理功能，可以实现常见的访问控制模型，例如：

1. `ACL (Access Control List)`：基于资源和用户的直接授权。
2. `RBAC (Role-Based Access Control)`：通过角色来分配权限。
3. `ABAC (Attribute-Based Access Control)`：基于属性的动态访问控制。

Casbin 的核心特点包括：
- **策略驱动**：权限规则存储在策略文件中（如 CSV 或 JSON），便于管理和修改。
- **模型抽象**：使用 `.conf` 文件定义访问控制模型，提高灵活性。
- **高性能**：优化的算法确保快速进行权限判断。
- **易于集成**：可与主流框架（如 Gin、Beego）无缝结合。

在 Go 项目中，Casbin 常用于构建细粒度的权限系统，比如 API 接口访问控制。

```go
auth := r.Group("/")
auth.Use(middleware.JWTAuth())
```

auth.Use(middleware.JWTAuth()) applies the JWTAuth middleware to all routes registered under the auth group.

Any incoming request to routes like /admin or /user will first go through the JWTAuth middleware, which ensures that the request contains a valid JWT token.


```go

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, "Bearer ")
        if len(parts) != 2 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
            c.Abort()
            return
        }
        /**
        This code parses a JWT token from the given string tokenStr.
        It uses a callback function to provide the secret key for validation.
        The function returns the parsed token and an error if any occurs during parsing or validation.
        */
        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
            c.Abort()
            return
        }

        c.Set("userID", claims["user_id"])
        c.Set("role", claims["role"])
        c.Next()
    }
}

```

## Casbin 配置

* model.conf

### 1. `request_definition`（请求定义）
```ini
[request_definition]
r = sub, obj, act
```
- 这部分定义了访问请求的结构。
- `r` 是一个请求，包含三个参数：
  - `sub`（主体）：通常是用户或角色。
  - `obj`（对象）：被访问的资源（如 `/admin`、`/user`）。
  - `act`（操作）：对资源执行的操作（如 [GET](file:///Users/walter.fan/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/routergroup.go#L37-L37)、[POST](file:///Users/walter.fan/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/routergroup.go#L38-L38)）。

---

### 2. `policy_definition`（策略定义）
```ini
[policy_definition]
p = sub, obj, act
```
- 这部分定义了策略（权限规则）的结构。
- `p` 是一条策略，也包含三个字段：
  - `sub`：有权执行操作的主体（用户或角色）。
  - `obj`：可访问的资源。
  - `act`：允许的操作。

---

### 3. `policy_effect`（策略效果）
```ini
[policy_effect]
e = some(where (p.eft == allow))
```
- 这部分定义了策略的效果。
- `e = some(where (p.eft == allow))` 表示只要有一条策略允许该请求（即 `p.eft == allow`），整个请求就视为允许。

---

### 4. `matchers`（匹配器）
```ini
[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
```
- 这部分定义了如何将请求与策略进行匹配。
- [m](file:///Users/walter.fan/go/pkg/mod/github.com/casbin/casbin/v2@v2.105.0/enforcer_synced.go#L29-L29) 是一个布尔表达式，表示只有当请求中的 `sub`、`obj`、`act` 都与某条策略完全匹配时，才认为该策略适用于当前请求。


这个配置文件实现了一个 **经典的 RBAC（基于角色的访问控制）模型**，其核心思想是：
- 每个请求由 `sub`（角色或用户）、`obj`（资源）、`act`（操作）组成。
- 系统会查找是否有对应的策略 `p` 匹配该请求。
- 如果存在匹配且策略允许（`allow`），则授权通过。

* [policy.csv](./config/policy.csv) 文件，可以定义具体的访问规则，例如：
```csv
p, admin, /admin, GET
p, user, /user, GET
```
表示：
- `admin` 角色可以访问 `/admin` 接口的 [GET](file:///Users/walter.fan/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/routergroup.go#L37-L37) 请求。
- [user](file:///Users/walter.fan/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/auth.go#L26-L26) 角色可以访问 `/user` 接口的 [GET](file:///Users/walter.fan/go/pkg/mod/github.com/gin-gonic/gin@v1.10.0/routergroup.go#L37-L37) 请求。