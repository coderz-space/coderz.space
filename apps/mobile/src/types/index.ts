// ─── Primitives ───────────────────────────────────────────────

export type UserRole = 'mentor' | 'mentee' | 'admin' | 'super_admin';
export type OrgRole = 'admin' | 'mentor' | 'mentee';
export type BootcampRole = 'mentor' | 'mentee';
export type Difficulty = 'easy' | 'medium' | 'hard';

// Maps exactly to backend assignment_problems.status
export type ProblemStatus = 'pending' | 'attempted' | 'completed';

// Mentee-facing label (maps to problem status + doubt flags)
export type MenteeStatus =
  | 'not_started'
  | 'discussion_needed'
  | 'revision_needed'
  | 'completed';

export type AssignmentStatus = 'active' | 'completed' | 'expired';

// ─── Auth ─────────────────────────────────────────────────────

export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;            // platform-level role
  avatarUrl?: string;
  emailVerified: boolean;
  createdAt: string;
}

export interface AuthTokens {
  accessToken: string;       // JWT, short-lived
  // refreshToken managed via HttpOnly cookie by backend
}

export interface LoginPayload {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  accessToken: string;
  orgRole: OrgRole;          // role inside the active org context
  bootcampRole?: BootcampRole;
  activeOrgId?: string;
  activeBootcampId?: string;
}

// ─── Organization ─────────────────────────────────────────────

export interface Organization {
  id: string;
  name: string;
  slug: string;
  description?: string;
  createdAt: string;
}

export interface OrgMember {
  id: string;               // organization_member.id
  userId: string;
  organizationId: string;
  role: OrgRole;
  joinedAt: string;
  user: Pick<User, 'id' | 'name' | 'email' | 'avatarUrl'>;
}

// ─── Bootcamp ─────────────────────────────────────────────────

export interface Bootcamp {
  id: string;
  organizationId: string;
  name: string;
  description?: string;
  startDate?: string;
  endDate?: string;
  isActive: boolean;
  createdAt: string;
}

export interface BootcampEnrollment {
  id: string;               // bootcamp_enrollment.id
  bootcampId: string;
  organizationMemberId: string;
  role: BootcampRole;
  enrolledAt: string;
}

// ─── Problem & Tags ───────────────────────────────────────────

export interface Tag {
  id: string;
  name: string;
  organizationId: string;
}

export interface ProblemResource {
  id: string;
  problemId: string;
  title: string;
  url: string;
}

export interface Problem {
  id: string;
  organizationId: string;
  createdBy: string;         // organization_member.id
  title: string;
  description: string;
  difficulty: Difficulty;
  externalLink?: string;
  tags: Tag[];
  resources: ProblemResource[];
  createdAt: string;
}

// ─── Assignments ──────────────────────────────────────────────

export interface AssignmentGroup {
  id: string;
  bootcampId: string;
  createdBy: string;
  title: string;
  description?: string;
  deadlineDays: number;
  createdAt: string;
  problems?: Problem[];      // populated via assignment_group_problems
}

export interface AssignmentProblem {
  id: string;
  assignmentId: string;
  problemId: string;
  problem: Problem;
  status: ProblemStatus;
  // Mapped to menteeStatus in UI layer
  menteeStatus: MenteeStatus;
  solutionLink?: string;
  notes?: string;            // "write note" in wireframe
  remarkForSelf?: string;    // "remark for self" in wireframe
  remarkForMentor?: string;  // "remark for mentor" in wireframe
  completedAt?: string;
  doubt?: Doubt;
}

export interface Assignment {
  id: string;
  assignmentGroupId: string;
  bootcampEnrollmentId: string;
  assignedBy: string;
  assignedAt: string;
  deadlineAt: string;
  status: AssignmentStatus;
  assignmentGroup: AssignmentGroup;
  problems: AssignmentProblem[];
  // Computed
  totalProblems: number;
  completedProblems: number;
  progressPercent: number;
}

// ─── Doubts ───────────────────────────────────────────────────

export interface Doubt {
  id: string;
  assignmentProblemId: string;
  raisedBy: string;          // org_member.id
  message: string;
  resolved: boolean;
  resolvedBy?: string;
  resolvedAt?: string;
  createdAt: string;
}

// ─── Leaderboard ──────────────────────────────────────────────

export interface LeaderboardEntry {
  id: string;
  bootcampId: string;
  bootcampEnrollmentId: string;
  problemsCompleted: number;
  problemsAttempted: number;
  completionRate: number;
  streakDays: number;
  score: number;
  rank: number;
  calculatedAt: string;
  // Joined
  user: Pick<User, 'id' | 'name' | 'avatarUrl'>;
}

// ─── Navigation Param Lists ───────────────────────────────────

export type RootStackParamList = {
  Auth: undefined;
  MentorApp: undefined;
  MenteeApp: undefined;
};

export type AuthStackParamList = {
  Login: undefined;
};

export type MenteeStackParamList = {
  Dashboard: undefined;
  AssignmentDetail: { assignmentId: string };
  ProblemDetail: {
    assignmentProblemId: string;
    problemTitle: string;
    isCompleted?: boolean;
  };
  CompletedProblems: undefined;
  Profile: undefined;
};

export type MentorStackParamList = {
  Dashboard: undefined;
  MenteeList: { day?: string };
  AssignTask: { menteeEnrollmentId?: string };
  QuestionBank: undefined;
  MenteeProgress: { enrollmentId: string; menteeName: string };
  Profile: undefined;
};

// ─── API Response wrappers ────────────────────────────────────

export interface ApiResponse<T> {
  data: T;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

// ─── App-level session context ────────────────────────────────

// Stored after login — the "active context" the user is operating in
export interface AppSession {
  user: User;
  accessToken: string;
  orgRole: OrgRole;
  bootcampRole?: BootcampRole;
  activeOrgId: string;
  activeBootcampId: string;
  orgMemberId: string;       // organization_member.id — critical for all API calls
  bootcampEnrollmentId: string; // bootcamp_enrollment.id
}