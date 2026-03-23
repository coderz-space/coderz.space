import apiClient from '../apiClient';
import { IOrgService } from '../interfaces';
import { Organization, OrgMember, OrgRole } from '../../../types';

export const orgLive: IOrgService = {
  async createOrg({ name, slug, description }) {
    const { data } = await apiClient.post('/orgs', { name, slug, description });
    return data;
  },

  async getMyOrgs() {
    const { data } = await apiClient.get('/orgs');
    return data.data;
  },

  async getOrg(slug) {
    const { data } = await apiClient.get(`/orgs/${slug}`);
    return data.data;
  },

  async getOrgMembers({ orgSlug, role }) {
    const { data } = await apiClient.get(`/orgs/${orgSlug}/members`, { params: { role } });
    return data.data;
  },

  async addMember({ orgSlug, userId, role }) {
    const { data } = await apiClient.post(`/orgs/${orgSlug}/members`, { user_id: userId, role });
    return data.data;
  },

  async updateMemberRole({ orgSlug, userId, role }) {
    const { data } = await apiClient.patch(`/orgs/${orgSlug}/members/${userId}`, { role });
    return data.data;
  },

  async removeMember({ orgSlug, userId }) {
    await apiClient.delete(`/orgs/${orgSlug}/members/${userId}`);
  },
};