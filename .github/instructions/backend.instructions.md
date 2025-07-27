# Backend Development Instructions - Go Language

## プロジェクト概要
Web TodoアプリのバックエンドAPIをGo言語で構築します。RESTful APIを提供し、認証、Todo CRUD操作を実装します。

## 技術スタック
- **Language**: Go 1.21+
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **ORM**: GORM (gorm.io/gorm)
- **Authentication**: JWT (github.com/golang-jwt/jwt/v5)
- **Migration**: golang-migrate
- **Validation**: github.com/go-playground/validator/v10
- **Configuration**: Viper (github.com/spf13/viper)
- **Logging**: logrus (github.com/sirupsen/logrus)
- **Testing**: testify (github.com/stretchr/testify)
- **Documentation**: Swagger (github.com/swaggo/gin-swagger)

## プロジェクト構造
```
backend/
├── cmd/
│   └── server/
│       └── main.go               # アプリケーションエントリーポイント
├── internal/                     # プライベートアプリケーションコード
│   ├── config/                   # 設定管理
│   │   └── config.go
│   ├── handler/                  # HTTPハンドラー
│   │   ├── auth.go
│   │   ├── todo.go
│   │   └── user.go
│   ├── middleware/               # HTTPミドルウェア
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── model/                    # データモデル
│   │   ├── todo.go
│   │   └── user.go
│   ├── repository/               # データアクセス層
│   │   ├── todo.go
│   │   └── user.go
│   ├── service/                  # ビジネスロジック層
│   │   ├── auth.go
│   │   ├── todo.go
│   │   └── user.go
│   └── utils/                    # ユーティリティ関数
│       ├── hash.go
│       ├── jwt.go
│       └── validator.go
├── pkg/                          # 外部パッケージ
│   └── database/
│       └── postgres.go
├── migrations/                   # データベースマイグレーション
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_todos_table.up.sql
│   └── 002_create_todos_table.down.sql
├── docs/                         # Swagger生成ドキュメント
├── docker/
│   └── Dockerfile
├── scripts/                      # ビルド・デプロイスクリプト
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## アーキテクチャパターン
Clean Architectureを採用し、依存関係を明確に分離します。

### レイヤー構成
1. **Handler層**: HTTPリクエスト/レスポンス処理
2. **Service層**: ビジネスロジック
3. **Repository層**: データアクセス
4. **Model層**: ドメインモデル

### 依存関係注入
```go
type Dependencies struct {
    UserService auth.UserService
    TodoService todo.TodoService
    DB          *gorm.DB
    Config      *config.Config
}

func NewDependencies(cfg *config.Config) (*Dependencies, error) {
    db, err := database.NewPostgresConnection(cfg.Database)
    if err != nil {
        return nil, err
    }

    userRepo := repository.NewUserRepository(db)
    todoRepo := repository.NewTodoRepository(db)

    userService := service.NewUserService(userRepo)
    todoService := service.NewTodoService(todoRepo)

    return &Dependencies{
        UserService: userService,
        TodoService: todoService,
        DB:          db,
        Config:      cfg,
    }, nil
}
```

## データモデル

### User Model
```go
// internal/model/user.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        string         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
    Name      string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
    Password  string         `json:"-" gorm:"not null" validate:"required,min=8"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    Todos     []Todo         `json:"todos,omitempty" gorm:"foreignKey:UserID"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
    User  *User  `json:"user"`
    Token string `json:"token"`
}
```

### Todo Model
```go
// internal/model/todo.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type Priority string

const (
    PriorityLow    Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh   Priority = "high"
)

type Todo struct {
    ID          string         `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Title       string         `json:"title" gorm:"not null" validate:"required,min=1,max=200"`
    Description *string        `json:"description" gorm:"type:text"`
    Completed   bool           `json:"completed" gorm:"default:false"`
    Priority    Priority       `json:"priority" gorm:"type:varchar(10);default:'medium'" validate:"oneof=low medium high"`
    DueDate     *time.Time     `json:"due_date"`
    UserID      string         `json:"user_id" gorm:"type:uuid;not null"`
    User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateTodoRequest struct {
    Title       string     `json:"title" validate:"required,min=1,max=200"`
    Description *string    `json:"description"`
    Priority    Priority   `json:"priority" validate:"oneof=low medium high"`
    DueDate     *time.Time `json:"due_date"`
}

type UpdateTodoRequest struct {
    Title       *string    `json:"title" validate:"omitempty,min=1,max=200"`
    Description *string    `json:"description"`
    Completed   *bool      `json:"completed"`
    Priority    *Priority  `json:"priority" validate:"omitempty,oneof=low medium high"`
    DueDate     *time.Time `json:"due_date"`
}
```

## API設計

### 認証エンドポイント
```go
// POST /api/auth/register - ユーザー登録
// POST /api/auth/login    - ログイン
// POST /api/auth/refresh  - トークンリフレッシュ

// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Register request"
// @Success 201 {object} model.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    // 実装
}
```

### Todo エンドポイント
```go
// GET    /api/todos      - Todo一覧取得
// POST   /api/todos      - Todo作成
// GET    /api/todos/:id  - Todo詳細取得
// PUT    /api/todos/:id  - Todo更新
// DELETE /api/todos/:id  - Todo削除

// @Summary Get todos
// @Tags todos
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param completed query bool false "Filter by completion status"
// @Success 200 {object} PaginatedResponse{data=[]model.Todo}
// @Failure 401 {object} ErrorResponse
// @Security BearerAuth
// @Router /todos [get]
func (h *TodoHandler) GetTodos(c *gin.Context) {
    // 実装
}
```

## エラーハンドリング

### 標準エラーレスポンス
```go
type ErrorResponse struct {
    Error   string                 `json:"error"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func HandleError(c *gin.Context, err error, statusCode int) {
    var response ErrorResponse
    
    switch e := err.(type) {
    case validator.ValidationErrors:
        response.Error = "validation_error"
        response.Message = "入力データに誤りがあります"
        response.Details = formatValidationErrors(e)
        statusCode = http.StatusBadRequest
    default:
        response.Error = "internal_error"
        response.Message = err.Error()
    }
    
    c.JSON(statusCode, response)
}
```

## セキュリティ

### JWT認証
```go
// internal/utils/jwt.go
func GenerateToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
        "iat":     time.Now().Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
}
```

### ミドルウェア
```go
// internal/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := utils.ValidateToken(tokenString)
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        c.Set("user_id", claims["user_id"])
        c.Next()
    }
}
```

## 設定管理
```go
// internal/config/config.go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
    Port string `mapstructure:"port" env:"PORT" envDefault:"8080"`
    Host string `mapstructure:"host" env:"HOST" envDefault:"localhost"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host" env:"DB_HOST" envDefault:"localhost"`
    Port     string `mapstructure:"port" env:"DB_PORT" envDefault:"5432"`
    User     string `mapstructure:"user" env:"DB_USER" envDefault:"postgres"`
    Password string `mapstructure:"password" env:"DB_PASSWORD"`
    Name     string `mapstructure:"name" env:"DB_NAME" envDefault:"todoapp"`
    SSLMode  string `mapstructure:"ssl_mode" env:"DB_SSL_MODE" envDefault:"disable"`
}

type JWTConfig struct {
    Secret     string        `mapstructure:"secret" env:"JWT_SECRET"`
    Expiration time.Duration `mapstructure:"expiration" env:"JWT_EXPIRATION" envDefault:"24h"`
}
```

## テスト戦略

### ユニットテスト
```go
// internal/service/todo_test.go
func TestTodoService_CreateTodo(t *testing.T) {
    mockRepo := &mocks.TodoRepository{}
    service := NewTodoService(mockRepo)

    req := &model.CreateTodoRequest{
        Title:    "Test Todo",
        Priority: model.PriorityMedium,
    }

    expectedTodo := &model.Todo{
        ID:       "test-id",
        Title:    req.Title,
        Priority: req.Priority,
        UserID:   "user-id",
    }

    mockRepo.On("Create", mock.AnythingOfType("*model.Todo")).Return(expectedTodo, nil)

    result, err := service.CreateTodo("user-id", req)

    assert.NoError(t, err)
    assert.Equal(t, expectedTodo, result)
    mockRepo.AssertExpectations(t)
}
```

## 環境変数
```env
# .env
# Server
PORT=8080
HOST=localhost

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=todoapp
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h

# Environment
GO_ENV=development
```

## ログ設定
```go
// Structured logging with logrus
log := logrus.New()
log.SetFormatter(&logrus.JSONFormatter{})
log.SetLevel(logrus.InfoLevel)

// ミドルウェアでのリクエストログ
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        
        logrus.WithFields(logrus.Fields{
            "status":     c.Writer.Status(),
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "ip":         c.ClientIP(),
            "duration":   time.Since(start),
            "user_agent": c.Request.UserAgent(),
        }).Info("Request processed")
    }
}
```

## パフォーマンス最適化
- データベースインデックスの適切な設計
- N+1問題の回避（Preloadの活用）
- ページネーション実装
- レスポンスキャッシュ（Redis）

## 命名規則
- パッケージ: lowercase
- 関数・メソッド: CamelCase（公開）、camelCase（非公開）
- 構造体: PascalCase（公開）、camelCase（非公開）
- 定数: UPPER_SNAKE_CASE または PascalCase
- ファイル: snake_case

## Git関連
- feature/* ブランチで開発
- go fmt, go vet, golangci-lint を実行
- テストカバレッジ80%以上を目標
