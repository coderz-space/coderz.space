// src/services/api/mock/mentorMock.ts

import { IMentorService } from '../interfaces';
import {
  OrgMember,
  Assignment,
  AssignmentGroup,
  Problem,
  PaginatedResponse,
  Doubt,
} from '../../../types';
import {
  MOCK_MENTEES,
  MOCK_ACTIVE_ASSIGNMENT,
  MOCK_COMPLETED_ASSIGNMENT,
  MOCK_ASSIGNMENT_GROUP,
  MOCK_ASSIGNMENT_GROUP_2,
  MOCK_PROBLEMS,
} from './_mockData';

const delay = (ms: number) => new Promise<void>((r) => setTimeout(() => r(), ms));

const MOCK_DOUBTS: Doubt[] = [
  {
    id: 'doubt-1',
    assignmentProblemId: 'ap-2',
    raisedBy: 'orgmember-mentee-1',
    message: 'Not sure about the window shrink condition in sliding window',
    resolved: false,
    createdAt: new Date(Date.now() - 3600000).toISOString(),
  },
  {
    id: 'doubt-2',
    assignmentProblemId: 'ap-3',
    raisedBy: 'orgmember-mentee-2',
    message: 'How do we handle the case when the cache is empty initially?',
    resolved: false,
    createdAt: new Date(Date.now() - 7200000).toISOString(),
  },
];
let mutableDoubts = [...MOCK_DOUBTS];

// In-memory store for assignments created with custom problems
let customAssignments: Assignment[] = [];

export const mentorMock: IMentorService = {
  async getMentees(): Promise<OrgMember[]> {
    await delay(800);
    return MOCK_MENTEES;
  },

  async getMenteeProgress({ enrollmentId }): Promise<Assignment[]> {
    await delay(700);
    // Combine regular assignments with custom ones for this mentee
    const regular = [MOCK_ACTIVE_ASSIGNMENT, MOCK_COMPLETED_ASSIGNMENT];
    const custom = customAssignments.filter(a => a.bootcampEnrollmentId === enrollmentId);
    return [...regular, ...custom];
  },

  async getAssignmentGroups(): Promise<AssignmentGroup[]> {
    await delay(600);
    return [MOCK_ASSIGNMENT_GROUP, MOCK_ASSIGNMENT_GROUP_2];
  },

  async getProblems({ search, difficulty }): Promise<PaginatedResponse<Problem>> {
    await delay(500);
    let filtered = [...MOCK_PROBLEMS];
    if (search) {
      filtered = filtered.filter((p) =>
        p.title.toLowerCase().includes(search.toLowerCase())
      );
    }
    if (difficulty) {
      filtered = filtered.filter((p) => p.difficulty === difficulty);
    }
    return { data: filtered, total: filtered.length, page: 1, limit: 20 };
  },

  // Original assignToMentee: assign an existing assignment group
  async assignToMentee(params: {
    orgId: string;
    bootcampId: string;
    assignmentGroupId: string;
    bootcampEnrollmentId: string;
    deadlineAt: string;
  }): Promise<void> {
    await delay(1000);
    console.log('[Mock] Assigned group', params.assignmentGroupId, 'to', params.bootcampEnrollmentId);
  },

  // New method: assign custom problems directly (creating a new assignment group on the fly)
  async assignProblemsToMentee(params: {
    orgId: string;
    bootcampId: string;
    bootcampEnrollmentId: string;
    problemIds: string[];
    deadlineAt: string;
  }): Promise<void> {
    await delay(1000);
    // Find the problems from MOCK_PROBLEMS
    const problems = MOCK_PROBLEMS.filter(p => params.problemIds.includes(p.id));
    if (problems.length === 0) {
      throw new Error('No valid problems to assign');
    }
    // Create a temporary assignment group for this assignment
    const newGroupId = `ag_custom_${Date.now()}`;
    const newGroup: AssignmentGroup = {
      id: newGroupId,
      bootcampId: params.bootcampId,
      title: `Custom Assignment (${new Date().toLocaleDateString()})`,
      description: `Assigned ${problems.length} problems`,
      deadlineDays: Math.ceil((new Date(params.deadlineAt).getTime() - Date.now()) / (1000 * 60 * 60 * 24)),
      createdBy: 'orgmember-mentor-1', // mock mentor id
      createdAt: new Date().toISOString(),
      problems: problems, // populate for convenience
    };
    // Build assignment problems
    const assignmentProblems = problems.map((p, idx) => ({
      id: `ap_${newGroupId}_${p.id}`,
      assignmentId: `assign_${newGroupId}`,
      problemId: p.id,
      problem: p,
      status: 'pending' as const,
      menteeStatus: 'not_started' as const,
      solutionLink: undefined,
      notes: undefined,
      remarkForSelf: undefined,
      remarkForMentor: undefined,
    }));
    const newAssignment: Assignment = {
      id: `assign_${newGroupId}`,
      assignmentGroupId: newGroupId,
      bootcampEnrollmentId: params.bootcampEnrollmentId,
      assignedBy: 'orgmember-mentor-1',
      assignedAt: new Date().toISOString(),
      deadlineAt: params.deadlineAt,
      status: 'active',
      assignmentGroup: newGroup,
      problems: assignmentProblems,
      totalProblems: assignmentProblems.length,
      completedProblems: 0,
      progressPercent: 0,
    };
    customAssignments.push(newAssignment);
    console.log('[Mock] Assigned custom problems', params.problemIds, 'to', params.bootcampEnrollmentId);
  },

  async getPendingDoubts(): Promise<Doubt[]> {
    await delay(700);
    return mutableDoubts.filter((d) => !d.resolved);
  },

  async resolveDoubt({ doubtId }): Promise<Doubt> {
    await delay(600);
    mutableDoubts = mutableDoubts.map((d) =>
      d.id === doubtId
        ? { ...d, resolved: true, resolvedBy: 'orgmember-mentor-1', resolvedAt: new Date().toISOString() }
        : d
    );
    return mutableDoubts.find((d) => d.id === doubtId)!;
  },
};