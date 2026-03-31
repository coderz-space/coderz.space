"use client";

import { useRouter } from "next/navigation";

const DAYS = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"];

export default function AssignTasklistPage() {
  const router = useRouter();

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-purple-400 mb-2">Assign Tasklist</h1>
      <p className="text-gray-500 dark:text-gray-500 text-sm mb-8">
        Select a day to manage mentees and assign questions.
      </p>

      <div className="flex flex-col gap-3">
        {DAYS.map((day) => (
          <button
            key={day}
            onClick={() => router.push(`/mentor-dashboard/assign-tasklist/${day.toLowerCase()}`)}
            className="flex items-center justify-between px-6 py-4 rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 hover:border-purple-400 dark:hover:border-purple-600 hover:bg-purple-50 dark:hover:bg-purple-950/30 transition-colors text-left group"
          >
            <span className="font-semibold text-gray-800 dark:text-gray-100 group-hover:text-purple-600 dark:group-hover:text-purple-400 transition-colors">
              {day}
            </span>
            <span className="text-gray-400 dark:text-gray-600 group-hover:text-purple-400 transition-colors text-sm">
              →
            </span>
          </button>
        ))}
      </div>
    </div>
  );
}
