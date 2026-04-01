"use client";

import { use, useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { createAssignments, getMenteeRequests, getSheets } from "@/services/mentor";
import type { MenteeRequest, Sheet, SheetQuestion } from "@/types";

const DIFFICULTY_COLOR: Record<SheetQuestion["difficulty"], string> = {
  easy: "text-green-400",
  medium: "text-yellow-400",
  hard: "text-red-400",
};

export default function SheetAssignPage({ params }: { params: Promise<{ sheetId: string }> }) {
  const { sheetId } = use(params);
  const router = useRouter();
  const dropdownRef = useRef<HTMLDivElement>(null);
  const [sheet, setSheet] = useState<Sheet | null>(null);
  const [approvedMentees, setApprovedMentees] = useState<MenteeRequest[]>([]);
  const [selectedQuestions, setSelectedQuestions] = useState<Set<string>>(new Set());
  const [selectedMentees, setSelectedMentees] = useState<MenteeRequest[]>([]);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [assigned, setAssigned] = useState(false);

  useEffect(() => {
    const loadPage = async () => {
      const [requests, sheets] = await Promise.all([getMenteeRequests(), getSheets()]);
      setApprovedMentees(requests.filter((request) => request.status === "approved"));
      setSheet(sheets.find((item) => item.key === sheetId) ?? null);
    };

    void loadPage();
  }, [sheetId]);

  useEffect(() => {
    const handleClick = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setDropdownOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  const questions = sheet?.questions ?? [];

  const selectAllQuestions = () => {
    if (selectedQuestions.size === questions.length) {
      setSelectedQuestions(new Set());
      return;
    }

    setSelectedQuestions(new Set(questions.map((question) => question.id)));
  };

  const addMentee = (mentee: MenteeRequest) => {
    setSelectedMentees((current) => {
      if (current.some((item) => item.id === mentee.id)) {
        return current;
      }
      return [...current, mentee];
    });
    setDropdownOpen(false);
  };

  const removeMentee = (id: string) => {
    setSelectedMentees((current) => current.filter((mentee) => mentee.id !== id));
  };

  const selectAllMentees = () => {
    setSelectedMentees(approvedMentees);
    setDropdownOpen(false);
  };

  const handleAssign = async () => {
    if (!sheet || selectedQuestions.size === 0 || selectedMentees.length === 0) {
      return;
    }

    await createAssignments({
      menteeUsernames: selectedMentees.map((mentee) => mentee.username),
      sheetKey: sheet.key,
      questionIds: Array.from(selectedQuestions),
    });

    setAssigned(true);
    window.setTimeout(() => {
      setAssigned(false);
      setSelectedQuestions(new Set());
      setSelectedMentees([]);
    }, 2000);
  };

  const canAssign = selectedQuestions.size > 0 && selectedMentees.length > 0;
  const selectedMenteeIds = useMemo(
    () => new Set(selectedMentees.map((mentee) => mentee.id)),
    [selectedMentees]
  );

  return (
    <div className="max-w-3xl pb-28">
      <button onClick={() => router.back()} className="mb-4 flex items-center gap-1 text-xs text-gray-500 transition-colors hover:text-gray-300">
        ← Back
      </button>
      <h1 className="mb-1 text-2xl font-bold text-purple-400">{sheet?.name ?? sheetId}</h1>
      <p className="mb-6 text-sm text-gray-500">{questions.length} questions</p>

      <div className="mb-6 flex flex-wrap items-start gap-3">
        <div className="flex flex-1 flex-wrap gap-2">
          {selectedMentees.map((mentee) => (
            <span
              key={mentee.id}
              className="flex items-center gap-1.5 rounded-full border border-purple-700 bg-purple-900/60 px-3 py-1.5 text-xs font-medium text-purple-200"
            >
              {mentee.firstName} {mentee.lastName}
              <button
                onClick={() => removeMentee(mentee.id)}
                className="leading-none text-purple-400 transition-colors hover:text-white"
                aria-label={`Remove ${mentee.firstName}`}
              >
                ×
              </button>
            </span>
          ))}
          {selectedMentees.length === 0 ? <span className="py-1.5 text-xs text-gray-600">No mentees selected</span> : null}
        </div>

        <div className="relative shrink-0" ref={dropdownRef}>
          <button
            onClick={() => setDropdownOpen((value) => !value)}
            className="flex items-center gap-2 rounded-lg border border-gray-700 bg-gray-800 px-4 py-2 text-sm font-semibold text-gray-200 transition-colors hover:bg-gray-700"
          >
            Select Mentees
            <span className="text-xs text-gray-400">{dropdownOpen ? "▲" : "▼"}</span>
          </button>

          {dropdownOpen ? (
            <div className="absolute right-0 top-full z-20 mt-2 w-64 overflow-hidden rounded-xl border border-gray-700 bg-gray-900 shadow-xl">
              <button
                onClick={selectAllMentees}
                className="w-full border-b border-gray-800 px-4 py-2.5 text-left text-sm font-semibold text-purple-400 transition-colors hover:bg-gray-800"
              >
                All mentees
              </button>
              {approvedMentees.length === 0 ? <p className="px-4 py-3 text-xs text-gray-500">No approved mentees</p> : null}
              {approvedMentees.map((mentee) => (
                <button
                  key={mentee.id}
                  onClick={() => addMentee(mentee)}
                  className="flex w-full items-center justify-between px-4 py-2.5 text-left text-sm text-gray-200 transition-colors hover:bg-gray-800"
                >
                  <span>
                    {mentee.firstName} {mentee.lastName}
                  </span>
                  {selectedMenteeIds.has(mentee.id) ? <span className="text-xs text-purple-400">✓</span> : null}
                </button>
              ))}
            </div>
          ) : null}
        </div>
      </div>

      <div className="mb-3 flex items-center justify-between">
        <span className="text-xs text-gray-500">
          {selectedQuestions.size} of {questions.length} selected
        </span>
        <button onClick={selectAllQuestions} className="text-xs font-semibold text-purple-400 transition-colors hover:text-purple-300">
          {selectedQuestions.size === questions.length ? "Deselect All" : "Select All"}
        </button>
      </div>

      <div className="flex flex-col gap-3">
        {questions.map((question) => {
          const isSelected = selectedQuestions.has(question.id);
          return (
            <div
              key={question.id}
              className={`flex items-center justify-between gap-4 rounded-xl border px-5 py-4 transition-colors ${
                isSelected ? "border-purple-700 bg-purple-950/40" : "border-gray-800 bg-gray-900"
              }`}
            >
              <div className="min-w-0">
                <p className="truncate text-sm font-medium text-gray-100">{question.title}</p>
                <p className="mt-0.5 text-xs text-gray-500">
                  {question.topic} · <span className={DIFFICULTY_COLOR[question.difficulty]}>{question.difficulty}</span>
                </p>
              </div>
              <button
                onClick={() =>
                  setSelectedQuestions((current) => {
                    const next = new Set(current);
                    if (next.has(question.id)) {
                      next.delete(question.id);
                    } else {
                      next.add(question.id);
                    }
                    return next;
                  })
                }
                className={`shrink-0 rounded-lg px-4 py-1.5 text-xs font-semibold transition-colors ${
                  isSelected
                    ? "bg-purple-600 text-white hover:bg-purple-500"
                    : "border border-gray-700 bg-gray-800 text-gray-300 hover:bg-gray-700"
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
          disabled={!canAssign || assigned}
          className={`rounded-xl px-6 py-3 text-sm font-bold shadow-lg transition-all ${
            assigned
              ? "scale-95 bg-green-600 text-white"
              : canAssign
                ? "bg-purple-600 text-white hover:scale-105 hover:bg-purple-500 active:scale-95"
                : "cursor-not-allowed bg-gray-700 text-gray-500"
          }`}
        >
          {assigned
            ? "Assigned ✓"
            : `Assign${selectedQuestions.size > 0 ? ` (${selectedQuestions.size})` : ""} to ${
                selectedMentees.length > 0 ? `${selectedMentees.length} mentee${selectedMentees.length > 1 ? "s" : ""}` : "mentees"
              }`}
        </button>
      </div>
    </div>
  );
}
