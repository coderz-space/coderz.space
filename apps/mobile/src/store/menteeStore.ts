import { create } from 'zustand';
import { Assignment, AssignmentProblem, MenteeStatus, LeaderboardEntry, LeaderboardPeriod } from '../types';
import { menteeMock as service } from '../services/api/mock/menteeMock';
// When backend is ready, swap ↑ with:
// import { menteeLive as service } from '../services/api/live/menteeLive';

interface MenteeStore {
  activeAssignments: Assignment[];
  completedAssignments: Assignment[];
  leaderboard: LeaderboardEntry[];

  // ✅ NEW
  leaderboardPeriod: LeaderboardPeriod;
  setLeaderboardPeriod: (period: LeaderboardPeriod) => void;

  isLoadingAssignments: boolean;
  isLoadingCompleted: boolean;
  isLoadingLeaderboard: boolean;
  error: string | null;

  fetchMyAssignments: (params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }) => Promise<void>;

  fetchCompletedAssignments: (params: {
    orgId: string;
    bootcampId: string;
    enrollmentId: string;
  }) => Promise<void>;

  // ✅ UPDATED
  fetchLeaderboard: (params: {
    orgId: string;
    bootcampId: string;
    period?: LeaderboardPeriod;
  }) => Promise<void>;

  updateProblemProgress: (params: {
    orgId: string;
    bootcampId: string;
    assignmentId: string;
    assignmentProblemId: string;
    status?: MenteeStatus;
    solutionLink?: string;
    notes?: string;
    remarkForSelf?: string;
    remarkForMentor?: string;
  }) => Promise<AssignmentProblem>;

  raiseDoubt: (params: {
    orgId: string;
    bootcampId: string;
    assignmentId: string;
    assignmentProblemId: string;
    message: string;
  }) => Promise<void>;
}

export const useMenteeStore = create<MenteeStore>((set, get) => ({
  activeAssignments: [],
  completedAssignments: [],
  leaderboard: [],

  // ✅ NEW
  leaderboardPeriod: 'allTime',

  isLoadingAssignments: false,
  isLoadingCompleted: false,
  isLoadingLeaderboard: false,
  error: null,

  // ✅ NEW
  setLeaderboardPeriod: (period: LeaderboardPeriod) => {
    set({ leaderboardPeriod: period });
    // Optionally auto-refetch leaderboard with new period
    const state = get();
  },

  fetchMyAssignments: async (params) => {
    set({ isLoadingAssignments: true, error: null });
    try {
      const data = await service.getMyAssignments(params);
      set({ activeAssignments: data });
    } catch (e: any) {
      set({ error: e.message });
    } finally {
      set({ isLoadingAssignments: false });
    }
  },

  fetchCompletedAssignments: async (params) => {
    set({ isLoadingCompleted: true, error: null });
    try {
      const data = await service.getCompletedAssignments(params);
      set({ completedAssignments: data });
    } catch (e: any) {
      set({ error: e.message });
    } finally {
      set({ isLoadingCompleted: false });
    }
  },

  // ✅ UPDATED
  fetchLeaderboard: async (params) => {
    set({ isLoadingLeaderboard: true, error: null });
    try {
      const period = params.period || get().leaderboardPeriod;
      const data = await service.getLeaderboard(params, period);
      set({ leaderboard: data });
    } catch (e: any) {
      set({ error: e.message });
    } finally {
      set({ isLoadingLeaderboard: false });
    }
  },

  updateProblemProgress: async (params) => {
    const updated = await service.updateProblemProgress(params);
    // Optimistically update in store
    set((state) => ({
      activeAssignments: state.activeAssignments.map((a) =>
        a.id === params.assignmentId
          ? {
              ...a,
              problems: a.problems.map((p) =>
                p.id === params.assignmentProblemId ? updated : p,
              ),
              completedProblems: a.problems.filter(
                (p) =>
                  (p.id === params.assignmentProblemId ? updated : p).status === 'completed',
              ).length,
            }
          : a,
      ),
    }));
    return updated;
  },

  raiseDoubt: async (params) => {
    const doubt = await service.raiseDoubt(params);
    set((state) => ({
      activeAssignments: state.activeAssignments.map((a) =>
        a.id === params.assignmentId
          ? {
              ...a,
              problems: a.problems.map((p) =>
                p.id === params.assignmentProblemId
                  ? { ...p, doubt, menteeStatus: 'discussion_needed' as MenteeStatus }
                  : p,
              ),
            }
          : a,
      ),
    }));
  },
}));