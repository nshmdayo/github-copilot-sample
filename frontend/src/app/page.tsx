'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { TodoAPI } from '../../lib/todos';
import { Todo, TodoFilters, CreateTodoRequest } from '../../types/todo';
import { Layout } from '../../components/layouts/Layout';
import { TodoForm } from '../../components/forms/TodoForm';
import { TodoList } from '../../components/features/todos/TodoList';
import { LoginForm } from '../../components/forms/LoginForm';
import { RegisterForm } from '../../components/forms/RegisterForm';
import { Button } from '../../components/ui/Button';

export default function Home() {
  const { isAuthenticated, isLoading: authLoading, login, register } = useAuth();
  const [todos, setTodos] = useState<Todo[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const [filters, setFilters] = useState<TodoFilters>({});
  const [isLoginMode, setIsLoginMode] = useState(true);

  // Fetch Todo list
  const fetchTodos = async () => {
    if (!isAuthenticated) return;
    
    setIsLoading(true);
    try {
      const response = await TodoAPI.getTodos(1, 50, filters);
      setTodos(response.data);
    } catch (error) {
      console.error('Failed to fetch todos:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Re-fetch Todo list when authentication state and filter change
  useEffect(() => {
    fetchTodos();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, filters]);

  // Create Todo
  const handleCreateTodo = async (data: CreateTodoRequest) => {
    try {
      const newTodo = await TodoAPI.createTodo(data);
      setTodos(prevTodos => [newTodo, ...prevTodos]);
      setShowForm(false);
    } catch (error) {
      console.error('Failed to create todo:', error);
      throw error;
    }
  };

  // Toggle Todo completion status
  const handleToggleTodo = async (id: string, completed: boolean) => {
    try {
      const updatedTodo = await TodoAPI.toggleTodo(id, completed);
      setTodos(prevTodos =>
        prevTodos.map(todo =>
          todo.id === id ? updatedTodo : todo
        )
      );
    } catch (error) {
      console.error('Failed to toggle todo:', error);
      throw error;
    }
  };

  // Todo editing (edit form omitted in this sample)
  const handleEditTodo = (todo: Todo) => {
    console.log('Edit todo:', todo);
    // TODO: Implement edit form
  };

  // Delete Todo
  const handleDeleteTodo = async (id: string) => {
    try {
      await TodoAPI.deleteTodo(id);
      setTodos(prevTodos => prevTodos.filter(todo => todo.id !== id));
    } catch (error) {
      console.error('Failed to delete todo:', error);
      throw error;
    }
  };

  // Login process
  const handleLogin = async (data: { email: string; password: string }) => {
    await login(data.email, data.password);
  };

  // Registration process
  const handleRegister = async (data: { name: string; email: string; password: string }) => {
    await register(data.name, data.email, data.password);
  };

  // Loading authentication state
  if (authLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  // Show login/registration form if not authenticated
  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
        {isLoginMode ? (
          <LoginForm
            onSubmit={handleLogin}
            onToggleForm={() => setIsLoginMode(false)}
          />
        ) : (
          <RegisterForm
            onSubmit={handleRegister}
            onToggleForm={() => setIsLoginMode(true)}
          />
        )}
      </div>
    );
  }

  // Show Todo app if authenticated
  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h2 className="text-2xl font-bold text-gray-900">My Todos</h2>
          <Button onClick={() => setShowForm(!showForm)}>
            {showForm ? 'Close Form' : 'Create New Todo'}
          </Button>
        </div>

        {showForm && (
          <div className="bg-white p-6 rounded-lg border border-gray-200 shadow-sm">
            <h3 className="text-lg font-medium text-gray-900 mb-4">
              Create New Todo
            </h3>
            <TodoForm onSubmit={handleCreateTodo} />
          </div>
        )}

        <TodoList
          todos={todos}
          isLoading={isLoading}
          onToggle={handleToggleTodo}
          onEdit={handleEditTodo}
          onDelete={handleDeleteTodo}
          onFiltersChange={setFilters}
          filters={filters}
        />
      </div>
    </Layout>
  );
}
