/**
 * menteeService.ts — Mentee management with API integration
 *
 * API Endpoints mapping:
 * - POST /api/auth/mentee-register → registerMentee
 * - GET /api/mentee-requests → getMenteeRequests
 * - PATCH /api/mentee-requests/:id → updateMenteeStatus
 * - POST /api/auth/mentee/login → loginMentee
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
  SheetId 
} from "@/types";
import { api, APIError } from "./api";

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

/**
 * Register a new mentee
 * @throws {APIError} If registration fails
 */
export async function registerMentee(
  data: Pick<MenteeRequest, "firstName" | "lastName" | "username" | "email" | "passwordHash">
): Promise<MenteeRequest> {
  try {
    const newRequest = await api.post<MenteeRequest>("/auth/mentee-register", data);
    return newRequest;
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Registration failed:", error.message);
    }
    throw error;
  }
}

/**
 * Get all mentee requests (admin view)
 * @throws {APIError} If fetch fails
 */
export async function getMenteeRequests(): Promise<MenteeRequest[]> {
  const cacheKey = getCacheKey("mentee:requests");
  const cached = getCached<MenteeRequest[]>(cacheKey);
  if (cached) return cached;

  try {
    const requests = await api.get<MenteeRequest[]>("/mentee-requests");
    setCache(cacheKey, requests);
    return requests;
  } catch (error) {
    if (error instanceof APIError) {
      console.warn("Failed to fetch mentee requests:", error.message);
    }
    return [];
  }
}

/**
 * Update mentee request status
 * @throws {APIError} If update fails
 */
export async function updateMenteeStatus(
  id: string,
  status: "pending" | "approved" | "rejected",
  assignedSheet?: SheetId
): Promise<void> {
  try {
    await api.patch(`/mentee-requests/${id}`, { 
      status, 
      ...(assignedSheet && { assignedSheet })
    });
    // Invalidate cache
    memoryCache.delete(getCacheKey("mentee:requests"));
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Failed to update mentee status:", error.message);
    }
    throw error;
  }
}

/**
 * Login mentee with username
 * @throws {APIError} If authentication fails
 */
export async function loginMentee(
  username: string,
  password: string
): Promise<{ token: string; refreshToken: string; mentee: MenteeRequest }> {
  try {
    const response = await api.post<{ token: string; refreshToken: string; mentee: MenteeRequest }>(
      "/auth/mentee/login",
      { username, password }
    );

    // Store tokens securely
    if (response.token) {
      localStorage.setItem("auth_token", response.token);
      localStorage.setItem("refresh_token", response.refreshToken);
    }

    return response;
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Login failed:", error.message);
    }
    throw error;
  }
}

/**
 * Login mentee with email
 * @throws {APIError} If authentication fails
 */
export async function loginMenteeByEmail(
  email: string,
  password: string
): Promise<{ token: string; refreshToken: string; mentee: MenteeRequest }> {
  try {
    const response = await api.post<{ token: string; refreshToken: string; mentee: MenteeRequest }>(
      "/auth/mentee/login",
      { email, password }
    );

    if (response.token) {
      localStorage.setItem("auth_token", response.token);
      localStorage.setItem("refresh_token", response.refreshToken);
    }

    return response;
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Email login failed:", error.message);
    }
    throw error;
  }
}

/**
 * Delete a mentee
 * @throws {APIError} If deletion fails
 */
export async function deleteMentee(id: string): Promise<void> {
  try {
    await api.delete(`/mentee-requests/${id}`);
    memoryCache.delete(getCacheKey("mentee:requests"));
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Failed to delete mentee:", error.message);
    }
    throw error;
  }
}


/**
 * Get questions for a mentee
 * @throws {APIError} If fetch fails
 */
export async function getMenteeQuestions(username: string): Promise<Question[]> {
  const cacheKey = getCacheKey("mentee:questions", username);
  const cached = getCached<Question[]>(cacheKey);
  if (cached) return cached;

  try {
    const questions = await api.get<Question[]>(`/mentees/${username}/questions`);
    setCache(cacheKey, questions);
    return questions;
  } catch (error) {
    if (error instanceof APIError && error.status !== 404) {
      console.warn("Failed to fetch mentee questions:", error.message);
    }
    return [];
  }
}

/**
 * Update individual question progress status
 * @throws {APIError} If update fails
 */
export async function updateQuestionProgress(
  username: string,
  questionId: string,
  progressStatus: QuestionProgressStatus
): Promise<void> {
  try {
    await api.patch(`/mentees/${username}/questions/${questionId}`, { progressStatus });
    // Invalidate cache
    memoryCache.delete(getCacheKey("mentee:questions", username));
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Failed to update question progress:", error.message);
    }
    throw error;
  }
}

/**
 * Update question notes (solution & resources)
 * @throws {APIError} If update fails
 */
export async function updateQuestionDetails(
  username: string,
  questionId: string,
  details: { solution?: string; resources?: string }
): Promise<void> {
  try {
    await api.patch(`/mentees/${username}/questions/${questionId}`, details);
    // Invalidate cache
    memoryCache.delete(getCacheKey("mentee:questions", username));
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Failed to update question details:", error.message);
    }
    throw error;
  }
}

/**
 * Get specific question detail for a mentee
 * @throws {APIError} If fetch fails
 */
export async function getQuestionDetail(
  username: string,
  questionId: string
): Promise<Question | null> {
  try {
    const question = await api.get<Question>(
      `/mentees/${username}/questions/${questionId}`
    );
    return question || null;
  } catch (error) {
    if (error instanceof APIError && error.status !== 404) {
      console.warn("Failed to fetch question detail:", error.message);
    }
    return null;
  }
}

/**
 * Assign a task to a mentee
 * @throws {APIError} If assignment fails
 */
export async function assignTaskToMentee(
  username: string,
  task: { title: string; description: string; difficulty: Question["difficulty"]; topic: string }
): Promise<Question> {
  try {
    const newTask = await api.post<Question>(
      `/mentees/${username}/questions`,
      task
    );
    // Invalidate cache
    memoryCache.delete(getCacheKey("mentee:questions", username));
    return newTask;
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Failed to assign task:", error.message);
    }
    throw error;
  }
}


/**
 * Get mentee's public profile
 * @throws {APIError} If fetch fails
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
  const cacheKey = getCacheKey("mentee:profile", profileUsername);
  const cached = getCached(cacheKey);
  if (cached) return cached;

  try {
    const profile = await api.get(
      `/mentees/${profileUsername}/profile`
    );
    setCache(cacheKey, profile);
    return profile || null;
  } catch (error) {
    if (error instanceof APIError && error.status !== 404) {
      console.warn("Failed to fetch mentee profile:", error.message);
    }
    return null;
  }
}

/**
 * Get leaderboard of top mentees
 * @throws {APIError} If fetch fails
 */
export async function getLeaderboard(): Promise<
  Array<{ username: string; firstName: string; lastName: string; solved: number }>
> {
  const cacheKey = getCacheKey("leaderboard");
  const cached = getCached(cacheKey);
  if (cached) return cached;

  try {
    const leaderboard = await api.get(
      "/leaderboard"
    );
    setCache(cacheKey, leaderboard);
    return leaderboard || [];
  } catch (error) {
    if (error instanceof APIError) {
      console.warn("Failed to fetch leaderboard:", error.message);
    }
    return [];
  }
}

/**
 * Get mentor profile
 * @throws {APIError} If fetch fails
 */
export async function getMentorProfile(): Promise<MentorProfile> {
  const cacheKey = getCacheKey("mentor:profile");
  const cached = getCached<MentorProfile>(cacheKey);
  if (cached) return cached;

  try {
    const profile = await api.get<MentorProfile>("/mentor/profile");
    if (profile) {
      setCache(cacheKey, profile);
      return profile;
    }
  } catch (error) {
    if (error instanceof APIError && error.status !== 401) {
      console.warn("Failed to fetch mentor profile:", error.message);
    }
  }

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
 * @throws {APIError} If update fails
 */
export async function updateMentorProfile(
  updates: Partial<Omit<MentorProfile, "joinedAt">>
): Promise<MentorProfile> {
  try {
    const updated = await api.patch<MentorProfile>("/mentor/profile", updates);
    // Invalidate cache
    memoryCache.delete(getCacheKey("mentor:profile"));
    return updated;
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Failed to update mentor profile:", error.message);
    }
    throw error;
  }
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
 * @throws {APIError} If update fails
 */
export async function updateMenteeProfile(
  username: string,
  fields: { bio?: string; github?: string; linkedin?: string }
): Promise<void> {
  try {
    await api.patch(`/mentees/${username}/profile`, fields);
    // Invalidate cache
    memoryCache.delete(getCacheKey("mentee:profile", username));
  } catch (error) {
    if (error instanceof APIError) {
      console.error("Failed to update mentee profile:", error.message);
    }
    throw error;
  }
}

/**
 * Update mentee password
 * @throws {APIError} If update fails
 */
export async function updateMenteePassword(
  username: string,
  currentPassword: string,
  newPassword: string
): Promise<{ ok: boolean; error?: string }> {
  try {
    const result = await api.patch<{ ok: boolean; error?: string }>(
      `/mentees/${username}/password`,
      { currentPassword, newPassword }
    );
    return result;
  } catch (error) {
    if (error instanceof APIError) {
      return { ok: false, error: error.message };
    }
    return { ok: false, error: "Password update failed" };
  }
}

/**
 * Update mentor password
 * @throws {APIError} If update fails
 */
export async function updateMentorPassword(
  currentPassword: string,
  newPassword: string
): Promise<{ ok: boolean; error?: string }> {
  try {
    const result = await api.patch<{ ok: boolean; error?: string }>(
      "/mentor/password",
      { currentPassword, newPassword }
    );
    return result;
  } catch (error) {
    if (error instanceof APIError) {
      return { ok: false, error: error.message };
    }
    return { ok: false, error: "Password update failed" };
  }
}

/**
 * Clear all caches (call on logout)
 */
export function clearCaches(): void {
  memoryCache.clear();
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
};


// ─── Get assigned tasks for a mentee (with progress applied) ─────────────────
// Replace with: GET /api/mentees/:username/assigned-tasks
export function getAssignedTasks(username: string): Question[] {
  if (typeof window === "undefined") return [];
  const tasks: Question[] = JSON.parse(localStorage.getItem(`coderz_assigned_tasks_${username}`) || "[]");
  const progressMap: Record<string, { progressStatus: QuestionProgressStatus; completedAt?: string }> = JSON.parse(
    localStorage.getItem(`coderz_question_progress_${username}`) || "{}"
  );
  return tasks.map((t) => {
    const p = progressMap[t.id];
    if (!p) return t;
    const inCompleted = p.progressStatus === "completed" || p.progressStatus === "revision_needed";
    return {
      ...t,
      progressStatus: p.progressStatus,
      status: inCompleted ? "completed" : "pending",
      completedAt: inCompleted ? p.completedAt : undefined,
    };
  });
}
