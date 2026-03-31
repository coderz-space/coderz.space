import type { Role } from "@/types";

// In-memory store — swap for fetch() calls when a backend is ready
let _selectedRole: Role | null = null;

export async function selectRole(role: Role): Promise<void> {
  _selectedRole = role;
  if (typeof window !== "undefined") {
    localStorage.setItem("coderz_selected_role", role);
  }
}

export async function getSelectedRole(): Promise<Role | null> {
  if (typeof window !== "undefined") {
    const stored = localStorage.getItem("coderz_selected_role") as Role | null;
    if (stored === "mentor" || stored === "mentee") {
      _selectedRole = stored;
    }
  }
  return _selectedRole;
}
