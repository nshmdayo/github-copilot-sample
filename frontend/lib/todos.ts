import { apiClient } from './api';
import { Todo, CreateTodoRequest, UpdateTodoRequest, TodoFilters, PaginatedResponse } from '../types/todo';

export class TodoAPI {
  // Todo一覧を取得
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

  // Todo詳細を取得
  static async getTodo(id: string): Promise<Todo> {
    return apiClient.get<Todo>(`/todos/${id}`);
  }

  // Todo作成
  static async createTodo(data: CreateTodoRequest): Promise<Todo> {
    return apiClient.post<Todo>('/todos', data);
  }

  // Todo更新
  static async updateTodo(id: string, data: UpdateTodoRequest): Promise<Todo> {
    return apiClient.put<Todo>(`/todos/${id}`, data);
  }

  // Todo削除
  static async deleteTodo(id: string): Promise<void> {
    return apiClient.delete<void>(`/todos/${id}`);
  }

  // Todo完了状態の切り替え
  static async toggleTodo(id: string, completed: boolean): Promise<Todo> {
    return apiClient.patch<Todo>(`/todos/${id}`, { completed });
  }
}

export default TodoAPI;
