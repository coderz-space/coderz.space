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
} from '../../types';

// ─── Auth ─────────────────────────────────────────────────────

export interface IAuthService {
  login(payload: LoginPayload): Promise<LoginResponse>;
  logout(): Promise<void>;
  refreshToken(): Promise<{ accessToken: string }>;
  getMe(): Promise<LoginResponse>;
}

// ─── Mentee ───────────────────────────────────────────────────

export interface IMenteeService {
  /**
   * Get all active assignments for the logged-in mentee
   * GET /orgs/:orgId/bootcamps/:bootcampId/enrollments/:enrollmentId/assignments
   */
  getMyAssignments(params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }): Promise<Assignment[]>;

  /**
   * Get a single assignment with all problems
   * GET /orgs/:orgId/bootcamps/:bootcampId/assignments/:assignmentId
   */
  getAssignmentDetail(params: {
    orgId: string;
    bootcampId: string;
    assignmentId: string;
  }): Promise<Assignment>;

  /**
   * Get completed assignments
   * GET /orgs/:orgId/bootcamps/:bootcampId/enrollments/:enrollmentId/assignments?status=completed
   */
  getCompletedAssignments(params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }): Promise<Assignment[]>;

  /**
   * Update a problem's status, solution link, notes, remarks
   * PATCH /orgs/:orgId/bootcamps/:bootcampId/assignments/:assignmentId/problems/:assignmentProblemId
   */
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

  /**
   * Raise a doubt on a problem
   * POST /orgs/:orgId/bootcamps/:bootcampId/assignments/:assignmentId/problems/:assignmentProblemId/doubts
   */
  raiseDoubt(params: {
    orgId: string;
    bootcampId: string;
    assignmentId: string;
    assignmentProblemId: string;
    message: string;
  }): Promise<Doubt>;

  /**
   * Get leaderboard for the bootcamp
   * GET /orgs/:orgId/bootcamps/:bootcampId/leaderboard
   */
  getLeaderboard(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<LeaderboardEntry[]>;
}

// ─── Mentor ───────────────────────────────────────────────────

export interface IMentorService {
  /**
   * Get all mentees in this bootcamp
   * GET /orgs/:orgId/bootcamps/:bootcampId/members?role=mentee
   */
  getMentees(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<OrgMember[]>;

  /**
   * Get a specific mentee's assignments / progress
   * GET /orgs/:orgId/bootcamps/:bootcampId/enrollments/:enrollmentId/assignments
   */
  getMenteeProgress(params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }): Promise<Assignment[]>;

  /**
   * Get all assignment groups (templates) in a bootcamp
   * GET /orgs/:orgId/bootcamps/:bootcampId/assignment-groups
   */
  getAssignmentGroups(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<AssignmentGroup[]>;

  /**
   * Get problems in the org problem bank
   * GET /orgs/:orgId/problems
   */
  getProblems(params: {
    orgId: string;
    search?: string;
    difficulty?: string;
    tagId?: string;
  }): Promise<PaginatedResponse<Problem>>;

  /**
   * Assign an assignment group to a mentee
   * POST /orgs/:orgId/bootcamps/:bootcampId/assignments
   */
  assignToMentee(params: {
    orgId: string;
    bootcampId: string;
    assignmentGroupId: string;
    bootcampEnrollmentId: string;
    deadlineAt: string;
  }): Promise<void>;

  /**
   * Get all unresolved doubts in the bootcamp
   * GET /orgs/:orgId/bootcamps/:bootcampId/doubts?resolved=false
   */
  getPendingDoubts(params: {
    orgId: string;
    bootcampId: string;
  }): Promise<Doubt[]>;

  /**
   * Resolve a doubt
   * PATCH /orgs/:orgId/bootcamps/:bootcampId/doubts/:doubtId/resolve
   */
  resolveDoubt(params: {
    orgId: string;
    bootcampId: string;
    doubtId: string;
  }): Promise<Doubt>;
}