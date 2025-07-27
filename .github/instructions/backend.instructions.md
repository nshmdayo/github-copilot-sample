# Backend Development Instructions - Go Language

## Project Overview
Build the backend API for a Web Todo application using Go language. Provide RESTful API and implement authentication and Todo CRUD operations.

## Technology Stack
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

## Project Structure
```
backend/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/                     # Private application code
│   ├── config/                   # Configuration management
│   │   └── config.go
│   ├── handler/                  # HTTP handlers
│   │   ├── auth.go
│   │   ├── todo.go
│   │   └── user.go
│   ├── middleware/               # HTTP middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── model/                    # Data models
│   │   ├── todo.go
│   │   └── user.go
│   ├── repository/               # Data access layer
│   │   ├── todo.go
│   │   └── user.go
│   ├── service/                  # Business logic layer
│   │   ├── auth.go
│   │   ├── todo.go
│   │   └── user.go
│   └── utils/                    # Utility functions
│       ├── hash.go
│       ├── jwt.go
│       └── validator.go
├── pkg/                          # External packages
│   └── database/
│       └── postgres.go
├── migrations/                   # Database migrations
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_todos_table.up.sql
│   └── 002_create_todos_table.down.sql
├── docs/                         # Swagger generated documentation
├── docker/
│   └── Dockerfile
├── scripts/                      # Build & deployment scripts
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## Architecture Pattern
Adopts Clean Architecture with clear separation of dependencies.

### Layer Composition
1. **Handler Layer**: HTTP request/response processing
2. **Service Layer**: Business logic
3. **Repository Layer**: Data access
4. **Model Layer**: Domain models

### Dependency Injection
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

## Data Models

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

## API Design

### Authentication Endpoints
```go

### Authentication Endpoints

// POST /api/auth/register - User registration
// POST /api/auth/login    - Login
// POST /api/auth/refresh  - Token refresh

// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Register request"
// @Success 201 {object} model.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    // Implementation
}
```

### Todo Endpoints
```go
// GET    /api/todos      - Get Todo list
// POST   /api/todos      - Create Todo
// GET    /api/todos/:id  - Get Todo details
// PUT    /api/todos/:id  - Update Todo
// DELETE /api/todos/:id  - Delete Todo

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
    // Implementation
}
```

## Error Handling

### Standard Error Response
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
        response.Message = "Input data is invalid"
        response.Details = formatValidationErrors(e)
        statusCode = http.StatusBadRequest
    default:
        response.Error = "internal_error"
        response.Message = err.Error()
    }
    
    c.JSON(statusCode, response)
}
```

## Security

### JWT Authentication
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

### Middleware
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

## Configuration Management
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

## Testing Strategy

### Unit Testing
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

## Environment Variables
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

## Logging Configuration
```go
// Structured logging with logrus
log := logrus.New()
log.SetFormatter(&logrus.JSONFormatter{})
log.SetLevel(logrus.InfoLevel)

// Request logging middleware
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

## Performance Optimization
- Proper database index design
- Avoiding N+1 problems (using Preload)
- Pagination implementation
- Response caching (Redis)

## Naming Conventions
- Package: lowercase
- Functions/Methods: CamelCase (public), camelCase (private)
- Structs: PascalCase (public), camelCase (private)
- Constants: UPPER_SNAKE_CASE or PascalCase
- Files: snake_case

## Git Guidelines
- Develop in feature/* branches
- Run go fmt, go vet, golangci-lint
- Aim for 80%+ test coverage
