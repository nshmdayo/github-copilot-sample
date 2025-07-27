package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nshmdayo/github-copilot-sample/backend/internal/model"
	"github.com/nshmdayo/github-copilot-sample/backend/internal/service"
)

// TodoHandler handles todo related requests
type TodoHandler struct {
	todoService service.TodoService
	validator   *validator.Validate
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
		validator:   validator.New(),
	}
}

// Create handles todo creation
// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param todo body model.CreateTodoRequest true "Todo creation data"
// @Success 201 {object} model.TodoResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var req model.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	todo, err := h.todoService.Create(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, todo.ToResponse(false))
}

// GetList handles todo list retrieval
// @Summary Get todo list
// @Description Get paginated list of todos with optional filters
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status" Enums(pending, completed)
// @Param priority query string false "Filter by priority" Enums(low, medium, high)
// @Param search query string false "Search in title and description"
// @Success 200 {object} model.TodoListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /todos [get]
func (h *TodoHandler) GetList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	var req model.TodoListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	response, err := h.todoService.GetList(userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetByID handles todo retrieval by ID
// @Summary Get todo by ID
// @Description Get a specific todo by its ID
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} model.TodoResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID is required"})
		return
	}

	todo, err := h.todoService.GetByID(userID.(string), todoID)
	if err != nil {
		if err.Error() == "todo not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, todo.ToResponse(false))
}

// Update handles todo updates
// @Summary Update todo
// @Description Update a specific todo by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Param todo body model.UpdateTodoRequest true "Todo update data"
// @Success 200 {object} model.TodoResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID is required"})
		return
	}

	var req model.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	todo, err := h.todoService.Update(userID.(string), todoID, &req)
	if err != nil {
		if err.Error() == "todo not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, todo.ToResponse(false))
}

// Delete handles todo deletion
// @Summary Delete todo
// @Description Delete a specific todo by its ID
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID is required"})
		return
	}

	err := h.todoService.Delete(userID.(string), todoID)
	if err != nil {
		if err.Error() == "todo not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// ToggleStatus handles todo status toggle
// @Summary Toggle todo status
// @Description Toggle the completion status of a todo
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} model.TodoResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /todos/{id}/toggle [patch]
func (h *TodoHandler) ToggleStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	todoID := c.Param("id")
	if todoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID is required"})
		return
	}

	todo, err := h.todoService.ToggleStatus(userID.(string), todoID)
	if err != nil {
		if err.Error() == "todo not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, todo.ToResponse(false))
}
