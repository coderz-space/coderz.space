import axios, { AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';

// ── Config ──────────────────────────────────────────────────────
// Swap BASE_URL when Go backend is ready. All other code stays the same.
const BASE_URL = __DEV__
  ? 'http://localhost:8080/api/v1'
  : 'https://api.coderzspace.com/api/v1';

// ── Create instance ─────────────────────────────────────────────
const apiClient: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: 15000,
  withCredentials: true,      // sends HttpOnly cookies (refresh token)
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
});

// ── Request interceptor: inject JWT ─────────────────────────────
apiClient.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    const token = await AsyncStorage.getItem('@access_token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// ── Response interceptor: handle 401, refresh token ─────────────
let isRefreshing = false;
let refreshSubscribers: Array<(token: string) => void> = [];

const subscribeTokenRefresh = (cb: (token: string) => void) => {
  refreshSubscribers.push(cb);
};

const onRefreshed = (token: string) => {
  refreshSubscribers.forEach((cb) => cb(token));
  refreshSubscribers = [];
};

apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve) => {
          subscribeTokenRefresh((token) => {
            if (originalRequest.headers) {
              originalRequest.headers['Authorization'] = `Bearer ${token}`;
            }
            resolve(apiClient(originalRequest));
          });
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // Refresh token is sent via HttpOnly cookie automatically
        const { data } = await apiClient.post('/auth/refresh');
        const newToken = data.accessToken;
        await AsyncStorage.setItem('@access_token', newToken);
        onRefreshed(newToken);
        isRefreshing = false;
        if (originalRequest.headers) {
          originalRequest.headers['Authorization'] = `Bearer ${newToken}`;
        }
        return apiClient(originalRequest);
      } catch (refreshError) {
        isRefreshing = false;
        // Clear session — force re-login
        await AsyncStorage.removeItem('@access_token');
        // Import dynamically to avoid circular dependency
        const { useAuthStore } = await import('../../store/authStore');
        useAuthStore.getState().logout();
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  },
);

export default apiClient;

// ── URL builder helpers ─────────────────────────────────────────
// Centralize path construction so swapping routes only changes here

export const ApiRoutes = {
  // Auth
  login: '/auth/login',
  logout: '/auth/logout',
  refresh: '/auth/refresh',
  me: '/auth/me',

  // Org
  orgMembers: (orgId: string) => `/orgs/${orgId}/members`,

  // Bootcamp
  bootcampMembers: (orgId: string, bootcampId: string) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/members`,

  // Problems
  problems: (orgId: string) => `/orgs/${orgId}/problems`,
  problemDetail: (orgId: string, problemId: string) =>
    `/orgs/${orgId}/problems/${problemId}`,

  // Assignment Groups
  assignmentGroups: (orgId: string, bootcampId: string) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/assignment-groups`,

  // Assignments
  assignments: (orgId: string, bootcampId: string) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/assignments`,
  assignmentDetail: (orgId: string, bootcampId: string, assignmentId: string) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/assignments/${assignmentId}`,
  enrollmentAssignments: (orgId: string, bootcampId: string, enrollmentId: string) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/enrollments/${enrollmentId}/assignments`,

  // Assignment Problems
  assignmentProblemUpdate: (
    orgId: string,
    bootcampId: string,
    assignmentId: string,
    apId: string,
  ) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/assignments/${assignmentId}/problems/${apId}`,

  // Doubts
  raiseDoubt: (orgId: string, bootcampId: string, assignmentId: string, apId: string) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/assignments/${assignmentId}/problems/${apId}/doubts`,
  doubts: (orgId: string, bootcampId: string) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/doubts`,
  resolveDoubt: (orgId: string, bootcampId: string, doubtId: string) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/doubts/${doubtId}/resolve`,

  // Leaderboard
  leaderboard: (orgId: string, bootcampId: string) =>
    `/orgs/${orgId}/bootcamps/${bootcampId}/leaderboard`,
} as const;