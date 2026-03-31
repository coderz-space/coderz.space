import type { Role } from "@/types";
import { api, APIError } from "./api";

/**
 * roleService.ts — Role management with API integration
 *
 * Provides secure role selection and retrieval following these practices:
 * - Server-side validation of role selection
 * - Secure session token storage
 * - Graceful error handling with fallback to localStorage
 */

let _selectedRole: Role | null = null;

/**
 * Select a role and persist to backend + localStorage
 * @throws {APIError} If API request fails
 */
export async function selectRole(role: Role): Promise<void> {
  try {
    // Validate role locally first (defense in depth)
    if (role !== "mentor" && role !== "mentee") {
      throw new Error("Invalid role");
    }

    // Send role selection to backend for persistent storage
    await api.post("/auth/select-role", { role });

    // Store locally as cache
    _selectedRole = role;
    if (typeof window !== "undefined") {
      localStorage.setItem("coderz_selected_role", role);
    }
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Failed to select role:", error.message);
    }
    throw error;
  }
}

/**
 * Retrieve the user's selected role from backend or cache
 * @returns The selected role or null if not set
 */
export async function getSelectedRole(): Promise<Role | null> {
  try {
    // Try to fetch from backend first
    if (typeof window !== "undefined") {
      const response = await api.get<{ role: Role }>("/auth/get-role");
      if (response?.role) {
        _selectedRole = response.role;
        localStorage.setItem("coderz_selected_role", response.role);
        return response.role;
      }
    }
  } catch (error) {
    if (error instanceof APIError && error.status !== 401) {
      console.warn("Failed to fetch role from backend, using cache:", error.message);
    }
  }

  // Fallback to localStorage cache
  if (typeof window !== "undefined") {
    const stored = localStorage.getItem("coderz_selected_role") as Role | null;
    if (stored === "mentor" || stored === "mentee") {
      _selectedRole = stored;
      return stored;
    }
  }

  return _selectedRole;
}

/**
 * Clear selected role on logout
 */
export function clearSelectedRole(): void {
  _selectedRole = null;
  if (typeof window !== "undefined") {
    localStorage.removeItem("coderz_selected_role");
  }
}

