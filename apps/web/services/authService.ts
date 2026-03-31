/**
 * authService.ts — Authentication with real backend API
 *
 * These functions call the actual backend endpoints:
 * - POST /api/v1/auth/signup   → signup
 * - POST /api/v1/auth/login    → login
 * - POST /api/v1/auth/refresh  → refreshToken
 * - POST /api/v1/auth/logout   → logout
 * - GET  /api/v1/auth/me       → getMe
 *
 * Auth tokens are managed via HttpOnly cookies (set by backend).
 */

import type { AuthResponse, AuthUser } from "@/types";
import { api, APIError } from "./api";

/**
 * Sign up a new user
 * @throws {APIError} If registration fails
 */
export async function signup(
  email: string,
  password: string,
  name: string
): Promise<AuthResponse> {
  try {
    const response = await api.post<AuthResponse>("/v1/auth/signup", {
      email,
      password,
      name,
    });
    return response;
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Signup failed:", error.message);
    }
    throw error;
  }
}

/**
 * Login with email and password
 * Backend sets HttpOnly cookies for access_token and refresh_token
 * @throws {APIError} If authentication fails
 */
export async function login(
  email: string,
  password: string
): Promise<AuthResponse> {
  try {
    const response = await api.post<AuthResponse>("/v1/auth/login", {
      email,
      password,
    });
    return response;
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Login failed:", error.message);
    }
    throw error;
  }
}

/**
 * Refresh the access token using the refresh_token cookie
 * @throws {APIError} If refresh fails
 */
export async function refreshToken(): Promise<{ success: boolean; data: { accessToken: string } }> {
  try {
    return await api.post("/v1/auth/refresh");
  } catch (error) {
    if (error instanceof APIError) {
      console.warn("Token refresh failed:", error.message);
    }
    throw error;
  }
}

/**
 * Logout — clears auth cookies on the backend
 */
export async function logout(): Promise<void> {
  try {
    await api.post("/v1/auth/logout");
  } catch (error) {
    // Logout should succeed even if session is already expired
    console.warn("Logout request failed:", error);
  }
}

/**
 * Get the current authenticated user's profile
 * @throws {APIError} If not authenticated
 */
export async function getMe(): Promise<AuthUser | null> {
  try {
    const response = await api.get<{ success: boolean; data: AuthUser }>("/v1/auth/me");
    return response?.data ?? null;
  } catch (error) {
    if (error instanceof APIError && error.status !== 401) {
      console.warn("Failed to fetch user profile:", error.message);
    }
    return null;
  }
}
