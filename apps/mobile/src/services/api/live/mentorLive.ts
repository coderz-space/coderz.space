import apiClient, { ApiRoutes } from '../apiClient';
import { IMentorService } from '../interfaces';
import { OrgMember, Assignment, AssignmentGroup, Problem, PaginatedResponse, Doubt } from '../../../types';

export const mentorLive: IMentorService = {
  async getMentees({ orgId, bootcampId }) {
    // GET /orgs/{org_id}/b/{bootcamp_id}/enrollments?role=mentee
    const { data } = await apiClient.get(
      `/orgs/${orgId}/b/${bootcampId}/enrollments`,
      { params: { role: 'mentee' } },
    );
    return data.items;
  },

  async getMenteeProgress({ enrollmentId, bootcampId }) {
    // GET /assignments?bootcamp_enrollment_id=X
    const { data } = await apiClient.get('/assignments', {
      params: { bootcamp_enrollment_id: enrollmentId, bootcamp_id: bootcampId },
    });
    return data.data;
  },

  async getAssignmentGroups({ bootcampId }) {
    // GET /b/{bootcamp_id}/agroups
    const { data } = await apiClient.get(`/b/${bootcampId}/agroups`);
    return data.data;
  },

  async getProblems({ orgId, search, difficulty, tagId }) {
    const { data } = await apiClient.get(
      ApiRoutes.problems(orgId),
      { params: { q: search, difficulty, tag_id: tagId, limit: 20 } },
    );
    return { data: data.items, total: data.total, page: data.page, limit: data.limit };
  },

  async assignToMentee({ assignmentGroupId, bootcampEnrollmentId, deadlineAt }) {
    // POST /assignments
    await apiClient.post('/assignments', {
      assignment_group_id: assignmentGroupId,
      bootcamp_enrollment_id: bootcampEnrollmentId,
      deadline_at: deadlineAt,
    });
  },

  async getPendingDoubts({ orgId, bootcampId }) {
    // GET /doubts?resolved=false (scoped by org/bootcamp via token context)
    const { data } = await apiClient.get('/doubts', { params: { resolved: false } });
    return data.data;
  },

  async resolveDoubt({ doubtId }) {
    // PATCH /doubts/{id}/resolve
    const { data } = await apiClient.patch(`/doubts/${doubtId}/resolve`);
    return data.data;
  },
};