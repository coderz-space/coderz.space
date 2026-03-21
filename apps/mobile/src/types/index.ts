// ─── User & Auth ──────────────────────────────────────────────

export type UserRole = 'mentor' | 'mentee' | 'admin';

export interface User {
  id: string;
  name: string;
  email: string;
  role: UserRole;
  avatarUrl?: string;
  cohort?: string;        // e.g. "Batch 12"
  joinedAt: string;       // ISO date string
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

// ─── Tasks ────────────────────────────────────────────────────

export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'review';
export type TaskDifficulty = 'easy' | 'medium' | 'hard';

export interface Task {
  id: string;
  title: string;
  description: string;
  status: TaskStatus;
  difficulty: TaskDifficulty;
  dueDate: string;        // ISO date string
  solutionUrl?: string;
  hasDoubt: boolean;
  doubtDescription?: string;
  assignedAt: string;
  completedAt?: string;
  tags: string[];
  points: number;
}

// ─── Weekly Goals ─────────────────────────────────────────────

export interface WeeklyGoal {
  id: string;
  weekNumber: number;
  startDate: string;
  endDate: string;
  totalTasks: number;
  completedTasks: number;
  tasks: Task[];
}

// ─── Mentee / Mentor Profiles ─────────────────────────────────

export interface MenteeProfile extends User {
  role: 'mentee';
  mentorId: string;
  currentStreak: number;    // days
  totalPoints: number;
  weeklyGoal?: WeeklyGoal;
  completionRate: number;   // 0-100
}

export interface MentorProfile extends User {
  role: 'mentor';
  menteeIds: string[];
  totalMentees: number;
  pendingDoubts: number;
}

// ─── Question Bank ────────────────────────────────────────────

export type QuestionCategory =
  | 'dsa'
  | 'system_design'
  | 'web'
  | 'mobile'
  | 'database'
  | 'other';

export interface Question {
  id: string;
  title: string;
  description: string;
  difficulty: TaskDifficulty;
  category: QuestionCategory;
  tags: string[];
  points: number;
  estimatedHours: number;
  resourceLinks?: string[];
}

// ─── Bootcamp / Cohort ────────────────────────────────────────

export interface Cohort {
  id: string;
  name: string;            // e.g. "Batch 12"
  startDate: string;
  endDate: string;
  mentorId: string;
  menteeIds: string[];
  isActive: boolean;
}

// ─── Navigation Param Lists ───────────────────────────────────

export type RootStackParamList = {
  Auth: undefined;
  MentorApp: undefined;
  MenteeApp: undefined;
};

export type AuthStackParamList = {
  Login: { title?: string };
  ForgotPassword: { title?: string };
};

export type MenteeStackParamList = {
  Dashboard: { title?: string };
  TaskDetail: { taskId: string; title?: string };
  Leaderboard: { title?: string };
  Profile: { title?: string };
};

export type MentorStackParamList = {
  Dashboard: { title?: string };
  MenteeProgress: { menteeId: string; title?: string };
  AssignTask: { menteeId?: string; title?: string };
  QuestionBank: { title?: string };
  Doubts: { title?: string };
  Profile: { title?: string };
};