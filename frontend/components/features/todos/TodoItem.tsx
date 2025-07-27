'use client';

import { useState } from 'react';
import { Todo } from '../../../types/todo';
import { Button } from '../../ui/Button';
import { formatDate, getPriorityColor, getPriorityLabel } from '../../../lib/utils';

interface TodoItemProps {
  todo: Todo;
  onToggle: (id: string, completed: boolean) => Promise<void>;
  onEdit: (todo: Todo) => void;
  onDelete: (id: string) => Promise<void>;
}

export function TodoItem({ todo, onToggle, onEdit, onDelete }: TodoItemProps) {
  const [isToggling, setIsToggling] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const handleToggle = async () => {
    setIsToggling(true);
    try {
      await onToggle(todo.id, !todo.completed);
    } finally {
      setIsToggling(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('このTodoを削除してもよろしいですか？')) {
      setIsDeleting(true);
      try {
        await onDelete(todo.id);
      } finally {
        setIsDeleting(false);
      }
    }
  };

  const isOverdue = todo.dueDate && new Date(todo.dueDate) < new Date() && !todo.completed;

  return (
    <div
      className={`p-4 border rounded-lg shadow-sm transition-colors ${
        todo.completed ? 'bg-gray-50 border-gray-200' : 'bg-white border-gray-300'
      } ${isOverdue ? 'border-red-300 bg-red-50' : ''}`}
    >
      <div className="flex items-start justify-between">
        <div className="flex items-start space-x-3 flex-1">
          <input
            type="checkbox"
            checked={todo.completed}
            onChange={handleToggle}
            disabled={isToggling}
            className="mt-1 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
          />
          
          <div className="flex-1 min-w-0">
            <h3
              className={`text-lg font-medium ${
                todo.completed ? 'line-through text-gray-500' : 'text-gray-900'
              }`}
            >
              {todo.title}
            </h3>
            
            {todo.description && (
              <p
                className={`mt-1 text-sm ${
                  todo.completed ? 'text-gray-400' : 'text-gray-600'
                }`}
              >
                {todo.description}
              </p>
            )}
            
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <span
                className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getPriorityColor(
                  todo.priority
                )}`}
              >
                優先度: {getPriorityLabel(todo.priority)}
              </span>
              
              {todo.dueDate && (
                <span
                  className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    isOverdue
                      ? 'bg-red-100 text-red-800'
                      : 'bg-blue-100 text-blue-800'
                  }`}
                >
                  期限: {formatDate(todo.dueDate)}
                </span>
              )}
              
              <span className="text-xs text-gray-500">
                作成日: {formatDate(todo.createdAt)}
              </span>
            </div>
          </div>
        </div>

        <div className="flex items-center space-x-2 ml-4">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onEdit(todo)}
            disabled={isToggling || isDeleting}
          >
            編集
          </Button>
          <Button
            variant="destructive"
            size="sm"
            onClick={handleDelete}
            isLoading={isDeleting}
            disabled={isToggling || isDeleting}
          >
            削除
          </Button>
        </div>
      </div>
    </div>
  );
}
