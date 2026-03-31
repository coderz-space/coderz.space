// ─── Role ────────────────────────────────────────────────────────────────────
export type Role = "mentor" | "mentee";

// ─── Backend Auth DTOs (matches /api/v1/auth/*) ──────────────────────────────
export interface AuthUser {
  id: string;
  name: string;
  email: string;
  emailVerified: boolean;
}

export interface AuthResponseData {
  accessToken: string;
  refreshToken: string;
  user: AuthUser;
}

export interface AuthResponse {
  success: boolean;
  data: AuthResponseData;
}

// ─── Backend Organization DTOs (matches /api/v1/organizations/*) ─────────────
export interface OrganizationData {
  id: string;
  name: string;
  slug: string;
  description: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface MemberData {
  id: string;
  organizationId: string;
  userId: string;
  role: string;
  joinedAt: string;
  name?: string;
  email?: string;
  avatarUrl?: string;
}

// ─── App UI State ─────────────────────────────────────────────────────────────
export interface AppState {
  showRoleCard: boolean;
  selectedRole: Role | null;
}

// ─── Component Props ──────────────────────────────────────────────────────────
export interface HeroSectionProps {
  onGetStarted: () => void;
}

export interface RoleCardProps {
  onSelectRole: (role: Role) => void;
  onClose: () => void;
}

// ─── Mentee Request (UI type — TODO: needs backend endpoints) ─────────────────
export type SheetId = "gfg-dsa-360" | "strivers-dsa-sheet";

export interface MenteeRequest {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  passwordHash: string;
  signedUpAt: string;
  status: "pending" | "approved" | "rejected";
  assignedSheet?: SheetId;
}

// ─── Mentor Profile ───────────────────────────────────────────────────────────
export interface MentorProfile {
  firstName: string;
  lastName: string;
  email: string;
  joinedAt: string;
}

// ─── Question (UI type — TODO: needs backend endpoints) ──────────────────────
export type QuestionProgressStatus = "not_started" | "discussion_needed" | "revision_needed" | "completed";

export interface Question {
  id: string;
  title: string;
  description: string;
  difficulty: "easy" | "medium" | "hard";
  topic: string;
  status: "pending" | "completed";          // overall bucket (pending/completed section)
  progressStatus: QuestionProgressStatus;   // mentee's self-reported progress
  assignedAt: string;       // ISO 8601
  completedAt?: string;     // ISO 8601, present only when status === "completed"
  solutionUrl?: string;     // link to solution submission
  solution?: string;        // mentee's written solution notes
  resources?: string;       // links / references used
}
