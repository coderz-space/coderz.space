"use client";

import { useEffect, useState } from "react";
import AccountProfileEditor from "@/components/dashboard/AccountProfileEditor";
import { getLeaderboard } from "@/services/leaderboard";
import { getMentorProfile, saveMentorProfile, updateMentorPassword } from "@/services/profile";
import type { MentorProfile } from "@/types";

export default function MentorMyProfilePage() {
  const [profile, setProfile] = useState<MentorProfile | null>(null);
  const [activeMentees, setActiveMentees] = useState(0);

  useEffect(() => {
    const loadPage = async () => {
      const [currentProfile, leaderboard] = await Promise.all([
        getMentorProfile(),
        getLeaderboard(),
      ]);

      setProfile(currentProfile);
      setActiveMentees(leaderboard.length);
    };

    void loadPage();
  }, []);

  return (
    <AccountProfileEditor
      title="My Profile"
      roleLabel="Mentor"
      profile={profile}
      stats={[
        { label: "Active Mentees", value: activeMentees },
        { label: "Role", value: "Mentor" },
      ]}
      onSave={async (nextProfile) => {
        const updated = await saveMentorProfile(nextProfile);
        setProfile(updated);
        return updated;
      }}
      onUpdatePassword={updateMentorPassword}
    />
  );
}
