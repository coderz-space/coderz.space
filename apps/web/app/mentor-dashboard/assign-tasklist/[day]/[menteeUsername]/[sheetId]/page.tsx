"use client";

import { use, useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { getMenteeRequests, assignTaskToMentee } from "@/services";
import { SHEET_QUESTIONS, SHEET_NAMES, type SheetQuestion } from "@/app/mentor-dashboard/master-tasklist/questionsData";
import type { MenteeRequest } from "@/types";

const DIFFICULTY_COLOR: Record<string, string> = {
  easy: "text-green-400",
  medium: "text-yellow-400",
  hard: "text-red-400",
};

function capitalize(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

export default function AssignQuestionsPage({
  params,
}: {
  params: Promise<{ day: string; menteeUsername: string; sheetId: string }>;
}) {
  const { day, menteeUsername, sheetId } = use(params);
  const router = useRouter();

  const questions: SheetQuestion[] = SHEET_QUESTIONS[sheetId] ?? [];
  const sheetName = SHEET_NAMES[sheetId] ?? sheetId;

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

  const [selectedQuestions, setSelectedQuestions] = useState<Set<string>>(new Set());
  const [assigned, setAssigned] = useState(false);

  const toggleQuestion = (id: string) => {
    setSelectedQuestions((prev) => {
      const next = new Set(prev);
      next.has(id) ? next.delete(id) : next.add(id);
      return next;
    });
  };

  const selectAll = () => {
    if (selectedQuestions.size === questions.length) {
      setSelectedQuestions(new Set());
    } else {
      setSelectedQuestions(new Set(questions.map((q) => q.id)));
    }
  };

  const handleAssign = async () => {
    if (selectedQuestions.size === 0) return;
    const selectedQs = questions.filter((q) => selectedQuestions.has(q.id));
    await Promise.all(selectedQs.map((q) =>
      assignTaskToMentee(menteeUsername, {
        title: q.title,
        description: "",
        difficulty: q.difficulty,
        topic: q.topic,
      })
    ));
    setAssigned(true);
    setTimeout(() => {
      setAssigned(false);
      setSelectedQuestions(new Set());
    }, 2000);
  };

  return (
    <div className="max-w-3xl pb-28">
      <button
        onClick={() => router.back()}
        className="text-xs text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 mb-4 flex items-center gap-1 transition-colors"
      >
        ← Back
      </button>

      <h1 className="text-2xl font-bold text-purple-400 mb-1">{sheetName}</h1>
      <p className="text-gray-500 text-sm mb-6">
        Assigning to:{" "}
        <span className="text-gray-700 dark:text-gray-300 font-medium">
          {mentee ? `${mentee.firstName} ${mentee.lastName}` : menteeUsername}
        </span>{" "}
        · {capitalize(day)} · {questions.length} questions
      </p>

      {/* Select all toggle */}
      <div className="flex items-center justify-between mb-3">
        <span className="text-xs text-gray-500">
          {selectedQuestions.size} of {questions.length} selected
        </span>
        <button
          onClick={selectAll}
          className="text-xs font-semibold text-purple-400 hover:text-purple-300 transition-colors"
        >
          {selectedQuestions.size === questions.length ? "Deselect All" : "Select All"}
        </button>
      </div>

      {/* Questions list */}
      <div className="flex flex-col gap-3">
        {questions.map((q) => {
          const isSelected = selectedQuestions.has(q.id);
          return (
            <div
              key={q.id}
              className={`flex items-center justify-between gap-4 rounded-xl px-5 py-4 border transition-colors ${
                isSelected
                  ? "bg-purple-950/40 border-purple-700"
                  : "bg-white dark:bg-gray-900 border-gray-200 dark:border-gray-800"
              }`}
            >
              <div className="min-w-0">
                <p className="text-gray-800 dark:text-gray-100 font-medium text-sm truncate">
                  {q.title}
                </p>
                <p className="text-xs text-gray-500 mt-0.5">
                  {q.topic} ·{" "}
                  <span className={DIFFICULTY_COLOR[q.difficulty]}>{q.difficulty}</span>
                </p>
              </div>
              <button
                onClick={() => toggleQuestion(q.id)}
                className={`text-xs font-semibold px-4 py-1.5 rounded-lg shrink-0 transition-colors ${
                  isSelected
                    ? "bg-purple-600 hover:bg-purple-500 text-white"
                    : "bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 border border-gray-200 dark:border-gray-700"
                }`}
              >
                {isSelected ? "Selected ✓" : "Select"}
              </button>
            </div>
          );
        })}
      </div>

      {/* Sticky Assign button */}
      <div className="fixed bottom-6 right-6 z-30">
        <button
          onClick={handleAssign}
          disabled={selectedQuestions.size === 0 || assigned}
          className={`px-6 py-3 rounded-xl text-sm font-bold shadow-lg transition-all ${
            assigned
              ? "bg-green-600 text-white scale-95"
              : selectedQuestions.size > 0
              ? "bg-purple-600 hover:bg-purple-500 text-white hover:scale-105 active:scale-95"
              : "bg-gray-300 dark:bg-gray-700 text-gray-500 cursor-not-allowed"
          }`}
        >
          {assigned
            ? "Assigned ✓"
            : `Assign ${selectedQuestions.size > 0 ? `(${selectedQuestions.size})` : ""} to ${
                mentee ? mentee.firstName : menteeUsername
              }`}
        </button>
      </div>
    </div>
  );
}
