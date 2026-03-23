import apiClient, { ApiRoutes } from '../apiClient';
import { IMenteeService } from '../interfaces';
import {
  Assignment, AssignmentProblem, MenteeStatus, Doubt, LeaderboardEntry,
} from '../../../types';

// Maps our MenteeStatus to backend ProblemStatus
const statusMap: Record<MenteeStatus, string> = {
  not_started: 'pending',
  discussion_needed: 'attempted',
  revision_needed: 'attempted',
  completed: 'completed',
};

export const menteeLive: IMenteeService = {
  async getMyAssignments({ orgId, bootcampId, enrollmentId }) {
    // GET /assignments?bootcamp_enrollment_id=X&status=active
    const { data } = await apiClient.get('/assignments', {
      params: { bootcamp_enrollment_id: enrollmentId, bootcamp_id: bootcampId, status: 'active' },
    });
    return data.data;
  },

  async getAssignmentDetail({ assignmentId }) {
    // GET /assignments/{assignment_id}
    const { data } = await apiClient.get(`/assignments/${assignmentId}`);
    return data.data;
  },

  async getCompletedAssignments({ enrollmentId, bootcampId }) {
    const { data } = await apiClient.get('/assignments', {
      params: { bootcamp_enrollment_id: enrollmentId, bootcamp_id: bootcampId, status: 'completed' },
    });
    return data.data;
  },

  async updateProblemProgress({ assignmentId, assignmentProblemId, status, solutionLink, notes, remarkForSelf, remarkForMentor }) {
    // PATCH /assignments/{assignment_id}/problems/{problem_id}
    // NOTE: remarkForSelf and remarkForMentor are stored in notes field as JSON
    // until backend adds dedicated fields
    const body: Record<string, any> = {};
    if (status) body.status = statusMap[status];
    if (solutionLink !== undefined) body.solution_link = solutionLink;
    if (notes !== undefined) body.notes = notes;
    // Extended fields — add when backend supports them
    if (remarkForSelf !== undefined) body.remark_for_self = remarkForSelf;
    if (remarkForMentor !== undefined) body.remark_for_mentor = remarkForMentor;

    const { data } = await apiClient.patch(
      `/assignments/${assignmentId}/problems/${assignmentProblemId}`,
      body,
    );
    return data.data;
  },

  async raiseDoubt({ assignmentProblemId, message }) {
    // POST /doubts
    const { data } = await apiClient.post('/doubts', { assignment_problem_id: assignmentProblemId, message });
    return data.data;
  },

  async getLeaderboard({ orgId, bootcampId }) {
    const { data } = await apiClient.get(
      ApiRoutes.leaderboard(orgId, bootcampId),
    );
    return data.data;
  },
};