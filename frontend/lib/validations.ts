import { z } from 'zod';

// Todo関連のバリデーションスキーマ
export const createTodoSchema = z.object({
  title: z.string()
    .min(1, 'タイトルは必須です')
    .max(200, 'タイトルは200文字以内で入力してください'),
  description: z.string()
    .max(1000, '説明は1000文字以内で入力してください')
    .optional(),
  priority: z.enum(['low', 'medium', 'high']),
  dueDate: z.string()
    .optional()
    .refine((date) => {
      if (!date) return true;
      const parsedDate = new Date(date);
      return parsedDate >= new Date();
    }, '期限は今日以降の日付を設定してください'),
});

export const updateTodoSchema = createTodoSchema.partial().extend({
  completed: z.boolean().optional(),
});

// ユーザー関連のバリデーションスキーマ
export const loginSchema = z.object({
  email: z.string()
    .min(1, 'メールアドレスは必須です')
    .email('有効なメールアドレスを入力してください'),
  password: z.string()
    .min(1, 'パスワードは必須です')
    .min(8, 'パスワードは8文字以上で入力してください'),
});

export const registerSchema = loginSchema.extend({
  name: z.string()
    .min(1, 'お名前は必須です')
    .min(2, 'お名前は2文字以上で入力してください')
    .max(100, 'お名前は100文字以内で入力してください'),
  confirmPassword: z.string()
    .min(1, 'パスワード確認は必須です'),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'パスワードが一致しません',
  path: ['confirmPassword'],
});

export const updateUserSchema = z.object({
  name: z.string()
    .min(2, 'お名前は2文字以上で入力してください')
    .max(100, 'お名前は100文字以内で入力してください')
    .optional(),
  email: z.string()
    .email('有効なメールアドレスを入力してください')
    .optional(),
});

export const changePasswordSchema = z.object({
  currentPassword: z.string()
    .min(1, '現在のパスワードは必須です'),
  newPassword: z.string()
    .min(8, 'パスワードは8文字以上で入力してください'),
  confirmPassword: z.string()
    .min(1, 'パスワード確認は必須です'),
}).refine((data) => data.newPassword === data.confirmPassword, {
  message: 'パスワードが一致しません',
  path: ['confirmPassword'],
});

// フィルター関連のバリデーションスキーマ
export const todoFiltersSchema = z.object({
  search: z.string().optional(),
  completed: z.boolean().optional(),
  priority: z.enum(['low', 'medium', 'high']).optional(),
});

// 型エクスポート
export type CreateTodoInput = z.infer<typeof createTodoSchema>;
export type UpdateTodoInput = z.infer<typeof updateTodoSchema>;
export type LoginInput = z.infer<typeof loginSchema>;
export type RegisterInput = z.infer<typeof registerSchema>;
export type UpdateUserInput = z.infer<typeof updateUserSchema>;
export type ChangePasswordInput = z.infer<typeof changePasswordSchema>;
export type TodoFiltersInput = z.infer<typeof todoFiltersSchema>;
