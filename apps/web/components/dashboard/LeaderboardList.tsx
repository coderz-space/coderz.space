"use client";

import type { LeaderboardEntry } from "@/types";

type LeaderboardListProps = {
  entries: LeaderboardEntry[];
  currentUsername?: string;
  onViewProfile: (username: string) => void;
};

const rankClasses = [
  "border-amber-300 bg-amber-50 dark:border-amber-900/60 dark:bg-amber-950/20",
  "border-slate-300 bg-slate-50 dark:border-slate-800 dark:bg-slate-900/40",
  "border-orange-300 bg-orange-50 dark:border-orange-900/60 dark:bg-orange-950/20",
];

export default function LeaderboardList({
  entries,
  currentUsername,
  onViewProfile,
}: LeaderboardListProps) {
  if (!entries.length) {
    return <p className="text-sm text-gray-500 dark:text-gray-400">No approved mentees yet.</p>;
  }

  return (
    <ol className="flex flex-col gap-3">
      {entries.map((entry, index) => {
        const isCurrentUser = currentUsername === entry.username;
        const isTopThree = index < 3;

        return (
          <li
            key={entry.username}
            className={`flex items-center justify-between gap-4 rounded-2xl border px-5 py-4 ${
              isCurrentUser
                ? "border-purple-400 bg-purple-50 dark:border-purple-700 dark:bg-purple-950/30"
                : isTopThree
                  ? rankClasses[index]
                  : "border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900/50"
            }`}
          >
            <div className="flex min-w-0 items-center gap-4">
              <span
                className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-sm font-semibold ${
                  isTopThree
                    ? "bg-white/70 text-gray-900 dark:bg-black/20 dark:text-white"
                    : "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300"
                }`}
              >
                #{index + 1}
              </span>

              <div className="min-w-0">
                <p className="truncate font-semibold text-gray-900 dark:text-gray-100">
                  {entry.firstName} {entry.lastName}
                  {isCurrentUser ? (
                    <span className="ml-2 text-xs font-normal text-purple-600 dark:text-purple-300">(you)</span>
                  ) : null}
                </p>
                <p className="text-xs text-gray-500 dark:text-gray-400">@{entry.username}</p>
              </div>
            </div>

            <div className="flex shrink-0 items-center gap-3">
              <div className="text-right">
                <p className="text-lg font-bold text-purple-600 dark:text-purple-300">{entry.solved}</p>
                <p className="text-xs text-gray-500 dark:text-gray-400">solved</p>
              </div>
              <button
                type="button"
                onClick={() => onViewProfile(entry.username)}
                className="rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-semibold text-gray-700 transition-colors hover:border-purple-500 hover:bg-purple-50 hover:text-purple-700 dark:border-gray-700 dark:text-gray-200 dark:hover:border-purple-600 dark:hover:bg-purple-950/30 dark:hover:text-white"
              >
                View Profile
              </button>
            </div>
          </li>
        );
      })}
    </ol>
  );
}
