package repository

import (
	"github.com/nshmdayo/github-copilot-sample/backend/internal/model"
	"gorm.io/gorm"
)

// TodoRepository defines the interface for todo data operations
type TodoRepository interface {
	Create(todo *model.Todo) error
	GetByID(id string) (*model.Todo, error)
	GetByUserID(userID string, req *model.TodoListRequest) ([]model.Todo, int64, error)
	Update(todo *model.Todo) error
	Delete(id string) error
	GetUserTodoByID(userID, todoID string) (*model.Todo, error)
}

// todoRepository implements TodoRepository interface
type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new todo repository
func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

// Create creates a new todo
func (r *todoRepository) Create(todo *model.Todo) error {
	return r.db.Create(todo).Error
}

// GetByID retrieves a todo by ID
func (r *todoRepository) GetByID(id string) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.Preload("User").Where("id = ?", id).First(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// GetUserTodoByID retrieves a todo by ID that belongs to a specific user
func (r *todoRepository) GetUserTodoByID(userID, todoID string) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// GetByUserID retrieves todos by user ID with pagination and filters
func (r *todoRepository) GetByUserID(userID string, req *model.TodoListRequest) ([]model.Todo, int64, error) {
	var todos []model.Todo
	var total int64

	query := r.db.Model(&model.Todo{}).Where("user_id = ?", userID)

	// Apply filters
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Priority != nil {
		query = query.Where("priority = ?", *req.Priority)
	}
	if req.Search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.Limit).Find(&todos).Error; err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

// Update updates a todo
func (r *todoRepository) Update(todo *model.Todo) error {
	return r.db.Save(todo).Error
}

// Delete deletes a todo by ID
func (r *todoRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Todo{}).Error
}
