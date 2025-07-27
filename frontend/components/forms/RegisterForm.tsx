'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { registerSchema, RegisterInput } from '../../lib/validations';
import { Button } from '../ui/Button';
import { Input } from '../ui/Input';

interface RegisterFormProps {
  onSubmit: (data: Omit<RegisterInput, 'confirmPassword'>) => Promise<void>;
  onToggleForm: () => void;
}

export function RegisterForm({ onSubmit, onToggleForm }: RegisterFormProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterInput>({
    resolver: zodResolver(registerSchema),
  });

  const handleFormSubmit = async (data: RegisterInput) => {
    setIsLoading(true);
    setError(null);
    
    try {
      // confirmPasswordを除いてAPIに送信
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { confirmPassword, ...submitData } = data;
      await onSubmit(submitData);
    } catch (err) {
      setError(err instanceof Error ? err.message : '予期しないエラーが発生しました');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md mx-auto">
      <div className="bg-white py-8 px-6 shadow-lg rounded-lg">
        <div className="text-center mb-6">
          <h2 className="text-3xl font-bold text-gray-900">新規登録</h2>
          <p className="mt-2 text-gray-600">新しいアカウントを作成してください</p>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <Input
            type="text"
            label="お名前"
            placeholder="山田太郎"
            error={errors.name?.message}
            {...register('name')}
          />

          <Input
            type="email"
            label="メールアドレス"
            placeholder="your-email@example.com"
            error={errors.email?.message}
            {...register('email')}
          />

          <Input
            type="password"
            label="パスワード"
            placeholder="8文字以上のパスワード"
            error={errors.password?.message}
            {...register('password')}
          />

          <Input
            type="password"
            label="パスワード確認"
            placeholder="同じパスワードを入力してください"
            error={errors.confirmPassword?.message}
            {...register('confirmPassword')}
          />

          <Button
            type="submit"
            className="w-full"
            isLoading={isLoading}
            disabled={isLoading}
          >
            アカウント作成
          </Button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-sm text-gray-600">
            既にアカウントをお持ちの方は{' '}
            <button
              type="button"
              onClick={onToggleForm}
              className="font-medium text-blue-600 hover:text-blue-500"
            >
              ログイン
            </button>
          </p>
        </div>
      </div>
    </div>
  );
}
