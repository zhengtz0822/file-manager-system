import request from './api';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: {
    id: number;
    username: string;
  };
}

export interface RegisterRequest {
  username: string;
  password: string;
}

// 登录
export async function login(data: LoginRequest): Promise<LoginResponse> {
  return request.post('/auth/login', data);
}

// 注册
export async function register(data: RegisterRequest): Promise<void> {
  return request.post('/auth/register', data);
}

// 登出
export async function logout(): Promise<void> {
  return request.post('/auth/logout');
}

// 获取当前用户信息
export async function getCurrentUser(): Promise<LoginResponse['user']> {
  return request.get('/user/me');
}
