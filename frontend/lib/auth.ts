import { apiClient } from './api';
import { User, LoginRequest, RegisterRequest, AuthResponse } from '../types/user';

export class AuthAPI {
  // User registration
  static async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/register', data);
    
    // Save token
    if (response.token) {
      apiClient.setToken(response.token);
    }
    
    return response;
  }

  // Login
  static async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/login', data);
    
    // Save token
    if (response.token) {
      apiClient.setToken(response.token);
    }
  
    return response;
  }

  // Logout
  static async logout(): Promise<void> {
    // Remove token
    apiClient.clearToken();
    
    // Server-side logout process if necessary
    try {
      await apiClient.post('/auth/logout');
    } catch (error) {
      // Ignore logout errors (e.g., when token is already invalid)
      console.warn('Logout error:', error);
    }
  }

  // Token refresh
  static async refreshToken(): Promise<AuthResponse> {
    return apiClient.post<AuthResponse>('/auth/refresh');
  }

  // Get current user information
  static async getCurrentUser(): Promise<User> {
    return apiClient.get<User>('/users/me');
  }

  // Update user information
  static async updateUser(data: Partial<User>): Promise<User> {
    return apiClient.put<User>('/users/me', data);
  }

  // Change password
  static async changePassword(currentPassword: string, newPassword: string): Promise<void> {
    return apiClient.post<void>('/auth/change-password', {
      currentPassword,
      newPassword,
    });
  }
}

export default AuthAPI;
