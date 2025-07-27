'use client';

import { useState } from 'react';
import { Todo, TodoFilters } from '../../../types/todo';
import { TodoItem } from './TodoItem';
import { Button } from '../../ui/Button';
import { Input } from '../../ui/Input';
import { Select } from '../../ui/Select';

interface TodoListProps {
  todos: Todo[];
  isLoading?: boolean;
  onToggle: (id: string, completed: boolean) => Promise<void>;
  onEdit: (todo: Todo) => void;
  onDelete: (id: string) => Promise<void>;
  onFiltersChange: (filters: TodoFilters) => void;
  filters: TodoFilters;
}

const priorityFilterOptions = [
  { value: '', label: 'すべての優先度' },
  { value: 'low', label: '低' },
  { value: 'medium', label: '中' },
  { value: 'high', label: '高' },
];

const completedFilterOptions = [
  { value: '', label: 'すべて' },
  { value: 'false', label: '未完了' },
  { value: 'true', label: '完了済み' },
];

export function TodoList({
  todos,
  isLoading,
  onToggle,
  onEdit,
  onDelete,
  onFiltersChange,
  filters,
}: TodoListProps) {
  const [searchTerm, setSearchTerm] = useState(filters.search || '');

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchTerm(value);
    onFiltersChange({ ...filters, search: value || undefined });
  };

  const handlePriorityChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    onFiltersChange({
      ...filters,
      priority: value ? (value as Todo['priority']) : undefined,
    });
  };

  const handleCompletedChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value;
    onFiltersChange({
      ...filters,
      completed: value === '' ? undefined : value === 'true',
    });
  };

  const clearFilters = () => {
    setSearchTerm('');
    onFiltersChange({});
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* フィルター */}
      <div className="bg-white p-4 rounded-lg border border-gray-200 shadow-sm">
        <h3 className="text-lg font-medium text-gray-900 mb-4">フィルター</h3>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Input
            placeholder="検索..."
            value={searchTerm}
            onChange={handleSearchChange}
          />
          <Select
            options={priorityFilterOptions}
            value={filters.priority || ''}
            onChange={handlePriorityChange}
            placeholder="優先度"
          />
          <Select
            options={completedFilterOptions}
            value={
              filters.completed === undefined
                ? ''
                : filters.completed.toString()
            }
            onChange={handleCompletedChange}
            placeholder="完了状態"
          />
          <Button variant="outline" onClick={clearFilters}>
            クリア
          </Button>
        </div>
      </div>

      {/* Todoリスト */}
      <div className="space-y-4">
        {todos.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-gray-500 text-lg">
              {Object.keys(filters).length > 0
                ? 'フィルター条件に一致するTodoがありません'
                : 'まだTodoがありません'}
            </div>
            <p className="text-gray-400 mt-2">
              新しいTodoを作成してみましょう！
            </p>
          </div>
        ) : (
          todos.map((todo) => (
            <TodoItem
              key={todo.id}
              todo={todo}
              onToggle={onToggle}
              onEdit={onEdit}
              onDelete={onDelete}
            />
          ))
        )}
      </div>
    </div>
  );
}
