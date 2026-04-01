"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getSheets } from "@/services/mentor";
import type { Sheet } from "@/types";

export default function MasterTasklistPage() {
  const router = useRouter();
  const [sheets, setSheets] = useState<Sheet[]>([]);

  useEffect(() => {
    void getSheets().then(setSheets);
  }, []);

  return (
    <div className="max-w-2xl">
      <h1 className="mb-2 text-2xl font-bold text-purple-400">Master Tasklist</h1>
      <p className="mb-8 text-sm text-gray-500">Select a sheet to assign questions to mentees.</p>

      <div className="flex flex-col gap-4">
        {sheets.map((sheet) => (
          <div
            key={sheet.key}
            className="flex items-center justify-between gap-4 rounded-xl border border-gray-700 bg-gray-900 px-6 py-5"
          >
            <div>
              <span className="text-base font-semibold text-gray-100">{sheet.name}</span>
              <p className="mt-1 text-xs text-gray-500">{sheet.questions.length} questions</p>
            </div>
            <button
              onClick={() => router.push(`/mentor-dashboard/master-tasklist/${sheet.key}`)}
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
