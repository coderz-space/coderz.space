/**
 * menteeService.ts — Mentee management service
 *
 * Auth functions (registerMentee, loginMentee) call the REAL backend.
 * All other functions are STUBS — they show a toast notification and
 * return mock data until backend endpoints are implemented.
 *
 * Real Backend Endpoints:
 * - POST /api/v1/auth/signup  → registerMentee
 * - POST /api/v1/auth/login   → loginMentee / loginMenteeByEmail
 *
 * Stubbed (no backend endpoint yet):
 * - GET /api/mentee-requests → getMenteeRequests
 * - PATCH /api/mentee-requests/:id → updateMenteeStatus
 * - DELETE /api/mentee-requests/:id → deleteMentee
 * - GET /api/mentees/:username/questions → getMenteeQuestions
 * - PATCH /api/mentees/:username/questions/:questionId → updateQuestionProgress/Details
 * - GET /api/mentees/:username/profile → getMenteeProfile
 * - GET /api/leaderboard → getLeaderboard
 * - GET /api/mentor/profile → getMentorProfile
 * - PATCH /api/mentor/profile → updateMentorProfile
 */

import type {
  MenteeRequest,
  Question,
  QuestionProgressStatus,
  MentorProfile,
  SheetId,
  AuthResponse,
} from "@/types";
import { api, APIError } from "./api";
import { showStubNotification } from "@/components/StubToast";

/**
 * Local cache for frequently accessed data to improve UX
 * These are fallbacks — primary source is always the backend
 */
const memoryCache = new Map<string, { data: any; timestamp: number }>();
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes

function getCacheKey(prefix: string, ...args: string[]): string {
  return `${prefix}:${args.join(':')}`;
}

function getCached<T>(key: string): T | null {
  const cached = memoryCache.get(key);
  if (!cached) return null;
  if (Date.now() - cached.timestamp > CACHE_TTL) {
    memoryCache.delete(key);
    return null;
  }
  return cached.data as T;
}

function setCache(key: string, data: any): void {
  memoryCache.set(key, { data, timestamp: Date.now() });
}

// ─── REAL API CALLS (connected to backend) ───────────────────────────────────

/**
 * Register a new mentee via POST /api/v1/auth/signup
 * @throws {APIError} If registration fails
 */
export async function registerMentee(
  data: Pick<MenteeRequest, "firstName" | "lastName" | "username" | "email" | "passwordHash">
): Promise<MenteeRequest> {
  try {
    const response = await api.post<AuthResponse>("/v1/auth/signup", {
      email: data.email,
      password: data.passwordHash, // backend hashes this
      name: `${data.firstName} ${data.lastName}`.trim(),
    });

    // Map backend response to MenteeRequest shape for UI compatibility
    return {
      id: response.data.user.id,
      firstName: data.firstName,
      lastName: data.lastName,
      username: data.username,
      email: data.email,
      passwordHash: "", // never store on client
      signedUpAt: new Date().toISOString(),
      status: "approved", // backend auto-approves on signup
    };
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Registration failed:", error.message);
    }
    throw error;
  }
}

/**
 * Login mentee with email (username login removed — backend only supports email)
 * Calls POST /api/v1/auth/login
 * @throws {APIError} If authentication fails
 */
export async function loginMentee(
  _usernameOrEmail: string,
  password: string
): Promise<{ token: string; refreshToken: string; mentee: MenteeRequest }> {
  // Backend only supports email login — treat the identifier as email
  return loginMenteeByEmail(_usernameOrEmail, password);
}

/**
 * Login mentee with email
 * Calls POST /api/v1/auth/login
 * @throws {APIError} If authentication fails
 */
export async function loginMenteeByEmail(
  email: string,
  password: string
): Promise<{ token: string; refreshToken: string; mentee: MenteeRequest }> {
  try {
    const response = await api.post<AuthResponse>("/v1/auth/login", {
      email,
      password,
    });

    // Tokens are set as HttpOnly cookies by backend — no localStorage needed
    // Map response to legacy shape for UI compatibility
    const nameParts = response.data.user.name.split(" ");
    return {
      token: response.data.accessToken,
      refreshToken: response.data.refreshToken,
      mentee: {
        id: response.data.user.id,
        firstName: nameParts[0] || "",
        lastName: nameParts.slice(1).join(" ") || "",
        username: response.data.user.email.split("@")[0], // derive username from email
        email: response.data.user.email,
        passwordHash: "",
        signedUpAt: new Date().toISOString(),
        status: "approved",
      },
    };
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Login failed:", error.message);
    }
    throw error;
  }
}

// ─── STUB FUNCTIONS (no backend endpoint yet) ────────────────────────────────

/**
 * Get all mentee requests (admin view)
 * TODO: Implement backend endpoint GET /api/mentee-requests
 */
export async function getMenteeRequests(): Promise<MenteeRequest[]> {
  showStubNotification("Get Mentee Requests");
  return [];
}

/**
 * Update mentee request status
 * TODO: Implement backend endpoint PATCH /api/mentee-requests/:id
 */
export async function updateMenteeStatus(
  id: string,
  status: "pending" | "approved" | "rejected",
  assignedSheet?: SheetId
): Promise<void> {
  showStubNotification("Update Mentee Status");
  // no-op stub
  void id; void status; void assignedSheet;
}

/**
 * Delete a mentee
 * TODO: Implement backend endpoint DELETE /api/mentee-requests/:id
 */
export async function deleteMentee(id: string): Promise<void> {
  showStubNotification("Delete Mentee");
  void id;
}

/**
 * Get questions for a mentee
 * TODO: Implement backend endpoint GET /api/mentees/:username/questions
 */
export async function getMenteeQuestions(username: string): Promise<Question[]> {
  showStubNotification("Get Mentee Questions");
  void username;
  return [];
}

/**
 * Update individual question progress status
 * TODO: Implement backend endpoint PATCH /api/mentees/:username/questions/:questionId
 */
export async function updateQuestionProgress(
  username: string,
  questionId: string,
  progressStatus: QuestionProgressStatus
): Promise<void> {
  showStubNotification("Update Question Progress");
  void username; void questionId; void progressStatus;
}

/**
 * Update question notes (solution & resources)
 * TODO: Implement backend endpoint PATCH /api/mentees/:username/questions/:questionId
 */
export async function updateQuestionDetails(
  username: string,
  questionId: string,
  details: { solution?: string; resources?: string }
): Promise<void> {
  showStubNotification("Update Question Details");
  void username; void questionId; void details;
}

/**
 * Get specific question detail for a mentee
 * TODO: Implement backend endpoint GET /api/mentees/:username/questions/:questionId
 */
export async function getQuestionDetail(
  username: string,
  questionId: string
): Promise<Question | null> {
  showStubNotification("Get Question Detail");
  void username; void questionId;
  return null;
}

/**
 * Assign a task to a mentee
 * TODO: Implement backend endpoint POST /api/mentees/:username/questions
 */
export async function assignTaskToMentee(
  username: string,
  task: { title: string; description: string; difficulty: Question["difficulty"]; topic: string }
): Promise<Question> {
  showStubNotification("Assign Task to Mentee");
  void username;
  return {
    id: crypto.randomUUID(),
    ...task,
    description: task.description,
    status: "pending",
    progressStatus: "not_started",
    assignedAt: new Date().toISOString(),
  };
}

/**
 * Get mentee's public profile
 * TODO: Implement backend endpoint GET /api/mentees/:username/profile
 */
export async function getMenteeProfile(profileUsername: string): Promise<{
  firstName: string;
  lastName: string;
  username: string;
  solved: number;
  joinedAt: string;
  bio?: string;
  github?: string;
  linkedin?: string;
} | null> {
  showStubNotification("Get Mentee Profile");
  void profileUsername;
  return null;
}

/**
 * Get leaderboard of top mentees
 * TODO: Implement backend endpoint GET /api/leaderboard
 */
export async function getLeaderboard(): Promise<
  Array<{ username: string; firstName: string; lastName: string; solved: number }>
> {
  showStubNotification("Get Leaderboard");
  return [];
}

/**
 * Get mentor profile
 * TODO: Implement backend endpoint GET /api/mentor/profile
 */
export async function getMentorProfile(): Promise<MentorProfile> {
  showStubNotification("Get Mentor Profile");
  // Fallback default profile
  return {
    firstName: "Mentor",
    lastName: "",
    username: "mentor",
    email: "",
    joinedAt: new Date().toISOString(),
  };
}

/**
 * Update mentor profile (async API call)
 * TODO: Implement backend endpoint PATCH /api/mentor/profile
 */
export async function updateMentorProfile(
  updates: Partial<Omit<MentorProfile, "joinedAt">>
): Promise<MentorProfile> {
  showStubNotification("Update Mentor Profile");
  return {
    firstName: updates.firstName ?? "Mentor",
    lastName: updates.lastName ?? "",
    username: updates.username ?? "mentor",
    email: updates.email ?? "",
    joinedAt: new Date().toISOString(),
    ...updates,
  };
}

/**
 * Save mentor profile (backwards compatible wrapper for existing code)
 * Calls updateMentorProfile with all profile fields
 */
export async function saveMentorProfile(
  profile: Partial<Omit<MentorProfile, "joinedAt">>
): Promise<MentorProfile> {
  return updateMentorProfile(profile);
}

/**
 * Update mentee profile fields
 * TODO: Implement backend endpoint PATCH /api/mentees/:username/profile
 */
export async function updateMenteeProfile(
  username: string,
  fields: { bio?: string; github?: string; linkedin?: string }
): Promise<void> {
  showStubNotification("Update Mentee Profile");
  void username; void fields;
}

/**
 * Update mentee password
 * TODO: Implement backend endpoint PATCH /api/mentees/:username/password
 */
export async function updateMenteePassword(
  username: string,
  currentPassword: string,
  newPassword: string
): Promise<{ ok: boolean; error?: string }> {
  showStubNotification("Update Mentee Password");
  void username; void currentPassword; void newPassword;
  return { ok: false, error: "Backend not implemented yet" };
}

/**
 * Update mentor password
 * TODO: Implement backend endpoint PATCH /api/mentor/password
 */
export async function updateMentorPassword(
  currentPassword: string,
  newPassword: string
): Promise<{ ok: boolean; error?: string }> {
  showStubNotification("Update Mentor Password");
  void currentPassword; void newPassword;
  return { ok: false, error: "Backend not implemented yet" };
}

/**
 * Clear all caches (call on logout)
 */
export function clearCaches(): void {
  memoryCache.clear();
}

// ─── STUB: Get assigned tasks for a mentee (with progress applied) ───────────
// TODO: Replace with real API call: GET /api/mentees/:username/assigned-tasks
export function getAssignedTasks(username: string): Question[] {
  showStubNotification("Get Assigned Tasks");
  void username;
  return [];
}

export default {
  registerMentee,
  getMenteeRequests,
  updateMenteeStatus,
  loginMentee,
  loginMenteeByEmail,
  deleteMentee,
  getMenteeQuestions,
  updateQuestionProgress,
  updateQuestionDetails,
  getQuestionDetail,
  assignTaskToMentee,
  getMenteeProfile,
  getLeaderboard,
  getMentorProfile,
  updateMentorProfile,
  updateMenteeProfile,
  updateMenteePassword,
  updateMentorPassword,
  clearCaches,
  getAssignedTasks,
};
