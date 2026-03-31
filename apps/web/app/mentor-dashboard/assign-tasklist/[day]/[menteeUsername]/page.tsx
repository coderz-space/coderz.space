"use client";

import { use, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { getMenteeRequests } from "@/services";
import type { MenteeRequest } from "@/types";

const TASKLISTS = [
  { id: "gfg-dsa-360", name: "GFG DSA 360" },
  { id: "strivers-dsa-sheet", name: "Striver's DSA Sheet" },
];

function capitalize(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

export default function MenteeSheetSelectorPage({
  params,
}: {
  params: Promise<{ day: string; menteeUsername: string }>;
}) {
  const { day, menteeUsername } = use(params);
  const router = useRouter();

  const [mentee, setMentee] = useState<MenteeRequest | null>(null);

  useEffect(() => {
    const loadMentee = async () => {
      const requests = await getMenteeRequests();
      const found = requests.find(
        (r) => r.username === menteeUsername && r.status === "approved"
      );
      setMentee(found || null);
    };
    loadMentee();
  }, [menteeUsername]);

  return (
    <div className="max-w-2xl">
      <button
        onClick={() => router.back()}
        className="text-xs text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 mb-4 flex items-center gap-1 transition-colors"
      >
        ← Back
      </button>

      <h1 className="text-2xl font-bold text-purple-400 mb-1">Assign Questions</h1>
      <p className="text-gray-500 text-sm mb-8">
        To:{" "}
        <span className="text-gray-700 dark:text-gray-300 font-medium">
          {mentee ? `${mentee.firstName} ${mentee.lastName}` : menteeUsername}
        </span>{" "}
        · {capitalize(day)}
      </p>

      <div className="flex flex-col gap-4">
        {TASKLISTS.map((list) => (
          <div
            key={list.id}
            className="flex items-center justify-between gap-4 border border-gray-200 dark:border-gray-700 rounded-xl px-6 py-5 bg-white dark:bg-gray-900"
          >
            <span className="text-gray-800 dark:text-gray-100 font-semibold text-base">
              {list.name}
            </span>
            <button
              onClick={() =>
                router.push(
                  `/mentor-dashboard/assign-tasklist/${day}/${menteeUsername}/${list.id}`
                )
              }
              className="text-sm font-semibold px-4 py-2 rounded-lg bg-purple-700 hover:bg-purple-600 text-white transition-colors shrink-0"
            >
              Assign Questions
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
