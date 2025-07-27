import { apiClient } from './api';
import { Todo, CreateTodoRequest, UpdateTodoRequest, TodoFilters, PaginatedResponse } from '../types/todo';

export class TodoAPI {
  // Get Todo list
  static async getTodos(
    page: number = 1,
    limit: number = 10,
    filters?: TodoFilters
  ): Promise<PaginatedResponse<Todo>> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });

    if (filters?.completed !== undefined) {
      params.append('completed', filters.completed.toString());
    }
    if (filters?.priority) {
      params.append('priority', filters.priority);
    }
    if (filters?.search) {
      params.append('search', filters.search);
    }

    return apiClient.get<PaginatedResponse<Todo>>(`/todos?${params.toString()}`);
  }

  // Get Todo details
  static async getTodo(id: string): Promise<Todo> {
    return apiClient.get<Todo>(`/todos/${id}`);
  }

  // Create Todo
  static async createTodo(data: CreateTodoRequest): Promise<Todo> {
    return apiClient.post<Todo>('/todos', data);
  }

  // Update Todo
  static async updateTodo(id: string, data: UpdateTodoRequest): Promise<Todo> {
    return apiClient.put<Todo>(`/todos/${id}`, data);
  }

  // Delete Todo
  static async deleteTodo(id: string): Promise<void> {
    return apiClient.delete<void>(`/todos/${id}`);
  }

  // Toggle Todo completion status
  static async toggleTodo(id: string, completed: boolean): Promise<Todo> {
    return apiClient.patch<Todo>(`/todos/${id}`, { completed });
  }
}

export default TodoAPI;
