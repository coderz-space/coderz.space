"use client";

import { use, useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import { getMenteeRequests } from "@/services";
import { SHEET_QUESTIONS, SHEET_NAMES, type SheetQuestion } from "../questionsData";
import type { MenteeRequest } from "@/types";

const DIFFICULTY_COLOR: Record<string, string> = {
  easy: "text-green-400",
  medium: "text-yellow-400",
  hard: "text-red-400",
};

export default function SheetAssignPage({ params }: { params: Promise<{ sheetId: string }> }) {
  const { sheetId } = use(params);
  const router = useRouter();

  const questions: SheetQuestion[] = SHEET_QUESTIONS[sheetId] ?? [];
  const sheetName = SHEET_NAMES[sheetId] ?? sheetId;

  const [approvedMentees, setApprovedMentees] = useState<MenteeRequest[]>([]);
  const [selectedQuestions, setSelectedQuestions] = useState<Set<string>>(new Set());
  const [selectedMentees, setSelectedMentees] = useState<MenteeRequest[]>([]);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [assigned, setAssigned] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Load approved mentees
  useEffect(() => {
    const loadMentees = async () => {
      const requests = await getMenteeRequests();
      const approved = requests.filter((r) => r.status === "approved");
      setApprovedMentees(approved);
    };
    loadMentees();
  }, []);

  // Close dropdown on outside click
  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setDropdownOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  const toggleQuestion = (id: string) => {
    setSelectedQuestions((prev) => {
      const next = new Set(prev);
      next.has(id) ? next.delete(id) : next.add(id);
      return next;
    });
  };

  const selectAllQuestions = () => {
    if (selectedQuestions.size === questions.length) {
      setSelectedQuestions(new Set());
    } else {
      setSelectedQuestions(new Set(questions.map((q) => q.id)));
    }
  };

  const addMentee = (mentee: MenteeRequest) => {
    if (!selectedMentees.find((m) => m.id === mentee.id)) {
      setSelectedMentees((prev) => [...prev, mentee]);
    }
    setDropdownOpen(false);
  };

  const removeMentee = (id: string) => {
    setSelectedMentees((prev) => prev.filter((m) => m.id !== id));
  };

  const selectAllMentees = () => {
    setSelectedMentees(approvedMentees);
    setDropdownOpen(false);
  };

  const handleAssign = () => {
    if (selectedQuestions.size === 0 || selectedMentees.length === 0) return;
    // Replace with: POST /api/assign-questions { menteeIds, questionIds, sheetId }
    const selectedQs = questions.filter((q) => selectedQuestions.has(q.id));
    selectedMentees.forEach((mentee) => {
      const key = `coderz_assigned_tasks_${mentee.username}`;
      const existing = JSON.parse(
        (typeof window !== "undefined" && localStorage.getItem(key)) || "[]"
      );
      const newTasks = selectedQs.map((q) => ({
        id: crypto.randomUUID(),
        title: q.title,
        description: "",
        difficulty: q.difficulty,
        topic: q.topic,
        status: "pending",
        progressStatus: "not_started",
        assignedAt: new Date().toISOString(),
      }));
      localStorage.setItem(key, JSON.stringify([...existing, ...newTasks]));
    });
    setAssigned(true);
    setTimeout(() => {
      setAssigned(false);
      setSelectedQuestions(new Set());
      setSelectedMentees([]);
    }, 2000);
  };

  const canAssign = selectedQuestions.size > 0 && selectedMentees.length > 0;

  return (
    <div className="max-w-3xl pb-28">
      {/* Header */}
      <button
        onClick={() => router.back()}
        className="text-xs text-gray-500 hover:text-gray-300 mb-4 flex items-center gap-1 transition-colors"
      >
        ← Back
      </button>
      <h1 className="text-2xl font-bold text-purple-400 mb-1">{sheetName}</h1>
      <p className="text-gray-500 text-sm mb-6">{questions.length} questions</p>

      {/* Mentee selector row */}
      <div className="flex items-start gap-3 mb-6 flex-wrap">
        {/* Selected mentee chips */}
        <div className="flex flex-wrap gap-2 flex-1">
          {selectedMentees.map((m) => (
            <span
              key={m.id}
              className="flex items-center gap-1.5 bg-purple-900/60 border border-purple-700 text-purple-200 text-xs font-medium px-3 py-1.5 rounded-full"
            >
              {m.firstName} {m.lastName}
              <button
                onClick={() => removeMentee(m.id)}
                className="text-purple-400 hover:text-white transition-colors leading-none"
                aria-label={`Remove ${m.firstName}`}
              >
                ×
              </button>
            </span>
          ))}
          {selectedMentees.length === 0 && (
            <span className="text-xs text-gray-600 py-1.5">No mentees selected</span>
          )}
        </div>

        {/* Dropdown */}
        <div className="relative shrink-0" ref={dropdownRef}>
          <button
            onClick={() => setDropdownOpen((v) => !v)}
            className="text-sm font-semibold px-4 py-2 rounded-lg bg-gray-800 hover:bg-gray-700 border border-gray-700 text-gray-200 transition-colors flex items-center gap-2"
          >
            Select Mentees
            <span className="text-gray-400 text-xs">{dropdownOpen ? "▲" : "▼"}</span>
          </button>

          {dropdownOpen && (
            <div className="absolute right-0 top-full mt-2 w-64 bg-gray-900 border border-gray-700 rounded-xl shadow-xl z-20 overflow-hidden">
              {/* All button */}
              <button
                onClick={selectAllMentees}
                className="w-full text-left px-4 py-2.5 text-sm font-semibold text-purple-400 hover:bg-gray-800 border-b border-gray-800 transition-colors"
              >
                All mentees
              </button>
              {approvedMentees.length === 0 && (
                <p className="px-4 py-3 text-xs text-gray-500">No approved mentees</p>
              )}
              {approvedMentees.map((m) => (
                <button
                  key={m.id}
                  onClick={() => addMentee(m)}
                  className="w-full text-left px-4 py-2.5 text-sm text-gray-200 hover:bg-gray-800 transition-colors flex items-center justify-between"
                >
                  <span>{m.firstName} {m.lastName}</span>
                  {selectedMentees.find((s) => s.id === m.id) && (
                    <span className="text-purple-400 text-xs">✓</span>
                  )}
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Select all questions toggle */}
      <div className="flex items-center justify-between mb-3">
        <span className="text-xs text-gray-500">
          {selectedQuestions.size} of {questions.length} selected
        </span>
        <button
          onClick={selectAllQuestions}
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
                  : "bg-gray-900 border-gray-800"
              }`}
            >
              <div className="min-w-0">
                <p className="text-gray-100 font-medium text-sm truncate">{q.title}</p>
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
                    : "bg-gray-800 hover:bg-gray-700 text-gray-300 border border-gray-700"
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
          disabled={!canAssign || assigned}
          className={`px-6 py-3 rounded-xl text-sm font-bold shadow-lg transition-all ${
            assigned
              ? "bg-green-600 text-white scale-95"
              : canAssign
              ? "bg-purple-600 hover:bg-purple-500 text-white hover:scale-105 active:scale-95"
              : "bg-gray-700 text-gray-500 cursor-not-allowed"
          }`}
        >
          {assigned
            ? "Assigned ✓"
            : `Assign ${selectedQuestions.size > 0 ? `(${selectedQuestions.size})` : ""} to ${
                selectedMentees.length > 0 ? `${selectedMentees.length} mentee${selectedMentees.length > 1 ? "s" : ""}` : "mentees"
              }`}
        </button>
      </div>
    </div>
  );
}
