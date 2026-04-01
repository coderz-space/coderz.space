import type { AxiosError, AxiosInstance, AxiosRequestConfig } from "axios";

export interface ApiEnvelope<T> {
  success?: boolean;
  status?: string;
  message?: string;
  data?: T;
  error?: {
    code?: number;
    message?: string;
  };
}

let axiosInstance: AxiosInstance | null = null;

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

function hasDataEnvelope<T>(value: unknown): value is ApiEnvelope<T> {
  return isRecord(value) && "data" in value;
}

function extractMessage(value: unknown, fallback: string): string {
  if (isRecord(value)) {
    const message = value.message;
    if (typeof message === "string" && message.trim()) {
      return message;
    }

    const error = value.error;
    if (isRecord(error) && typeof error.message === "string" && error.message.trim()) {
      return error.message;
    }
  }

  return fallback;
}

async function getAxiosInstance(): Promise<AxiosInstance> {
  if (typeof window === "undefined") {
    throw new Error("API client can only be used in the browser.");
  }

  if (axiosInstance) {
    return axiosInstance;
  }

  const axios = await import("axios").then((module) => module.default);
  const baseURL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api";

  axiosInstance = axios.create({
    baseURL,
    timeout: 15000,
    withCredentials: true,
    headers: {
      "Content-Type": "application/json",
      "X-Requested-With": "XMLHttpRequest",
    },
  });

  axiosInstance.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const request = error.config;
      const isRefreshRequest = request?.url?.includes("/v1/auth/refresh");

      if (error.response?.status === 401 && request && !isRefreshRequest) {
        try {
          const axios = await import("axios").then((module) => module.default);
          const base = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api";
          await axios.post(`${base}/v1/auth/refresh`, {}, { withCredentials: true });
          return axiosInstance!(request);
        } catch {
          if (typeof window !== "undefined") {
            window.location.href = "/";
          }
        }
      }

      return Promise.reject(new APIError(error));
    }
  );

  return axiosInstance;
}

export class APIError extends Error {
  public readonly status: number;
  public readonly data?: unknown;

  constructor(error: AxiosError) {
    super(extractMessage(error.response?.data, error.message));
    this.name = "APIError";
    this.status = error.response?.status ?? 500;
    this.data = error.response?.data;
  }
}

export async function requestRaw<T>(config: AxiosRequestConfig): Promise<T> {
  const client = await getAxiosInstance();
  const response = await client.request<T>(config);
  return response.data;
}

export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  const response = await requestRaw<T | ApiEnvelope<T>>(config);
  if (hasDataEnvelope<T>(response)) {
    return response.data as T;
  }
  return response as T;
}

export const api = {
  get: <T>(url: string, config?: AxiosRequestConfig) =>
    request<T>({ ...config, method: "GET", url }),
  post: <T, B extends object | undefined = undefined>(url: string, data?: B, config?: AxiosRequestConfig) =>
    request<T>({ ...config, method: "POST", url, data }),
  put: <T, B extends object | undefined = undefined>(url: string, data?: B, config?: AxiosRequestConfig) =>
    request<T>({ ...config, method: "PUT", url, data }),
  patch: <T, B extends object | undefined = undefined>(url: string, data?: B, config?: AxiosRequestConfig) =>
    request<T>({ ...config, method: "PATCH", url, data }),
  delete: <T>(url: string, config?: AxiosRequestConfig) =>
    request<T>({ ...config, method: "DELETE", url }),
  rawPost: <T, B extends object | undefined = undefined>(url: string, data?: B, config?: AxiosRequestConfig) =>
    requestRaw<T>({ ...config, method: "POST", url, data }),
  rawGet: <T>(url: string, config?: AxiosRequestConfig) =>
    requestRaw<T>({ ...config, method: "GET", url }),
};

export async function checkAPIHealth(): Promise<boolean> {
  try {
    const response = await api.get<{ status?: string }>("/health");
    return response.status === "ok";
  } catch {
    return false;
  }
}
