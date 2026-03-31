"use client";

import { useRouter } from "next/navigation";

const TASKLISTS = [
  { id: "gfg-dsa-360", name: "GFG DSA 360" },
  { id: "strivers-dsa-sheet", name: "Striver's DSA Sheet" },
];

export default function MasterTasklistPage() {
  const router = useRouter();

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-purple-400 mb-2">Master Tasklist</h1>
      <p className="text-gray-500 text-sm mb-8">Select a sheet to assign questions to mentees.</p>

      <div className="flex flex-col gap-4">
        {TASKLISTS.map((list) => (
          <div
            key={list.id}
            className="flex items-center justify-between gap-4 border border-gray-700 rounded-xl px-6 py-5 bg-gray-900"
          >
            <span className="text-gray-100 font-semibold text-base">{list.name}</span>
            <button
              onClick={() => router.push(`/mentor-dashboard/master-tasklist/${list.id}`)}
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
