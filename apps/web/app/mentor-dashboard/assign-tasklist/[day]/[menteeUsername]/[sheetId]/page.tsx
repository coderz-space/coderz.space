"use client";

import { use, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { createAssignments, getMenteeRequests, getSheets } from "@/services/mentor";
import type { MenteeRequest, Sheet, SheetQuestion } from "@/types";

const DIFFICULTY_COLOR: Record<SheetQuestion["difficulty"], string> = {
  easy: "text-green-400",
  medium: "text-yellow-400",
  hard: "text-red-400",
};

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export default function AssignQuestionsPage({
  params,
}: {
  params: Promise<{ day: string; menteeUsername: string; sheetId: string }>;
}) {
  const { day, menteeUsername, sheetId } = use(params);
  const router = useRouter();
  const [mentee, setMentee] = useState<MenteeRequest | null>(null);
  const [sheet, setSheet] = useState<Sheet | null>(null);
  const [selectedQuestions, setSelectedQuestions] = useState<Set<string>>(new Set());
  const [assigned, setAssigned] = useState(false);

  useEffect(() => {
    const loadPage = async () => {
      const [requests, sheets] = await Promise.all([getMenteeRequests(), getSheets()]);
      setMentee(requests.find((request) => request.username === menteeUsername) ?? null);
      setSheet(sheets.find((item) => item.key === sheetId) ?? null);
    };

    void loadPage();
  }, [menteeUsername, sheetId]);

  const questions = sheet?.questions ?? [];

  const toggleQuestion = (id: string) => {
    setSelectedQuestions((current) => {
      const next = new Set(current);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  };

  const allSelected = useMemo(
    () => questions.length > 0 && selectedQuestions.size === questions.length,
    [questions.length, selectedQuestions.size]
  );

  const selectAll = () => {
    if (allSelected) {
      setSelectedQuestions(new Set());
      return;
    }

    setSelectedQuestions(new Set(questions.map((question) => question.id)));
  };

  const handleAssign = async () => {
    if (selectedQuestions.size === 0) {
      return;
    }

    await createAssignments({
      day,
      menteeUsernames: [menteeUsername],
      sheetKey: sheetId as Sheet["key"],
      questionIds: Array.from(selectedQuestions),
    });

    setAssigned(true);
    window.setTimeout(() => {
      setAssigned(false);
      setSelectedQuestions(new Set());
    }, 2000);
  };

  return (
    <div className="max-w-3xl pb-28">
      <button onClick={() => router.back()} className="mb-4 flex items-center gap-1 text-xs text-gray-500 transition-colors hover:text-gray-700 dark:hover:text-gray-300">
        ← Back
      </button>

      <h1 className="mb-1 text-2xl font-bold text-purple-400">{sheet?.name ?? sheetId}</h1>
      <p className="mb-6 text-sm text-gray-500">
        Assigning to:{" "}
        <span className="font-medium text-gray-700 dark:text-gray-300">
          {mentee ? `${mentee.firstName} ${mentee.lastName}` : menteeUsername}
        </span>{" "}
        · {capitalize(day)} · {questions.length} questions
      </p>

      <div className="mb-3 flex items-center justify-between">
        <span className="text-xs text-gray-500">
          {selectedQuestions.size} of {questions.length} selected
        </span>
        <button onClick={selectAll} className="text-xs font-semibold text-purple-400 transition-colors hover:text-purple-300">
          {allSelected ? "Deselect All" : "Select All"}
        </button>
      </div>

      <div className="flex flex-col gap-3">
        {questions.map((question) => {
          const isSelected = selectedQuestions.has(question.id);
          return (
            <div
              key={question.id}
              className={`flex items-center justify-between gap-4 rounded-xl border px-5 py-4 transition-colors ${
                isSelected
                  ? "border-purple-700 bg-purple-950/40"
                  : "border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900"
              }`}
            >
              <div className="min-w-0">
                <p className="truncate text-sm font-medium text-gray-800 dark:text-gray-100">{question.title}</p>
                <p className="mt-0.5 text-xs text-gray-500">
                  {question.topic} · <span className={DIFFICULTY_COLOR[question.difficulty]}>{question.difficulty}</span>
                </p>
              </div>
              <button
                onClick={() => toggleQuestion(question.id)}
                className={`shrink-0 rounded-lg px-4 py-1.5 text-xs font-semibold transition-colors ${
                  isSelected
                    ? "bg-purple-600 text-white hover:bg-purple-500"
                    : "border border-gray-200 bg-gray-100 text-gray-700 hover:bg-gray-200 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
                }`}
              >
                {isSelected ? "Selected ✓" : "Select"}
              </button>
            </div>
          );
        })}
      </div>

      <div className="fixed bottom-6 right-6 z-30">
        <button
          onClick={() => void handleAssign()}
          disabled={selectedQuestions.size === 0 || assigned}
          className={`rounded-xl px-6 py-3 text-sm font-bold shadow-lg transition-all ${
            assigned
              ? "scale-95 bg-green-600 text-white"
              : selectedQuestions.size > 0
                ? "bg-purple-600 text-white hover:scale-105 hover:bg-purple-500 active:scale-95"
                : "cursor-not-allowed bg-gray-300 text-gray-500 dark:bg-gray-700"
          }`}
        >
          {assigned
            ? "Assigned ✓"
            : `Assign${selectedQuestions.size > 0 ? ` (${selectedQuestions.size})` : ""} to ${
                mentee ? mentee.firstName : menteeUsername
              }`}
        </button>
      </div>
    </div>
  );
}
