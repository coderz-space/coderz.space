"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import LeaderboardList from "@/components/dashboard/LeaderboardList";
import { getLeaderboard } from "@/services/leaderboard";
import type { LeaderboardEntry } from "@/types";

export default function MentorLeaderboardPage() {
  const router = useRouter();
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);

  useEffect(() => {
    const loadLeaderboard = async () => {
      const leaderboard = await getLeaderboard();
      setEntries(leaderboard);
    };

    void loadLeaderboard();
  }, []);

  return (
    <div className="max-w-3xl">
      <h1 className="mb-2 text-2xl font-bold text-gray-900 dark:text-white">Leaderboard</h1>
      <p className="mb-8 text-sm text-gray-500 dark:text-gray-400">
        Ranked by completed Algo Buddy assignments.
      </p>

      <LeaderboardList
        entries={entries}
        onViewProfile={(username) => router.push(`/mentor-dashboard/profile/${username}`)}
      />
    </div>
  );
}
