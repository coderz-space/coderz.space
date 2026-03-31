import type { Role } from "@/types";

/**
 * roleService.ts — Role management (local-only)
 *
 * The backend does not have role selection endpoints.
 * Role is stored in localStorage and used locally for routing.
 * TODO: Add backend endpoint for role persistence when available.
 */

let _selectedRole: Role | null = null;

/**
 * Select a role and persist to localStorage
 * (No backend endpoint exists — local storage only)
 */
export async function selectRole(role: Role): Promise<void> {
  // Validate role locally
  if (role !== "mentor" && role !== "mentee") {
    throw new Error("Invalid role");
  }

  _selectedRole = role;
  if (typeof window !== "undefined") {
    localStorage.setItem("coderz_selected_role", role);
  }
}

/**
 * Retrieve the user's selected role from cache or localStorage
 * @returns The selected role or null if not set
 */
export async function getSelectedRole(): Promise<Role | null> {
  if (_selectedRole) return _selectedRole;

  // Fallback to localStorage cache
  if (typeof window !== "undefined") {
    const stored = localStorage.getItem("coderz_selected_role") as Role | null;
    if (stored === "mentor" || stored === "mentee") {
      _selectedRole = stored;
      return stored;
    }
  }

  return null;
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
