import type { MentorProfile, Profile } from "@/types";
import { api } from "./api";

type ProfilePayload = {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  solved: number;
  joinedAt: string;
  bio?: string;
  github?: string;
  linkedin?: string;
};

function toMentorProfile(profile: ProfilePayload): MentorProfile {
  return {
    firstName: profile.firstName,
    lastName: profile.lastName,
    username: profile.username,
    email: profile.email,
    solved: profile.solved,
    joinedAt: profile.joinedAt,
    bio: profile.bio,
    github: profile.github,
    linkedin: profile.linkedin,
  };
}

export async function getMenteeProfile(username: string): Promise<Profile | null> {
  try {
    return await api.get<Profile>(`/v1/app/mentees/${username}/profile`);
  } catch {
    return null;
  }
}

export async function getMyProfile(): Promise<MentorProfile> {
  const profile = await api.get<ProfilePayload>("/v1/app/me/profile");
  return toMentorProfile(profile);
}

export async function getMentorProfile(): Promise<MentorProfile> {
  return getMyProfile();
}

export async function updateMyProfile(profile: {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  bio?: string;
  github?: string;
  linkedin?: string;
}): Promise<MentorProfile> {
  const updated = await api.patch<ProfilePayload, typeof profile>("/v1/app/me/profile", profile);
  return toMentorProfile(updated);
}

export async function saveMentorProfile(profile: {
  firstName: string;
  lastName: string;
  username: string;
  email: string;
  bio?: string;
  github?: string;
  linkedin?: string;
}): Promise<MentorProfile> {
  return updateMyProfile(profile);
}

export async function updateMenteeProfile(
  username: string,
  fields: { bio?: string; github?: string; linkedin?: string }
): Promise<MentorProfile> {
  const current = await getMyProfile();
  if (current.username !== username) {
    throw new Error("ACCESS_DENIED");
  }

  return updateMyProfile({
    firstName: current.firstName,
    lastName: current.lastName,
    username: current.username,
    email: current.email,
    bio: fields.bio ?? current.bio,
    github: fields.github ?? current.github,
    linkedin: fields.linkedin ?? current.linkedin,
  });
}

export async function updateMyPassword(currentPassword: string, newPassword: string): Promise<void> {
  await api.patch<Record<string, never>, { currentPassword: string; newPassword: string }>(
    "/v1/app/me/password",
    { currentPassword, newPassword }
  );
}

export async function updateMenteePassword(
  username: string,
  currentPassword: string,
  newPassword: string
): Promise<{ ok: boolean; error?: string }> {
  const current = await getMyProfile();
  if (current.username !== username) {
    return { ok: false, error: "ACCESS_DENIED" };
  }

  try {
    await updateMyPassword(currentPassword, newPassword);
    return { ok: true };
  } catch (error) {
    return {
      ok: false,
      error: error instanceof Error ? error.message : "PASSWORD_UPDATE_FAILED",
    };
  }
}

export async function updateMentorPassword(
  currentPassword: string,
  newPassword: string
): Promise<{ ok: boolean; error?: string }> {
  try {
    await updateMyPassword(currentPassword, newPassword);
    return { ok: true };
  } catch (error) {
    return {
      ok: false,
      error: error instanceof Error ? error.message : "PASSWORD_UPDATE_FAILED",
    };
  }
}
