import { create } from 'zustand';
import { OrgMember, Assignment, AssignmentGroup, Problem, Doubt, PaginatedResponse } from '../types';
import { mentorMock as service } from '../services/api/mock/mentorMock';
// Swap with live service when backend ready

interface MentorStore {
  mentees: OrgMember[];
  menteeProgress: Record<string, Assignment[]>;   // keyed by enrollmentId
  assignmentGroups: AssignmentGroup[];
  problems: Problem[];
  pendingDoubts: Doubt[];
  isLoadingMentees: boolean;
  isLoadingGroups: boolean;
  isLoadingDoubts: boolean;
  error: string | null;

  fetchMentees: (params: { orgId: string; bootcampId: string }) => Promise<void>;
  fetchMenteeProgress: (params: { orgId: string; bootcampId: string; enrollmentId: string }) => Promise<void>;
  fetchAssignmentGroups: (params: { orgId: string; bootcampId: string }) => Promise<void>;
  fetchProblems: (params: { orgId: string; search?: string; difficulty?: string }) => Promise<void>;
  fetchPendingDoubts: (params: { orgId: string; bootcampId: string }) => Promise<void>;
  assignToMentee: (params: {
    orgId: string;
    bootcampId: string;
    assignmentGroupId: string;
    bootcampEnrollmentId: string;
    deadlineAt: string;
  }) => Promise<void>;
  assignProblemsToMentee: (params: {
  orgId: string;
  bootcampId: string;
  bootcampEnrollmentId: string;
  problemIds: string[];
  deadlineAt: string;
}) => Promise<void>;
  resolveDoubt: (params: { orgId: string; bootcampId: string; doubtId: string }) => Promise<void>;
}

export const useMentorStore = create<MentorStore>((set) => ({
  mentees: [],
  menteeProgress: {},
  assignmentGroups: [],
  problems: [],
  pendingDoubts: [],
  isLoadingMentees: false,
  isLoadingGroups: false,
  isLoadingDoubts: false,
  error: null,

  fetchMentees: async (params) => {
    set({ isLoadingMentees: true, error: null });
    try {
      const data = await service.getMentees(params);
      set({ mentees: data });
    } catch (e: any) {
      set({ error: e.message });
    } finally {
      set({ isLoadingMentees: false });
    }
  },

  fetchMenteeProgress: async (params) => {
    try {
      const data = await service.getMenteeProgress(params);
      set((s) => ({
        menteeProgress: { ...s.menteeProgress, [params.enrollmentId]: data },
      }));
    } catch (e: any) {
      set({ error: e.message });
    }
  },

  fetchAssignmentGroups: async (params) => {
    set({ isLoadingGroups: true });
    try {
      const data = await service.getAssignmentGroups(params);
      set({ assignmentGroups: data });
    } catch (e: any) {
      set({ error: e.message });
    } finally {
      set({ isLoadingGroups: false });
    }
  },

  fetchProblems: async (params) => {
    try {
      const res = await service.getProblems(params);
      set({ problems: res.data });
    } catch (e: any) {
      set({ error: e.message });
    }
  },

  fetchPendingDoubts: async (params) => {
    set({ isLoadingDoubts: true });
    try {
      const data = await service.getPendingDoubts(params);
      set({ pendingDoubts: data });
    } catch (e: any) {
      set({ error: e.message });
    } finally {
      set({ isLoadingDoubts: false });
    }
  },

  assignToMentee: async (params) => {
    await service.assignToMentee(params);
  },

  assignProblemsToMentee: async (params) => {
  await service.assignProblemsToMentee(params);
},

  resolveDoubt: async (params) => {
    await service.resolveDoubt(params);
    set((s) => ({
      pendingDoubts: s.pendingDoubts.filter((d) => d.id !== params.doubtId),
    }));
  },
}));