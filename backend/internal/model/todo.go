package model

import (
	"time"

	"gorm.io/gorm"
)

// Priority represents the priority level of a todo
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Status represents the status of a todo
type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
)

// Todo represents a todo item in the system
type Todo struct {
	ID          string         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title" validate:"required,min=1,max=200"`
	Description string         `json:"description" validate:"max=1000"`
	Priority    Priority       `gorm:"type:varchar(10);default:'medium'" json:"priority" validate:"oneof=low medium high"`
	Status      Status         `gorm:"type:varchar(20);default:'pending'" json:"status" validate:"oneof=pending completed"`
	UserID      string         `gorm:"type:uuid;not null;index" json:"user_id"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name for Todo model
func (Todo) TableName() string {
	return "todos"
}

// IsCompleted returns true if the todo is completed
func (t *Todo) IsCompleted() bool {
	return t.Status == StatusCompleted
}

// MarkAsCompleted marks the todo as completed
func (t *Todo) MarkAsCompleted() {
	t.Status = StatusCompleted
}

// MarkAsPending marks the todo as pending
func (t *Todo) MarkAsPending() {
	t.Status = StatusPending
}

// CreateTodoRequest represents the request payload for creating a todo
type CreateTodoRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=200" example:"Buy groceries"`
	Description string     `json:"description" validate:"max=1000" example:"Buy milk, eggs, and bread"`
	Priority    Priority   `json:"priority" validate:"oneof=low medium high" example:"medium"`
	DueDate     *time.Time `json:"due_date,omitempty" example:"2024-02-01T10:00:00Z"`
}

// UpdateTodoRequest represents the request payload for updating a todo
type UpdateTodoRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=200" example:"Buy groceries"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=1000" example:"Buy milk, eggs, and bread"`
	Priority    *Priority  `json:"priority,omitempty" validate:"omitempty,oneof=low medium high" example:"high"`
	Status      *Status    `json:"status,omitempty" validate:"omitempty,oneof=pending completed" example:"completed"`
	DueDate     *time.Time `json:"due_date,omitempty" example:"2024-02-01T10:00:00Z"`
}

// TodoResponse represents the response payload for todo data
type TodoResponse struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Priority    Priority      `json:"priority"`
	Status      Status        `json:"status"`
	UserID      string        `json:"user_id"`
	DueDate     *time.Time    `json:"due_date,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	User        *UserResponse `json:"user,omitempty"`
}

// ToResponse converts Todo to TodoResponse
func (t *Todo) ToResponse(includeUser bool) TodoResponse {
	response := TodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Priority:    t.Priority,
		Status:      t.Status,
		UserID:      t.UserID,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	if includeUser {
		userResponse := t.User.ToResponse()
		response.User = &userResponse
	}

	return response
}

// TodoListRequest represents the request parameters for listing todos
type TodoListRequest struct {
	Page     int       `form:"page" validate:"min=1" example:"1"`
	Limit    int       `form:"limit" validate:"min=1,max=100" example:"10"`
	Status   *Status   `form:"status" validate:"omitempty,oneof=pending completed" example:"pending"`
	Priority *Priority `form:"priority" validate:"omitempty,oneof=low medium high" example:"high"`
	Search   string    `form:"search" validate:"max=200" example:"groceries"`
}

// TodoListResponse represents the response payload for todo list
type TodoListResponse struct {
	Data       []TodoResponse `json:"data"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}
