package service

import (
	"errors"
	"math"

	"github.com/nshmdayo/github-copilot-sample/backend/internal/model"
	"github.com/nshmdayo/github-copilot-sample/backend/internal/repository"
	"gorm.io/gorm"
)

// TodoService defines the interface for todo operations
type TodoService interface {
	Create(userID string, req *model.CreateTodoRequest) (*model.Todo, error)
	GetByID(userID, todoID string) (*model.Todo, error)
	GetList(userID string, req *model.TodoListRequest) (*model.TodoListResponse, error)
	Update(userID, todoID string, req *model.UpdateTodoRequest) (*model.Todo, error)
	Delete(userID, todoID string) error
	ToggleStatus(userID, todoID string) (*model.Todo, error)
}

// todoService implements TodoService interface
type todoService struct {
	todoRepo repository.TodoRepository
}

// NewTodoService creates a new todo service
func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}

// Create creates a new todo
func (s *todoService) Create(userID string, req *model.CreateTodoRequest) (*model.Todo, error) {
	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      model.StatusPending,
		UserID:      userID,
		DueDate:     req.DueDate,
	}

	if err := s.todoRepo.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// GetByID retrieves a todo by ID for a specific user
func (s *todoService) GetByID(userID, todoID string) (*model.Todo, error) {
	todo, err := s.todoRepo.GetUserTodoByID(userID, todoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	return todo, nil
}

// GetList retrieves todos for a user with pagination and filters
func (s *todoService) GetList(userID string, req *model.TodoListRequest) (*model.TodoListResponse, error) {
	// Set default values
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	todos, total, err := s.todoRepo.GetByUserID(userID, req)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	todoResponses := make([]model.TodoResponse, len(todos))
	for i, todo := range todos {
		todoResponses[i] = todo.ToResponse(false)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	response := &model.TodoListResponse{
		Data:       todoResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}

	return response, nil
}

// Update updates a todo
func (s *todoService) Update(userID, todoID string, req *model.UpdateTodoRequest) (*model.Todo, error) {
	todo, err := s.todoRepo.GetUserTodoByID(userID, todoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	if req.Status != nil {
		todo.Status = *req.Status
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}

	if err := s.todoRepo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// Delete deletes a todo
func (s *todoService) Delete(userID, todoID string) error {
	// Check if todo exists and belongs to user
	_, err := s.todoRepo.GetUserTodoByID(userID, todoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("todo not found")
		}
		return err
	}

	return s.todoRepo.Delete(todoID)
}

// ToggleStatus toggles the completion status of a todo
func (s *todoService) ToggleStatus(userID, todoID string) (*model.Todo, error) {
	todo, err := s.todoRepo.GetUserTodoByID(userID, todoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	// Toggle status
	if todo.IsCompleted() {
		todo.MarkAsPending()
	} else {
		todo.MarkAsCompleted()
	}

	if err := s.todoRepo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}
