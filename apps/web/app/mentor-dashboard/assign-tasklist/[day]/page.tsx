"use client";

import { use, useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import { getMenteeRequests, getAssignedTasks } from "@/services";
import Modal from "@/components/Modal";
import type { MenteeRequest, Question } from "@/types";

const STORAGE_KEY = (day: string) => `coderz_day_mentees_${day}`;

const DIFFICULTY_COLOR: Record<string, string> = {
  easy: "text-green-400",
  medium: "text-yellow-400",
  hard: "text-red-400",
};

function capitalize(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

export default function DayPage({ params }: { params: Promise<{ day: string }> }) {
  const { day } = use(params);
  const router = useRouter();
  const dayLabel = capitalize(day);

  const [allApproved, setAllApproved] = useState<MenteeRequest[]>([]);

  const [assignedIds, setAssignedIds] = useState<string[]>(() => {
    if (typeof window === "undefined") return [];
    return JSON.parse(localStorage.getItem(STORAGE_KEY(day)) || "[]");
  });

  const [dropdownOpen, setDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Modal state: which mentee's tasks to show
  const [taskModal, setTaskModal] = useState<{ mentee: MenteeRequest; tasks: Question[] } | null>(null);

  // Load approved mentees
  useEffect(() => {
    const loadMentees = async () => {
      const requests = await getMenteeRequests();
      const approved = requests.filter((r) => r.status === "approved");
      setAllApproved(approved);
    };
    loadMentees();
  }, []);

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setDropdownOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  const save = (ids: string[]) => {
    setAssignedIds(ids);
    localStorage.setItem(STORAGE_KEY(day), JSON.stringify(ids));
  };

  const addMentee = (id: string) => {
    if (!assignedIds.includes(id)) save([...assignedIds, id]);
    setDropdownOpen(false);
  };

  const removeMentee = (id: string) => save(assignedIds.filter((i) => i !== id));

  const openTaskModal = (m: MenteeRequest) => {
    const tasks = getAssignedTasks(m.username);
    setTaskModal({ mentee: m, tasks });
  };

  const assignedMentees = allApproved.filter((m) => assignedIds.includes(m.id));
  const unassigned = allApproved.filter((m) => !assignedIds.includes(m.id));

  return (
    <div className="max-w-2xl">
      <button
        onClick={() => router.back()}
        className="text-xs text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 mb-4 flex items-center gap-1 transition-colors"
      >
        ← Back
      </button>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-purple-400">{dayLabel}</h1>
          <p className="text-gray-500 text-sm mt-0.5">
            {assignedMentees.length} mentee{assignedMentees.length !== 1 ? "s" : ""} assigned
          </p>
        </div>

        {/* Add mentee dropdown */}
        <div className="relative" ref={dropdownRef}>
          <button
            onClick={() => setDropdownOpen((v) => !v)}
            className="flex items-center gap-2 text-sm font-semibold px-4 py-2 rounded-lg bg-purple-700 hover:bg-purple-600 text-white transition-colors"
          >
            + Add Mentee
            <span className="text-purple-300 text-xs">{dropdownOpen ? "▲" : "▼"}</span>
          </button>

          {dropdownOpen && (
            <div className="absolute right-0 top-full mt-2 w-64 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-xl shadow-xl z-20 overflow-hidden">
              {unassigned.length === 0 ? (
                <p className="px-4 py-3 text-xs text-gray-500">All mentees already added</p>
              ) : (
                unassigned.map((m) => (
                  <button
                    key={m.id}
                    onClick={() => addMentee(m.id)}
                    className="w-full text-left px-4 py-2.5 text-sm text-gray-700 dark:text-gray-200 hover:bg-purple-50 dark:hover:bg-gray-800 transition-colors flex items-center gap-3"
                  >
                    <div className="w-7 h-7 rounded-full bg-purple-700 flex items-center justify-center text-xs font-bold text-white shrink-0">
                      {m.firstName[0]}{m.lastName[0]}
                    </div>
                    <span>{m.firstName} {m.lastName}</span>
                  </button>
                ))
              )}
            </div>
          )}
        </div>
      </div>

      {/* Mentee list */}
      {assignedMentees.length === 0 ? (
        <div className="text-center py-16 text-gray-400 dark:text-gray-600 text-sm border border-dashed border-gray-200 dark:border-gray-800 rounded-xl">
          No mentees assigned to {dayLabel} yet. Use "Add Mentee" to get started.
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {assignedMentees.map((m) => {
            const taskCount = getAssignedTasks(m.username).length;
            return (
              <div
                key={m.id}
                className="flex items-center justify-between gap-4 px-5 py-4 rounded-xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900"
              >
                <div className="flex items-center gap-3 min-w-0">
                  <div className="w-9 h-9 rounded-full bg-purple-700 flex items-center justify-center text-sm font-bold text-white shrink-0">
                    {m.firstName[0]}{m.lastName[0]}
                  </div>
                  <div className="min-w-0">
                    <p className="font-semibold text-gray-800 dark:text-gray-100 text-sm truncate">
                      {m.firstName} {m.lastName}
                    </p>
                    <p className="text-xs text-gray-500">
                      @{m.username}
                      {m.assignedSheet && (
                        <span className="ml-2 text-purple-400">
                          · {m.assignedSheet === "gfg-dsa-360" ? "GFG DSA 360" : "Striver's DSA Sheet"}
                        </span>
                      )}
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-2 shrink-0">
                  <button
                    onClick={() => openTaskModal(m)}
                    className="text-xs font-semibold px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300 border border-gray-200 dark:border-gray-700 transition-colors"
                  >
                    Assigned Tasks
                    {taskCount > 0 && (
                      <span className="ml-1.5 bg-purple-600 text-white text-xs rounded-full px-1.5 py-0.5">
                        {taskCount}
                      </span>
                    )}
                  </button>
                  <button
                    onClick={() => {
                      const sheet = m.assignedSheet ?? "gfg-dsa-360";
                      router.push(`/mentor-dashboard/assign-tasklist/${day}/${m.username}/${sheet}`);
                    }}
                    className="text-xs font-semibold px-3 py-1.5 rounded-lg bg-purple-700 hover:bg-purple-600 text-white transition-colors"
                  >
                    Assign Questions
                  </button>
                  <button
                    onClick={() => removeMentee(m.id)}
                    className="text-xs px-2 py-1.5 rounded-lg text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-500 transition-colors"
                    aria-label="Remove mentee"
                  >
                    ✕
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Assigned Tasks Modal */}
      {taskModal && (
        <Modal onClose={() => setTaskModal(null)} className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-2xl w-full max-w-lg max-h-[80vh] flex flex-col shadow-2xl">
            {/* Header */}
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100 dark:border-gray-800">
              <div>
                <h2 className="text-base font-bold text-purple-400">Assigned Tasks</h2>
                <p className="text-xs text-gray-500 mt-0.5">
                  {taskModal.mentee.firstName} {taskModal.mentee.lastName} · {taskModal.tasks.length} task{taskModal.tasks.length !== 1 ? "s" : ""}
                </p>
              </div>
              <button
                onClick={() => setTaskModal(null)}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 text-lg leading-none transition-colors"
              >
                ✕
              </button>
            </div>

            {/* Task list */}
            <div className="overflow-y-auto flex-1 px-6 py-4">
              {taskModal.tasks.length === 0 ? (
                <p className="text-sm text-gray-500 text-center py-8">No tasks assigned yet.</p>
              ) : (
                <div className="flex flex-col gap-3">
                  {taskModal.tasks.map((t, i) => (
                    <div
                      key={t.id}
                      className="flex items-start justify-between gap-3 px-4 py-3 rounded-xl border border-gray-100 dark:border-gray-800 bg-gray-50 dark:bg-gray-800/50"
                    >
                      <div className="min-w-0">
                        <p className="text-sm font-medium text-gray-800 dark:text-gray-100 truncate">
                          {i + 1}. {t.title}
                        </p>
                        <p className="text-xs text-gray-500 mt-0.5">
                          {t.topic} ·{" "}
                          <span className={DIFFICULTY_COLOR[t.difficulty]}>{t.difficulty}</span>
                        </p>
                        <div className="flex flex-col gap-0.5 mt-1.5">
                          <span className="text-xs text-gray-400">
                            Assigned on {new Date(t.assignedAt).toLocaleDateString()}
                          </span>
                          {t.completedAt && (
                            <span className="text-xs text-green-400">
                              Solved on {new Date(t.completedAt).toLocaleDateString()}
                            </span>
                          )}
                        </div>
                      </div>
                      <span className={`text-xs font-semibold px-2 py-0.5 rounded-full shrink-0 ${
                        t.status === "completed"
                          ? "bg-green-100 dark:bg-green-900/40 text-green-600 dark:text-green-400"
                          : "bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600 dark:text-yellow-400"
                      }`}>
                        {t.status}
                      </span>
                    </div>
                  ))}
                </div>
              )}
            </div>
        </Modal>
      )}
    </div>
  );
}
