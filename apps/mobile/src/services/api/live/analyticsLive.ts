import apiClient, { ApiRoutes } from '../apiClient';
import { IAnalyticsService } from '../interfaces';
import { LeaderboardEntry, Poll, PollResult, PollVote } from '../../../types';

export const analyticsLive: IAnalyticsService = {
  async getLeaderboard({ orgId, bootcampId }) {
    const { data } = await apiClient.get(ApiRoutes.leaderboard(orgId, bootcampId));
    return data.data;
  },

  async getMyLeaderboardEntry({ orgId, bootcampId, enrollmentId }) {
    const { data } = await apiClient.get(
      `${ApiRoutes.leaderboard(orgId, bootcampId)}/${enrollmentId}`,
    );
    return data.data;
  },

  async createPoll({ orgId, bootcampId, problemId, question }) {
    const { data } = await apiClient.post(
      `/orgs/${orgId}/bootcamps/${bootcampId}/polls`,
      { problem_id: problemId, question },
    );
    return data.data;
  },

  async listPolls({ orgId, bootcampId }) {
    const { data } = await apiClient.get(`/orgs/${orgId}/bootcamps/${bootcampId}/polls`);
    return data.data;
  },

  async getPoll({ orgId, bootcampId, pollId }) {
    const { data } = await apiClient.get(
      `/orgs/${orgId}/bootcamps/${bootcampId}/polls/${pollId}`,
    );
    return data.data;
  },

  async votePoll({ orgId, bootcampId, pollId, vote }) {
    await apiClient.put(
      `/orgs/${orgId}/bootcamps/${bootcampId}/polls/${pollId}/vote`,
      { vote },
    );
  },

  async getPollResults({ orgId, bootcampId, pollId }) {
    const { data } = await apiClient.get(
      `/orgs/${orgId}/bootcamps/${bootcampId}/polls/${pollId}/results`,
    );
    return data.data;
  },
};