import { create } from 'zustand';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { AppSession } from '../types';

interface AuthStore {
  session: AppSession | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isBootstrapping: boolean;   // true on first app launch while checking storage

  setSession: (session: AppSession) => void;
  logout: () => Promise<void>;
  setLoading: (loading: boolean) => void;
  bootstrapAuth: () => Promise<void>;
}

export const useAuthStore = create<AuthStore>((set, get) => ({
  session: null,
  isAuthenticated: false,
  isLoading: false,
  isBootstrapping: true,

  setSession: (session) =>
    set({ session, isAuthenticated: true, isLoading: false }),

  logout: async () => {
    await AsyncStorage.removeItem('@access_token');
    await AsyncStorage.removeItem('@session');
    set({ session: null, isAuthenticated: false });
  },

  setLoading: (isLoading) => set({ isLoading }),

  bootstrapAuth: async () => {
    try {
      const raw = await AsyncStorage.getItem('@session');
      if (raw) {
        const session: AppSession = JSON.parse(raw);
        set({ session, isAuthenticated: true });
      }
    } catch {
      // corrupted storage — start fresh
    } finally {
      set({ isBootstrapping: false });
    }
  },
}));