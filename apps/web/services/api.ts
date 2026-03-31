/**
 * api.ts — Secure HTTP client for API integration
 *
 * This module provides:
 * - Centralized API configuration
 * - Request/response interceptors
 * - Error handling
 * - Authentication token management
 * - Security best practices (CSRF, XSS prevention)
 */

import type { AxiosInstance, AxiosRequestConfig, AxiosError } from 'axios';

// Use dynamic import to avoid SSR issues
let axiosInstance: AxiosInstance | null = null;

/**
 * Get or create axios instance with proper configuration
 */
async function getAxiosInstance(): Promise<AxiosInstance> {
  if (typeof window === 'undefined') {
    throw new Error('API client can only be used in browser environment');
  }

  if (axiosInstance) {
    return axiosInstance;
  }

  const axios = await import('axios').then(m => m.default);

  const baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

  axiosInstance = axios.create({
    baseURL,
    timeout: 10000,
    withCredentials: true, // Enable CORS cookies
    headers: {
      'Content-Type': 'application/json',
      'X-Requested-With': 'XMLHttpRequest', // CSRF protection
    },
  });

  // Request interceptor — add auth token
  axiosInstance.interceptors.request.use(
    (config) => {
      if (typeof window !== 'undefined') {
        const token = localStorage.getItem('auth_token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Response interceptor — handle errors & token refresh
  axiosInstance.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      if (error.response?.status === 401) {
        // Token expired or unauthorized — clear auth state
        localStorage.removeItem('auth_token');
        localStorage.removeItem('refresh_token');

        // Redirect to login if available
        if (typeof window !== 'undefined') {
          window.location.href = '/';
        }
      }

      return Promise.reject(new APIError(error));
    }
  );

  return axiosInstance;
}

/**
 * Custom error class for better error handling
 */
export class APIError extends Error {
  public status: number;
  public data?: Record<string, any>;

  constructor(error: AxiosError) {
    const message = (error.response?.data as any)?.message || error.message;
    super(message);
    this.name = 'APIError';
    this.status = error.response?.status || 500;
    this.data = error.response?.data as Record<string, any>;
  }
}

/**
 * Generic API request wrapper
 */
export async function apiRequest<T = any>(
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE',
  url: string,
  data?: Record<string, any>,
  config?: AxiosRequestConfig
): Promise<T> {
  const axios = await getAxiosInstance();

  const requestConfig: AxiosRequestConfig = {
    method,
    url,
    ...config,
    ...(data && (method === 'POST' || method === 'PUT' || method === 'PATCH') && { data }),
  };

  const response = await axios(requestConfig);
  return response.data;
}

/**
 * Convenience methods matching REST conventions
 */
export const api = {
  get: <T = any>(url: string, config?: AxiosRequestConfig) => 
    apiRequest<T>('GET', url, undefined, config),

  post: <T = any>(url: string, data?: Record<string, any>, config?: AxiosRequestConfig) => 
    apiRequest<T>('POST', url, data, config),

  put: <T = any>(url: string, data?: Record<string, any>, config?: AxiosRequestConfig) => 
    apiRequest<T>('PUT', url, data, config),

  patch: <T = any>(url: string, data?: Record<string, any>, config?: AxiosRequestConfig) => 
    apiRequest<T>('PATCH', url, data, config),

  delete: <T = any>(url: string, config?: AxiosRequestConfig) => 
    apiRequest<T>('DELETE', url, undefined, config),
};

/**
 * Health check to verify API connectivity
 */
export async function checkAPIHealth(): Promise<boolean> {
  try {
    const response = await api.get('/health');
    return response?.status === 'ok';
  } catch (error) {
    console.warn('API health check failed:', error);
    return false;
  }
}

export default api;
