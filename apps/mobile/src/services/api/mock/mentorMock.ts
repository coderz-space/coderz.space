import { IMentorService } from '../interfaces';
import {
  OrgMember, Assignment, AssignmentGroup, Problem,
  PaginatedResponse, Doubt,
} from '../../../types';
import {
  MOCK_MENTEES, MOCK_ACTIVE_ASSIGNMENT, MOCK_COMPLETED_ASSIGNMENT,
  MOCK_ASSIGNMENT_GROUP, MOCK_ASSIGNMENT_GROUP_2, MOCK_PROBLEMS,
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

export const mentorMock: IMentorService = {
  async getMentees(): Promise<OrgMember[]> {
    await delay(800);
    return MOCK_MENTEES;
  },

  async getMenteeProgress({ enrollmentId }): Promise<Assignment[]> {
    await delay(700);
    return [MOCK_ACTIVE_ASSIGNMENT, MOCK_COMPLETED_ASSIGNMENT];
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
        p.title.toLowerCase().includes(search.toLowerCase()),
      );
    }
    if (difficulty) {
      filtered = filtered.filter((p) => p.difficulty === difficulty);
    }
    return { data: filtered, total: filtered.length, page: 1, limit: 20 };
  },

  async assignToMentee({ assignmentGroupId, bootcampEnrollmentId }): Promise<void> {
    await delay(1000);
    console.log('[Mock] Assigned group', assignmentGroupId, 'to', bootcampEnrollmentId);
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
        : d,
    );
    return mutableDoubts.find((d) => d.id === doubtId)!;
  },
};