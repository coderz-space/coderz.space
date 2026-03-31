import {
  LoginPayload,
  LoginResponse,
  Assignment,
  AssignmentProblem,
  MenteeStatus,
  Problem,
  AssignmentGroup,
  OrgMember,
  LeaderboardEntry,
  Doubt,
  PaginatedResponse,

  // ✅ NEW TYPES
  SignupPayload,
  SignupResponse,
  Poll,
  PollResult,
  PollVote,
  Organization,
  OrgRole,
  Bootcamp,
  BootcampEnrollment,
  BootcampRole,
  Difficulty,
  Tag,
} from '../../types';

// ─── Auth ─────────────────────────────────────────────────────

export interface IAuthService {
  login(payload: LoginPayload): Promise<LoginResponse>;
  logout(): Promise<void>;
  refreshToken(): Promise<{ accessToken: string }>;
  getMe(): Promise<LoginResponse>;
}

// ✅ UPDATED AUTH (V2)
export interface IAuthServiceV2 extends IAuthService {
  signup(payload: SignupPayload): Promise<SignupResponse>;
  forgotPassword(email: string): Promise<void>;
  resetPassword(params: { token: string; newPassword: string }): Promise<void>;
}

// ─── Mentee ───────────────────────────────────────────────────

export interface IMenteeService {
  getMyAssignments(params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }): Promise<Assignment[]>;

  getAssignmentDetail(params: {
    orgId: string;
    bootcampId: string;
    assignmentId: string;
  }): Promise<Assignment>;

  getCompletedAssignments(params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }): Promise<Assignment[]>;

  updateProblemProgress(params: {
    orgId: string;
    bootcampId: string;
    assignmentId: string;
    assignmentProblemId: string;
    status?: MenteeStatus;
    solutionLink?: string;
    notes?: string;
    remarkForSelf?: string;
    remarkForMentor?: string;
  }): Promise<AssignmentProblem>;

  raiseDoubt(params: {
    orgId: string;
    bootcampId: string;
    assignmentId: string;
    assignmentProblemId: string;
    message: string;
  }): Promise<Doubt>;

  getLeaderboard(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<LeaderboardEntry[]>;
}

// ─── Mentor ───────────────────────────────────────────────────

export interface IMentorService {
  getMentees(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<OrgMember[]>;

  getMenteeProgress(params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }): Promise<Assignment[]>;

  getAssignmentGroups(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<AssignmentGroup[]>;

  getProblems(params: {
    orgId: string;
    search?: string;
    difficulty?: string;
    tagId?: string;
  }): Promise<PaginatedResponse<Problem>>;

  assignToMentee(params: {
    orgId: string;
    bootcampId: string;
    assignmentGroupId: string;
    bootcampEnrollmentId: string;
    deadlineAt: string;
  }): Promise<void>;

  getPendingDoubts(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<Doubt[]>;

  resolveDoubt(params: {
    orgId: string;
    bootcampId: string;
    doubtId: string;
  }): Promise<Doubt>;
}

// ─── Org Service ──────────────────────────────────────────────

export interface IOrgService {
  createOrg(params: {
    name: string;
    slug: string;
    description?: string;
  }): Promise<{ id: string; status: string }>;

  getMyOrgs(): Promise<Organization[]>;

  getOrg(slug: string): Promise<Organization>;

  getOrgMembers(params: {
    orgSlug: string;
    role?: OrgRole;
  }): Promise<OrgMember[]>;

  addMember(params: {
    orgSlug: string;
    userId: string;
    role: OrgRole;
  }): Promise<OrgMember>;

  updateMemberRole(params: {
    orgSlug: string;
    userId: string;
    role: OrgRole;
  }): Promise<OrgMember>;

  removeMember(params: {
    orgSlug: string;
    userId: string;
  }): Promise<void>;
}

// ─── Bootcamp Service ─────────────────────────────────────────

export interface IBootcampService {
  listBootcamps(params: {
    orgId: string;
  }): Promise<Bootcamp[]>;

  getBootcamp(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<Bootcamp>;

  getEnrollments(params: {
    orgId: string;
    bootcampId: string;
    role?: BootcampRole;
  }): Promise<BootcampEnrollment[]>;

  enroll(params: {
    orgId: string;
    bootcampId: string;
    orgMemberId: string;
    role: BootcampRole;
  }): Promise<BootcampEnrollment>;

  removeEnrollment(params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }): Promise<void>;
}

// ─── Problem Service ──────────────────────────────────────────

export interface IProblemService {
  createProblem(params: {
    orgId: string;
    title: string;
    description?: string;
    difficulty: Difficulty;
    externalLink?: string;
  }): Promise<Problem>;

  listProblems(params: {
    orgId: string;
    q?: string;
    difficulty?: Difficulty;
    tagId?: string;
    page?: number;
    limit?: number;
  }): Promise<PaginatedResponse<Problem>>;

  getProblem(params: {
    orgId: string;
    problemId: string;
  }): Promise<Problem>;

  updateProblem(params: {
    orgId: string;
    problemId: string;
    title?: string;
    description?: string;
    difficulty?: Difficulty;
    externalLink?: string;
  }): Promise<Problem>;

  deleteProblem(params: {
    orgId: string;
    problemId: string;
  }): Promise<void>;

  createTag(params: {
    orgId: string;
    name: string;
  }): Promise<Tag>;

  listTags(params: {
    orgId: string;
    q?: string;
  }): Promise<Tag[]>;

  attachTags(params: {
    orgId: string;
    problemId: string;
    tagIds: string[];
  }): Promise<void>;

  detachTag(params: {
    orgId: string;
    problemId: string;
    tagId: string;
  }): Promise<void>;
}

// ─── Doubt Service ────────────────────────────────────────────

export interface IDoubtService {
  createDoubt(params: {
    assignmentProblemId: string;
    message: string;
  }): Promise<Doubt>;

  listDoubts(params: {
    resolved?: boolean;
    raisedBy?: string;
  }): Promise<{ data: Doubt[]; nextCursor?: string }>;

  getDoubt(doubtId: string): Promise<Doubt>;

  resolveDoubt(params: {
    doubtId: string;
    note?: string;
  }): Promise<Doubt>;

  deleteDoubt(doubtId: string): Promise<void>;

  getMyDoubts(params: {
    resolved?: boolean;
  }): Promise<Doubt[]>;
}

// ─── Analytics Service ────────────────────────────────────────

export interface IAnalyticsService {
  getLeaderboard(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<LeaderboardEntry[]>;

  getMyLeaderboardEntry(params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }): Promise<LeaderboardEntry>;

  createPoll(params: {
    orgId: string;
    bootcampId: string;
    problemId: string;
    question: string;
  }): Promise<Poll>;

  listPolls(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<Poll[]>;

  getPoll(params: {
    orgId: string;
    bootcampId: string;
    pollId: string;
  }): Promise<Poll>;

  votePoll(params: {
    orgId: string;
    bootcampId: string;
    pollId: string;
    vote: PollVote;
  }): Promise<void>;

  getPollResults(params: {
    orgId: string;
    bootcampId: string;
    pollId: string;
  }): Promise<PollResult>;
}