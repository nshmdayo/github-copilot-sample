'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { createTodoSchema, CreateTodoInput } from '../../lib/validations';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';
import { Textarea } from '../ui/Textarea';
import { Select } from '../ui/Select';

interface TodoFormProps {
  onSubmit: (data: CreateTodoInput) => Promise<void>;
  isLoading?: boolean;
  defaultValues?: Partial<CreateTodoInput>;
}

const priorityOptions = [
  { value: 'low', label: '低' },
  { value: 'medium', label: '中' },
  { value: 'high', label: '高' },
];

export function TodoForm({ onSubmit, isLoading, defaultValues }: TodoFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<CreateTodoInput>({
    resolver: zodResolver(createTodoSchema),
    defaultValues: {
      priority: 'medium',
      ...defaultValues,
    },
  });

  const handleFormSubmit = async (data: CreateTodoInput) => {
    try {
      await onSubmit(data);
      reset();
    } catch (error) {
      console.error('Form submission error:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <Input
        label="タイトル"
        placeholder="Todoのタイトルを入力してください"
        error={errors.title?.message}
        {...register('title')}
      />

      <Textarea
        label="説明（任意）"
        placeholder="詳細な説明を入力してください"
        error={errors.description?.message}
        {...register('description')}
        rows={3}
      />

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Select
          label="優先度"
          options={priorityOptions}
          error={errors.priority?.message}
          {...register('priority')}
        />

        <Input
          type="date"
          label="期限（任意）"
          error={errors.dueDate?.message}
          {...register('dueDate')}
        />
      </div>

      <div className="flex justify-end space-x-3">
        <Button
          type="button"
          variant="outline"
          onClick={() => reset()}
          disabled={isLoading}
        >
          リセット
        </Button>
        <Button
          type="submit"
          isLoading={isLoading}
          disabled={isLoading}
        >
          Todoを作成
        </Button>
      </div>
    </form>
  );
}
