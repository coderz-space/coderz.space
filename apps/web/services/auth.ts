import type { AuthResponse, LoginResult, MenteeRequest, MenteeSignupInput } from "@/types";
import { api } from "./api";
import { getAppContext } from "./appContext";

export async function registerMentee(input: MenteeSignupInput): Promise<MenteeRequest> {
  const response = await api.post<{
    requestId: string;
    status: "pending" | "approved" | "rejected";
    username: string;
    email: string;
  }, {
    firstName: string;
    lastName: string;
    username: string;
    email: string;
    password: string;
  }>("/v1/app/auth/mentee-signup", {
    firstName: input.firstName,
    lastName: input.lastName,
    username: input.username,
    email: input.email,
    password: input.password,
  });

  return {
    id: response.requestId,
    firstName: input.firstName,
    lastName: input.lastName,
    username: response.username,
    email: response.email,
    signedUpAt: new Date().toISOString(),
    status: response.status,
  };
}

export async function loginMenteeByEmail(email: string, password: string): Promise<LoginResult> {
  const auth = await api.rawPost<AuthResponse, { email: string; password: string }>("/v1/auth/login", {
    email,
    password,
  });

  const context = await getAppContext();
  return {
    auth: auth.data,
    context,
  };
}

export async function loginMentee(identifier: string, password: string): Promise<LoginResult> {
  return loginMenteeByEmail(identifier, password);
}

export async function logout(): Promise<void> {
  try {
    await api.post<Record<string, never>>("/v1/auth/logout", {});
  } catch {
    // Logout should not block the redirect path.
  }
}
