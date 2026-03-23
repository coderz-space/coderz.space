import apiClient, { ApiRoutes } from '../apiClient';
import { IAuthServiceV2 } from '../interfaces';
import { LoginPayload, LoginResponse, SignupPayload, SignupResponse } from '../../../types';
import AsyncStorage from '@react-native-async-storage/async-storage';

export const authLive: IAuthServiceV2 = {
  async signup(payload: SignupPayload): Promise<SignupResponse> {
    const { data } = await apiClient.post(ApiRoutes.signup, payload);
    return data.data;
  },

  async login(payload: LoginPayload): Promise<LoginResponse> {
    const { data } = await apiClient.post(ApiRoutes.login, payload);
    const { accessToken, user } = data.data;
    await AsyncStorage.setItem('@access_token', accessToken);
    // NOTE: orgRole & bootcampRole require /auth/me or org context endpoint
    // after login, call getMe() and then fetch org membership
    return { user, accessToken, orgRole: 'mentee' };
  },

  async logout(): Promise<void> {
    await apiClient.post(ApiRoutes.logout);
    await AsyncStorage.removeItem('@access_token');
    await AsyncStorage.removeItem('@session');
  },

  async refreshToken() {
    const { data } = await apiClient.post(ApiRoutes.refresh);
    const { accessToken } = data.data;
    await AsyncStorage.setItem('@access_token', accessToken);
    return { accessToken };
  },

  async getMe(): Promise<LoginResponse> {
    const { data } = await apiClient.get(ApiRoutes.me);
    return { user: data.data, accessToken: '', orgRole: 'mentee' };
  },

  async forgotPassword(email: string): Promise<void> {
    await apiClient.post('/auth/forgot-password', { email });
  },

  async resetPassword({ token, newPassword }): Promise<void> {
    await apiClient.post('/auth/reset-password', { token, newPassword });
  },
};