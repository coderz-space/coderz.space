// ─── Role ────────────────────────────────────────────────────────────────────
export type Role = "mentor" | "mentee";

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

// ─── Mentee Request (mirrors future DB schema) ────────────────────────────────
// When backend is ready, replace localStorage calls with POST /api/mentee-requests
export type SheetId = "gfg-dsa-360" | "strivers-dsa-sheet";

export interface MenteeRequest {
  id: string;
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  passwordHash: string;
  signedUpAt: string;
  status: "pending" | "approved" | "rejected";
  assignedSheet?: SheetId;
  bio?: string;
  github?: string;
  linkedin?: string;
}

// ─── Mentor Profile ───────────────────────────────────────────────────────────
export interface MentorProfile {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  bio?: string;
  github?: string;
  linkedin?: string;
  joinedAt: string;
}

// ─── Question (mirrors future DB schema) ─────────────────────────────────────
// When backend is ready, replace dummy data with GET /api/mentees/:username/questions
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
