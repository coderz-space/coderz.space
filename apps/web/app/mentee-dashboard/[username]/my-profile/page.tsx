"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import AccountProfileEditor from "@/components/dashboard/AccountProfileEditor";
import { getMyProfile, updateMenteePassword, updateMyProfile } from "@/services/profile";
import type { MentorProfile } from "@/types";

export default function MenteeMyProfilePage() {
  const { username } = useParams() as { username: string };
  const router = useRouter();
  const [profile, setProfile] = useState<MentorProfile | null>(null);

  useEffect(() => {
    const loadProfile = async () => {
      const currentProfile = await getMyProfile();
      setProfile(currentProfile);

      if (currentProfile.username !== username) {
        router.replace(`/mentee-dashboard/${currentProfile.username}/my-profile`);
      }
    };

    void loadProfile();
  }, [router, username]);

  return (
    <AccountProfileEditor
      title="My Profile"
      roleLabel="Mentee"
      profile={profile}
      stats={[
        { label: "Solved", value: profile?.solved ?? 0 },
        { label: "Role", value: "Mentee" },
      ]}
      onSave={async (nextProfile) => {
        const updated = await updateMyProfile(nextProfile);
        setProfile(updated);

        if (updated.username !== username) {
          router.replace(`/mentee-dashboard/${updated.username}/my-profile`);
        }

        return updated;
      }}
      onUpdatePassword={(currentPassword, newPassword) =>
        updateMenteePassword(profile?.username ?? username, currentPassword, newPassword)
      }
    />
  );
}
