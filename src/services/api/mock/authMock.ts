import AsyncStorage from '@react-native-async-storage/async-storage';
import { IAuthService } from '../interfaces';
import { LoginPayload, LoginResponse } from '../../../types';
import { MOCK_USER_MENTEE, MOCK_USER_MENTOR } from './_mockData';

const delay = (ms: number) => new Promise<void>((r) => setTimeout(() => r(), ms));

// Test credentials:
// mentor@test.com / any password  →  mentor session
// mentee@test.com / any password  →  mentee session
export const authMock: IAuthService = {
  async login({ email, password }: LoginPayload): Promise<LoginResponse> {
    await delay(1200);

    if (!email || !password) {
      throw new Error('Email and password are required');
    }

    const isMentor = email.toLowerCase().includes('mentor') ||
      email === 'priya@coderz.space';

    const mockToken = `mock_jwt_${Date.now()}`;
    await AsyncStorage.setItem('@access_token', mockToken);

    if (isMentor) {
      return {
        user: MOCK_USER_MENTOR,
        accessToken: mockToken,
        orgRole: 'mentor',
        bootcampRole: 'mentor',
        activeOrgId: 'org-1',
        activeBootcampId: 'bootcamp-1',
      };
    }

    return {
      user: MOCK_USER_MENTEE,
      accessToken: mockToken,
      orgRole: 'mentee',
      bootcampRole: 'mentee',
      activeOrgId: 'org-1',
      activeBootcampId: 'bootcamp-1',
    };
  },

  async logout(): Promise<void> {
    await delay(300);
    await AsyncStorage.removeItem('@access_token');
  },

  async refreshToken() {
    await delay(500);
    return { accessToken: `mock_jwt_refreshed_${Date.now()}` };
  },

  async getMe(): Promise<LoginResponse> {
    await delay(800);
    const token = await AsyncStorage.getItem('@access_token');
    if (!token) throw new Error('No session');
    return {
      user: MOCK_USER_MENTEE,
      accessToken: token,
      orgRole: 'mentee',
      bootcampRole: 'mentee',
      activeOrgId: 'org-1',
      activeBootcampId: 'bootcamp-1',
    };
  },

  // ✅ NEW METHOD
  async changePassword({
    currentPassword,
    newPassword,
  }: {
    currentPassword: string;
    newPassword: string;
  }): Promise<void> {
    await delay(1000);

    // Retrieve stored password (mock)
    let storedPassword = await AsyncStorage.getItem('@mock_password');

    if (!storedPassword) {
      // First time default
      storedPassword = 'pass123';
      await AsyncStorage.setItem('@mock_password', storedPassword);
    }

    // Validate current password
    if (currentPassword !== storedPassword) {
      throw new Error('Current password is incorrect');
    }

    // Validate new password strength
    if (newPassword.length < 8) {
      throw new Error('New password must be at least 8 characters');
    }

    if (!/[A-Za-z]/.test(newPassword) || !/[0-9]/.test(newPassword)) {
      throw new Error('New password must contain at least one letter and one number');
    }

    // Update password
    await AsyncStorage.setItem('@mock_password', newPassword);

    console.log('[Mock] Password changed successfully');
  },
};