/**
 * menteeService.ts — localStorage stub
 *
 * BACKEND INTEGRATION GUIDE:
 * Each function maps 1-to-1 to a REST endpoint. When the backend is ready:
 *   1. Replace the localStorage read/write with a fetch() call to the endpoint shown.
 *   2. Keep the function signature identical — no component changes needed.
 *
 * Storage key: "coderz_mentee_requests"  →  DB table: mentee_requests
 */

import type { MenteeRequest, Question, QuestionProgressStatus } from "@/types";

// ─── Dummy questions assigned to every approved mentee ───────────────────────
// Replace with: GET /api/mentees/:username/questions
const DUMMY_QUESTIONS: Question[] = [
  {
    id: "q1",
    title: "Two Sum",
    description: "Given an array of integers, return indices of the two numbers that add up to a target.",
    difficulty: "easy",
    topic: "Arrays",
    status: "pending",
    progressStatus: "not_started",
    assignedAt: new Date().toISOString(),
  },
  {
    id: "q2",
    title: "Valid Parentheses",
    description: "Given a string containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.",
    difficulty: "easy",
    topic: "Stack",
    status: "pending",
    progressStatus: "not_started",
    assignedAt: new Date().toISOString(),
  },
  {
    id: "q3",
    title: "Merge Two Sorted Lists",
    description: "Merge two sorted linked lists and return it as a new sorted list.",
    difficulty: "easy",
    topic: "Linked List",
    status: "completed",
    progressStatus: "completed",
    assignedAt: new Date(Date.now() - 86400000 * 3).toISOString(),
    completedAt: new Date(Date.now() - 86400000).toISOString(),
  },
  {
    id: "q4",
    title: "Binary Search",
    description: "Given a sorted array of integers, implement binary search.",
    difficulty: "easy",
    topic: "Binary Search",
    status: "completed",
    progressStatus: "completed",
    assignedAt: new Date(Date.now() - 86400000 * 5).toISOString(),
    completedAt: new Date(Date.now() - 86400000 * 2).toISOString(),
  },
  {
    id: "q5",
    title: "Longest Substring Without Repeating Characters",
    description: "Find the length of the longest substring without repeating characters.",
    difficulty: "medium",
    topic: "Sliding Window",
    status: "pending",
    progressStatus: "not_started",
    assignedAt: new Date().toISOString(),
  },
];

// ─── Read all mentee requests ─────────────────────────────────────────────────
// Replace with: GET /api/mentee-requests
export function getMenteeRequests(): MenteeRequest[] {
  if (typeof window === "undefined") return [];
  return JSON.parse(localStorage.getItem("coderz_mentee_requests") || "[]");
}

// ─── Save all mentee requests ─────────────────────────────────────────────────
// Replace with: PUT /api/mentee-requests (bulk) or PATCH /api/mentee-requests/:id
function saveMenteeRequests(requests: MenteeRequest[]): void {
  localStorage.setItem("coderz_mentee_requests", JSON.stringify(requests));
}

// ─── Register a new mentee ────────────────────────────────────────────────────
// Replace with: POST /api/mentee-requests
export function registerMentee(
  data: Pick<MenteeRequest, "firstName" | "lastName" | "username" | "email" | "passwordHash">
): MenteeRequest {
  const existing = getMenteeRequests();
  const newRequest: MenteeRequest = {
    id: crypto.randomUUID(),
    ...data,
    signedUpAt: new Date().toISOString(),
    status: "pending",
  };
  saveMenteeRequests([...existing, newRequest]);
  return newRequest;
}

// ─── Update mentee request status ────────────────────────────────────────────
// Replace with: PATCH /api/mentee-requests/:id  { status, assignedSheet? }
export function updateMenteeStatus(
  id: string,
  status: "approved" | "rejected",
  assignedSheet?: import("@/types").SheetId
): void {
  const requests = getMenteeRequests();
  const updated = requests.map((r) =>
    r.id === id ? { ...r, status, ...(assignedSheet ? { assignedSheet } : {}) } : r
  );
  saveMenteeRequests(updated);
}

// ─── Find approved mentee by username + password ─────────────────────────────
// Replace with: POST /api/auth/mentee/login  { username, password }
export function loginMentee(
  username: string,
  password: string
): MenteeRequest | null {
  const requests = getMenteeRequests();
  const match = requests.find(
    (r) =>
      r.username === username &&
      r.passwordHash === password && // dev stub: plain text; use bcrypt in real backend
      r.status === "approved"
  );
  return match ?? null;
}

// ─── Find approved mentee by email + password ────────────────────────────────
// Replace with: POST /api/auth/mentee/login  { email, password }
export function loginMenteeByEmail(
  email: string,
  password: string
): MenteeRequest | null {
  const requests = getMenteeRequests();
  const match = requests.find(
    (r) =>
      r.email.toLowerCase() === email.toLowerCase() &&
      r.passwordHash === password && // dev stub: plain text; use bcrypt in real backend
      r.status === "approved"
  );
  return match ?? null;
}

// ─── Delete a mentee entirely ─────────────────────────────────────────────────
// Replace with: DELETE /api/mentee-requests/:id
export function deleteMentee(id: string): void {
  const requests = getMenteeRequests();
  saveMenteeRequests(requests.filter((r) => r.id !== id));
}

// ─── Get/save question detail notes per mentee ───────────────────────────────
// Storage key: "coderz_question_details_{username}"  →  DB table: question_details
// Replace with: GET/PATCH /api/mentees/:username/questions/:questionId/details
type DetailMap = Record<string, { solution?: string; resources?: string }>;

function getDetailMap(username: string): DetailMap {
  if (typeof window === "undefined") return {};
  return JSON.parse(localStorage.getItem(`coderz_question_details_${username}`) || "{}");
}

// Replace with: PATCH /api/mentees/:username/questions/:questionId  { solution, resources }
export function updateQuestionDetails(
  username: string,
  questionId: string,
  details: { solution?: string; resources?: string }
): void {
  const map = getDetailMap(username);
  map[questionId] = { ...map[questionId], ...details };
  localStorage.setItem(`coderz_question_details_${username}`, JSON.stringify(map));
}

// Replace with: GET /api/mentees/:username/questions/:questionId
export function getQuestionDetail(username: string, questionId: string): import("@/types").Question | null {
  const questions = getMenteeQuestions(username);
  const q = questions.find((q) => q.id === questionId) ?? null;
  if (!q) return null;
  const details = getDetailMap(username);
  return { ...q, ...details[questionId] };
}
// Storage key: "coderz_question_progress_{username}"  →  DB table: question_progress
// Replace with: GET/PATCH /api/mentees/:username/questions/:questionId/progress
type ProgressMap = Record<string, { progressStatus: QuestionProgressStatus; completedAt?: string }>;

function getProgressMap(username: string): ProgressMap {
  if (typeof window === "undefined") return {};
  return JSON.parse(localStorage.getItem(`coderz_question_progress_${username}`) || "{}");
}

// ─── Update a single question's progress status ───────────────────────────────
// Replace with: PATCH /api/mentees/:username/questions/:questionId  { progressStatus }
export function updateQuestionProgress(
  username: string,
  questionId: string,
  progressStatus: QuestionProgressStatus
): void {
  const map = getProgressMap(username);
  map[questionId] = {
    progressStatus,
    completedAt:
      progressStatus === "completed" || progressStatus === "revision_needed"
        ? (map[questionId]?.completedAt ?? new Date().toISOString())
        : undefined,
  };
  localStorage.setItem(`coderz_question_progress_${username}`, JSON.stringify(map));
}

// ─── Assign a task to a mentee ────────────────────────────────────────────────
// Replace with: POST /api/mentees/:username/tasks  { taskId, title, description, difficulty, topic }
export function assignTaskToMentee(
  username: string,
  task: { title: string; description: string; difficulty: Question["difficulty"]; topic: string }
): void {
  const key = `coderz_assigned_tasks_${username}`;
  const existing: Question[] = JSON.parse(
    (typeof window !== "undefined" && localStorage.getItem(key)) || "[]"
  );
  const newTask: Question = {
    id: crypto.randomUUID(),
    ...task,
    status: "pending",
    progressStatus: "not_started",
    assignedAt: new Date().toISOString(),
  };
  localStorage.setItem(key, JSON.stringify([...existing, newTask]));
}

// ─── Get public profile for a mentee ─────────────────────────────────────────
// Replace with: GET /api/mentees/:profileUsername/profile
export function getMenteeProfile(profileUsername: string): {
  firstName: string; lastName: string; username: string; solved: number; joinedAt: string;
  bio?: string; github?: string; linkedin?: string;
} | null {
  const mentee = getMenteeRequests().find(
    (r) => r.username === profileUsername && r.status === "approved"
  );
  if (!mentee) return null;
  const progress: ProgressMap = JSON.parse(
    (typeof window !== "undefined" && localStorage.getItem(`coderz_question_progress_${profileUsername}`)) || "{}"
  );
  const solved = Object.values(progress).filter(
    (p) => p.progressStatus === "completed" || p.progressStatus === "revision_needed"
  ).length;
  return {
    firstName: mentee.firstName, lastName: mentee.lastName, username: mentee.username,
    solved, joinedAt: mentee.signedUpAt,
    bio: mentee.bio, github: mentee.github, linkedin: mentee.linkedin,
  };
}

// ─── Leaderboard: all approved mentees ranked by solved question count ────────
// Replace with: GET /api/leaderboard
export function getLeaderboard(): { username: string; firstName: string; lastName: string; solved: number }[] {
  const approved = getMenteeRequests().filter((r) => r.status === "approved");
  return approved
    .map((r) => {
      const progress: ProgressMap = JSON.parse(
        (typeof window !== "undefined" && localStorage.getItem(`coderz_question_progress_${r.username}`)) || "{}"
      );
      const solved = Object.values(progress).filter(
        (p) => p.progressStatus === "completed" || p.progressStatus === "revision_needed"
      ).length;
      return { username: r.username, firstName: r.firstName, lastName: r.lastName, solved };
    })
    .sort((a, b) => b.solved - a.solved);
}

// ─── Get questions for a mentee (with persisted progress applied) ─────────────
// Replace with: GET /api/mentees/:username/questions
export function getMenteeQuestions(username: string): Question[] {
  const map = getProgressMap(username);
  return DUMMY_QUESTIONS.map((q) => {
    const override = map[q.id];
    if (!override) return q;
    const inCompletedBucket =
      override.progressStatus === "completed" ||
      override.progressStatus === "revision_needed";
    return {
      ...q,
      progressStatus: override.progressStatus,
      status: inCompletedBucket ? "completed" : "pending",
      completedAt: inCompletedBucket ? override.completedAt : undefined,
    };
  });
}

// ─── Mentor profile (localStorage stub) ──────────────────────────────────────
// Replace with: GET/PATCH /api/mentor/profile
import type { MentorProfile } from "@/types";

export function getMentorProfile(): MentorProfile {
  if (typeof window === "undefined") {
    return { firstName: "Mentor", lastName: "", username: "mentor", email: "", joinedAt: new Date().toISOString() };
  }
  const stored = localStorage.getItem("coderz_mentor_profile");
  if (stored) return JSON.parse(stored);
  return { firstName: "Mentor", lastName: "", username: "mentor", email: "", joinedAt: new Date().toISOString() };
}

export function saveMentorProfile(profile: MentorProfile): void {
  localStorage.setItem("coderz_mentor_profile", JSON.stringify(profile));
}

// ─── Update mentee profile fields ────────────────────────────────────────────
// Replace with: PATCH /api/mentees/:username/profile
export function updateMenteeProfile(
  username: string,
  fields: { bio?: string; github?: string; linkedin?: string }
): void {
  const requests = getMenteeRequests();
  const updated = requests.map((r) =>
    r.username === username ? { ...r, ...fields } : r
  );
  saveMenteeRequests(updated);
}

// ─── Update mentee password ───────────────────────────────────────────────────
// Replace with: PATCH /api/mentees/:username/password  { currentPassword, newPassword }
export function updateMenteePassword(
  username: string,
  currentPassword: string,
  newPassword: string
): { ok: boolean; error?: string } {
  const requests = getMenteeRequests();
  const mentee = requests.find((r) => r.username === username);
  if (!mentee) return { ok: false, error: "User not found." };
  if (mentee.passwordHash !== currentPassword) return { ok: false, error: "Current password is incorrect." };
  const updated = requests.map((r) =>
    r.username === username ? { ...r, passwordHash: newPassword } : r
  );
  saveMenteeRequests(updated);
  return { ok: true };
}

// ─── Update mentor password ───────────────────────────────────────────────────
// Replace with: PATCH /api/mentor/password  { currentPassword, newPassword }
export function updateMentorPassword(
  currentPassword: string,
  newPassword: string
): { ok: boolean; error?: string } {
  const stored = typeof window !== "undefined"
    ? localStorage.getItem("coderz_mentor_password")
    : null;
  const current = stored ?? "mentor123"; // default dev password
  if (current !== currentPassword) return { ok: false, error: "Current password is incorrect." };
  localStorage.setItem("coderz_mentor_password", newPassword);
  return { ok: true };
}

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
