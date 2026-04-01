"use client";

import { use, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getMenteeRequests, getSheets } from "@/services/mentor";
import type { MenteeRequest, Sheet } from "@/types";

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export default function MenteeSheetSelectorPage({
  params,
}: {
  params: Promise<{ day: string; menteeUsername: string }>;
}) {
  const { day, menteeUsername } = use(params);
  const router = useRouter();
  const [mentee, setMentee] = useState<MenteeRequest | null>(null);
  const [sheets, setSheets] = useState<Sheet[]>([]);

  useEffect(() => {
    const loadPage = async () => {
      const [requests, availableSheets] = await Promise.all([getMenteeRequests(), getSheets()]);
      setMentee(requests.find((request) => request.username === menteeUsername) ?? null);
      setSheets(availableSheets);
    };

    void loadPage();
  }, [menteeUsername]);

  return (
    <div className="max-w-2xl">
      <button onClick={() => router.back()} className="mb-4 flex items-center gap-1 text-xs text-gray-500 transition-colors hover:text-gray-700 dark:hover:text-gray-300">
        ← Back
      </button>

      <h1 className="mb-1 text-2xl font-bold text-purple-400">Assign Questions</h1>
      <p className="mb-8 text-sm text-gray-500">
        To:{" "}
        <span className="font-medium text-gray-700 dark:text-gray-300">
          {mentee ? `${mentee.firstName} ${mentee.lastName}` : menteeUsername}
        </span>{" "}
        · {capitalize(day)}
      </p>

      <div className="flex flex-col gap-4">
        {sheets.map((sheet) => (
          <div
            key={sheet.key}
            className="flex items-center justify-between gap-4 rounded-xl border border-gray-200 bg-white px-6 py-5 dark:border-gray-700 dark:bg-gray-900"
          >
            <div>
              <span className="text-base font-semibold text-gray-800 dark:text-gray-100">{sheet.name}</span>
              <p className="mt-1 text-xs text-gray-500">{sheet.questions.length} questions</p>
            </div>
            <button
              onClick={() => router.push(`/mentor-dashboard/assign-tasklist/${day}/${menteeUsername}/${sheet.key}`)}
              className="shrink-0 rounded-lg bg-purple-700 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-purple-600"
            >
              Assign Questions
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
