"use client";

import { useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { getLeaderboard } from "@/services";

type LeaderboardEntry = { username: string; firstName: string; lastName: string; solved: number };

const medals = ["🥇", "🥈", "🥉"];

const rankBg = [
  "border-yellow-500/40 bg-yellow-900/10",
  "border-gray-400/40 bg-gray-800/30",
  "border-orange-500/40 bg-orange-900/10",
];

export default function LeaderboardPage() {
  const { username } = useParams() as { username: string };
  const router = useRouter();
  const [board, setBoard] = useState<LeaderboardEntry[]>([]);

  useEffect(() => {
    const loadLeaderboard = async () => {
      const entries = await getLeaderboard();
      setBoard(entries);
    };
    loadLeaderboard();
  }, []);

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-purple-400 mb-2">Leaderboard</h1>
      <p className="text-sm text-gray-500 mb-8">Ranked by questions solved (completed + revision needed).</p>

      {board.length === 0 && (
        <p className="text-gray-500 text-sm">No approved mentees yet.</p>
      )}

      <ol className="flex flex-col gap-3">
        {board.map((entry, i) => {
          const isMe = entry.username === username;
          const topThree = i < 3;
          return (
            <li
              key={entry.username}
              className={`flex items-center justify-between gap-4 px-5 py-4 rounded-xl border ${
                isMe
                  ? "border-purple-500 bg-purple-900/20"
                  : topThree
                  ? rankBg[i]
                  : "border-gray-800 bg-gray-900/40"
              }`}
            >
              {/* Rank + name */}
              <div className="flex items-center gap-4 min-w-0">
                <span className="text-xl w-8 text-center shrink-0">
                  {topThree ? medals[i] : <span className="text-gray-500 text-sm font-bold">{i + 1}</span>}
                </span>
                <div className="min-w-0">
                  <p className={`font-semibold truncate ${isMe ? "text-white" : "text-gray-200"}`}>
                    {entry.firstName} {entry.lastName}
                    {isMe && <span className="ml-2 text-xs text-purple-300 font-normal">(you)</span>}
                  </p>
                  <p className="text-xs text-gray-500">@{entry.username}</p>
                </div>
              </div>

              {/* Solved count + View Profile */}
              <div className="flex items-center gap-3 shrink-0">
                <span className="text-purple-400 font-bold text-lg">{entry.solved}</span>
                <span className="text-xs text-gray-600">solved</span>
                <button
                  onClick={() => router.push(`/mentee-dashboard/${username}/profile/${entry.username}`)}
                  className="text-xs font-semibold px-3 py-1.5 rounded-lg bg-gray-800 hover:bg-purple-700 text-gray-300 hover:text-white transition-colors border border-gray-700 hover:border-purple-600"
                >
                  View Profile
                </button>
              </div>
            </li>
          );
        })}
      </ol>
    </div>
  );
}
