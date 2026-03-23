import apiClient from '../apiClient';
import { IDoubtService } from '../interfaces';
import { Doubt } from '../../../types';

export const doubtLive: IDoubtService = {
  async createDoubt({ assignmentProblemId, message }) {
    const { data } = await apiClient.post('/doubts', {
      assignment_problem_id: assignmentProblemId,
      message,
    });
    return data.data;
  },

  async listDoubts({ resolved, raisedBy } = {}) {
    const { data } = await apiClient.get('/doubts', { params: { resolved, raised_by: raisedBy } });
    return { data: data.data, nextCursor: data.next_cursor };
  },

  async getDoubt(doubtId) {
    const { data } = await apiClient.get(`/doubts/${doubtId}`);
    return data.data;
  },

  async resolveDoubt({ doubtId, note }) {
    const { data } = await apiClient.patch(`/doubts/${doubtId}/resolve`, { note });
    return data.data;
  },

  async deleteDoubt(doubtId) {
    await apiClient.delete(`/doubts/${doubtId}`);
  },

  async getMyDoubts({ resolved } = {}) {
    const { data } = await apiClient.get('/doubts/me', { params: { resolved } });
    return data.data;
  },
};