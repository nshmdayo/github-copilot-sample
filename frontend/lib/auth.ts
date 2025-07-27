import { apiClient } from './api';
import { User, LoginRequest, RegisterRequest, AuthResponse } from '../types/user';

export class AuthAPI {
  // ユーザー登録
  static async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/register', data);
    
    // トークンを保存
    if (response.token) {
      apiClient.setToken(response.token);
    }
    
    return response;
  }

  // ログイン
  static async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/login', data);
    
    // トークンを保存
    if (response.token) {
      apiClient.setToken(response.token);
    }
    
    return response;
  }

  // ログアウト
  static async logout(): Promise<void> {
    // トークンを削除
    apiClient.clearToken();
    
    // 必要に応じてサーバーサイドでのログアウト処理
    try {
      await apiClient.post('/auth/logout');
    } catch (error) {
      // ログアウトエラーは無視（既にトークンが無効な場合など）
      console.warn('Logout error:', error);
    }
  }

  // トークンリフレッシュ
  static async refreshToken(): Promise<AuthResponse> {
    return apiClient.post<AuthResponse>('/auth/refresh');
  }

  // 現在のユーザー情報を取得
  static async getCurrentUser(): Promise<User> {
    return apiClient.get<User>('/users/me');
  }

  // ユーザー情報を更新
  static async updateUser(data: Partial<User>): Promise<User> {
    return apiClient.put<User>('/users/me', data);
  }

  // パスワード変更
  static async changePassword(currentPassword: string, newPassword: string): Promise<void> {
    return apiClient.post<void>('/auth/change-password', {
      currentPassword,
      newPassword,
    });
  }
}

export default AuthAPI;
