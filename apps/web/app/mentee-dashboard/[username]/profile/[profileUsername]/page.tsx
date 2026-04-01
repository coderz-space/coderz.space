"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import PublicMenteeProfile from "@/components/dashboard/PublicMenteeProfile";
import { getMenteeQuestions } from "@/services/mentee";
import { getMenteeProfile } from "@/services/profile";
import type { Profile, Question } from "@/types";

export default function MenteeProfilePage() {
  const { username, profileUsername } = useParams() as {
    username: string;
    profileUsername: string;
  };
  const router = useRouter();
  const [profile, setProfile] = useState<Profile | null>(null);
  const [completedQuestions, setCompletedQuestions] = useState<Question[]>([]);

  useEffect(() => {
    const loadProfile = async () => {
      const [nextProfile, questions] = await Promise.all([
        getMenteeProfile(profileUsername),
        getMenteeQuestions(profileUsername),
      ]);

      setProfile(nextProfile);
      setCompletedQuestions(questions.filter((question) => question.status === "completed"));
    };

    void loadProfile();
  }, [profileUsername]);

  if (!profile) {
    return <p className="text-sm text-gray-500 dark:text-gray-400">Profile not found.</p>;
  }

  return (
    <PublicMenteeProfile
      viewerUsername={username}
      profile={profile}
      completedQuestions={completedQuestions}
      onBack={() => router.back()}
      onOpenQuestion={(questionId) =>
        router.push(`/mentee-dashboard/${username}/completed/${questionId}?owner=${profileUsername}`)
      }
    />
  );
}
