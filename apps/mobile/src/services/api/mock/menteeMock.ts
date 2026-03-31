import { IMenteeService } from '../interfaces';
import {
  Assignment, AssignmentProblem, MenteeStatus, Doubt, LeaderboardEntry,
} from '../../../types';
import {
  MOCK_ACTIVE_ASSIGNMENT, MOCK_COMPLETED_ASSIGNMENT, MOCK_LEADERBOARD,
} from './_mockData';

const delay = (ms: number) => new Promise<void>((r) => setTimeout(() => r(), ms));

// Local mutable copy so updates persist during session
let activeAssignment = { ...MOCK_ACTIVE_ASSIGNMENT };

export const menteeMock: IMenteeService = {
  async getMyAssignments(): Promise<Assignment[]> {
    await delay(800);
    return [activeAssignment];
  },

  async getAssignmentDetail({ assignmentId }): Promise<Assignment> {
    await delay(600);
    if (assignmentId === activeAssignment.id) return activeAssignment;
    if (assignmentId === MOCK_COMPLETED_ASSIGNMENT.id) return MOCK_COMPLETED_ASSIGNMENT;
    throw new Error('Assignment not found');
  },

  async getCompletedAssignments(): Promise<Assignment[]> {
    await delay(700);
    return [MOCK_COMPLETED_ASSIGNMENT];
  },

  async updateProblemProgress({ assignmentProblemId, status, solutionLink, notes, remarkForSelf, remarkForMentor }): Promise<AssignmentProblem> {
    await delay(600);

    const idx = activeAssignment.problems.findIndex((p) => p.id === assignmentProblemId);
    if (idx === -1) throw new Error('Problem not found in assignment');

    const updated: AssignmentProblem = {
      ...activeAssignment.problems[idx],
      ...(status && { menteeStatus: status }),
      ...(status === 'completed' && { status: 'completed', completedAt: new Date().toISOString() }),
      ...(status === 'discussion_needed' && { status: 'attempted' }),
      ...(status === 'revision_needed' && { status: 'attempted' }),
      ...(status === 'not_started' && { status: 'pending' }),
      ...(solutionLink !== undefined && { solutionLink }),
      ...(notes !== undefined && { notes }),
      ...(remarkForSelf !== undefined && { remarkForSelf }),
      ...(remarkForMentor !== undefined && { remarkForMentor }),
    };

    const updatedProblems = [...activeAssignment.problems];
    updatedProblems[idx] = updated;
    const completed = updatedProblems.filter((p) => p.status === 'completed').length;

    activeAssignment = {
      ...activeAssignment,
      problems: updatedProblems,
      completedProblems: completed,
      progressPercent: Math.round((completed / updatedProblems.length) * 100),
    };

    return updated;
  },

  async raiseDoubt({ assignmentProblemId, message }): Promise<Doubt> {
    await delay(700);
    const doubt: Doubt = {
      id: `doubt-${Date.now()}`,
      assignmentProblemId,
      raisedBy: 'orgmember-mentee-1',
      message,
      resolved: false,
      createdAt: new Date().toISOString(),
    };

    const idx = activeAssignment.problems.findIndex((p) => p.id === assignmentProblemId);
    if (idx !== -1) {
      const updatedProblems = [...activeAssignment.problems];
      updatedProblems[idx] = {
        ...updatedProblems[idx],
        doubt,
        menteeStatus: 'discussion_needed',
        status: 'attempted',
      };
      activeAssignment = { ...activeAssignment, problems: updatedProblems };
    }

    return doubt;
  },

  async getLeaderboard(): Promise<LeaderboardEntry[]> {
    await delay(900);
    return MOCK_LEADERBOARD;
  },
};