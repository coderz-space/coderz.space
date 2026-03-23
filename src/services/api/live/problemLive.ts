import apiClient, { ApiRoutes } from '../apiClient';
import { IProblemService } from '../interfaces';
import { Problem, Tag, Difficulty, PaginatedResponse } from '../../../types';

export const problemLive: IProblemService = {
  async createProblem({ orgId, title, description, difficulty, externalLink }) {
    const { data } = await apiClient.post(ApiRoutes.problems(orgId), {
      title, description, difficulty, external_link: externalLink,
    });
    return data.data;
  },

  async listProblems({ orgId, q, difficulty, tagId, page = 1, limit = 20 }) {
    const { data } = await apiClient.get(ApiRoutes.problems(orgId), {
      params: { q, difficulty, tag_id: tagId, page, limit },
    });
    return { data: data.items, total: data.total, page: data.page, limit: data.limit };
  },

  async getProblem({ orgId, problemId }) {
    const { data } = await apiClient.get(ApiRoutes.problemDetail(orgId, problemId));
    return data.data;
  },

  async updateProblem({ orgId, problemId, ...rest }) {
    const { data } = await apiClient.patch(
      ApiRoutes.problemDetail(orgId, problemId),
      { title: rest.title, description: rest.description, difficulty: rest.difficulty, external_link: rest.externalLink },
    );
    return data.data;
  },

  async deleteProblem({ orgId, problemId }) {
    await apiClient.delete(ApiRoutes.problemDetail(orgId, problemId));
  },

  async createTag({ orgId, name }) {
    const { data } = await apiClient.post(`/orgs/${orgId}/tags`, { name });
    return data.data;
  },

  async listTags({ orgId, q }) {
    const { data } = await apiClient.get(`/orgs/${orgId}/tags`, { params: { q } });
    return data.items;
  },

  async attachTags({ orgId, problemId, tagIds }) {
    await apiClient.post(`/orgs/${orgId}/problems/${problemId}/tags`, { tag_ids: tagIds });
  },

  async detachTag({ orgId, problemId, tagId }) {
    await apiClient.delete(`/orgs/${orgId}/problems/${problemId}/tags/${tagId}`);
  },
};