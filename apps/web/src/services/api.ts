import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import { message } from 'antd';

// 创建 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
});

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message: msg, data } = response.data;

    if (code === 200) {
      return data;
    }

    if (code === 401) {
      message.error('登录已过期，请重新登录');
      localStorage.removeItem('token');
      window.location.href = '/login';
      return Promise.reject(new Error(msg || '登录已过期'));
    }

    message.error(msg || '请求失败');
    return Promise.reject(new Error(msg || '请求失败'));
  },
  (error) => {
    if (error.response) {
      const { status } = error.response;
      if (status === 401) {
        message.error('登录已过期，请重新登录');
        localStorage.removeItem('token');
        window.location.href = '/login';
      } else if (status === 403) {
        message.error('没有权限访问');
      } else if (status === 500) {
        message.error('服务器错误');
      } else {
        message.error(error.response.data?.message || '请求失败');
      }
    } else {
      message.error('网络错误');
    }
    return Promise.reject(error);
  }
);

export default request;
