export type Role = "mentor" | "mentee";
export type SheetId = "gfg-dsa-360" | "strivers-dsa-sheet";
export type QuestionProgressStatus =
  | "not_started"
  | "discussion_needed"
  | "revision_needed"
  | "completed";

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

export interface AppUser {
  id: string;
  name: string;
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  bio?: string;
  github?: string;
  linkedin?: string;
}

export interface AppContext {
  role: Role | "unknown";
  accountStatus: string;
  user: AppUser;
  organization?: {
    id: string;
    name: string;
    slug: string;
  };
  bootcamp?: {
    id: string;
    name: string;
  };
  enrollment?: {
    id?: string;
    assignedSheet?: SheetId;
  };
}

export interface AppState {
  showRoleCard: boolean;
  selectedRole: Role | null;
}

export interface HeroSectionProps {
  onGetStarted: () => void;
}

export interface RoleCardProps {
  onSelectRole: (role: Role) => void;
  onClose: () => void;
}

export interface MenteeRequest {
  id: string;
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  signedUpAt: string;
  status: "pending" | "approved" | "rejected";
  assignedSheet?: SheetId;
  passwordHash?: string;
  bio?: string;
  github?: string;
  linkedin?: string;
}

export interface SheetQuestion {
  id: string;
  title: string;
  topic: string;
  difficulty: "easy" | "medium" | "hard";
}

export interface Sheet {
  key: SheetId;
  name: string;
  questions: SheetQuestion[];
}

export interface DayAssignmentMentee {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  assigned: boolean;
  assignedSheet?: SheetId;
}

export interface DayAssignments {
  day: string;
  mentees: DayAssignmentMentee[];
}

export interface Question {
  id: string;
  title: string;
  description: string;
  difficulty: "easy" | "medium" | "hard";
  topic: string;
  status: "pending" | "completed";
  progressStatus: QuestionProgressStatus;
  assignedAt: string;
  completedAt?: string;
  solution?: string;
  resources?: string;
}

export interface Profile {
  firstName: string;
  lastName: string;
  username: string;
  email?: string;
  solved: number;
  joinedAt: string;
  bio?: string;
  github?: string;
  linkedin?: string;
}

export interface MentorProfile extends Profile {
  email: string;
}

export interface LeaderboardEntry {
  username: string;
  firstName: string;
  lastName: string;
  solved: number;
}

export interface MenteeSignupInput {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  password: string;
}

export interface LoginResult {
  auth: AuthResponseData;
  context: AppContext;
}
