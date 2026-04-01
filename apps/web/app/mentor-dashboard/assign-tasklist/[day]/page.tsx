"use client";

import { use, useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import Modal from "@/components/Modal";
import { getMenteeQuestions } from "@/services/mentee";
import { getDayAssignments, updateDayAssignments } from "@/services/mentor";
import type { DayAssignmentMentee, Question } from "@/types";

const difficultyColor: Record<Question["difficulty"], string> = {
  easy: "text-green-400",
  medium: "text-yellow-400",
  hard: "text-red-400",
};

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

export default function DayPage({ params }: { params: Promise<{ day: string }> }) {
  const { day } = use(params);
  const router = useRouter();
  const dayLabel = capitalize(day);
  const [mentees, setMentees] = useState<DayAssignmentMentee[]>([]);
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [taskModal, setTaskModal] = useState<{ mentee: DayAssignmentMentee; tasks: Question[] } | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let active = true;

    void getDayAssignments(day).then((response) => {
      if (active) {
        setMentees(response.mentees);
      }
    });

    return () => {
      active = false;
    };
  }, [day]);

  useEffect(() => {
    const handleClick = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setDropdownOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  const assignedMentees = useMemo(() => mentees.filter((mentee) => mentee.assigned), [mentees]);
  const unassignedMentees = useMemo(() => mentees.filter((mentee) => !mentee.assigned), [mentees]);

  const persistAssignments = async (usernames: string[]) => {
    const response = await updateDayAssignments(day, usernames);
    setMentees(response.mentees);
  };

  const addMentee = async (username: string) => {
    const next = [...assignedMentees.map((mentee) => mentee.username), username];
    await persistAssignments(next);
    setDropdownOpen(false);
  };

  const removeMentee = async (username: string) => {
    const next = assignedMentees.map((mentee) => mentee.username).filter((value) => value !== username);
    await persistAssignments(next);
  };

  const openTaskModal = async (mentee: DayAssignmentMentee) => {
    const tasks = await getMenteeQuestions(mentee.username);
    setTaskModal({ mentee, tasks });
  };

  return (
    <div className="max-w-2xl">
      <button
        type="button"
        onClick={() => router.back()}
        className="mb-4 flex items-center gap-1 text-xs text-gray-500 transition-colors hover:text-gray-700 dark:hover:text-gray-300"
      >
        Back
      </button>

      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-purple-400">{dayLabel}</h1>
          <p className="mt-0.5 text-sm text-gray-500">
            {assignedMentees.length} mentee{assignedMentees.length !== 1 ? "s" : ""} assigned
          </p>
        </div>

        <div className="relative" ref={dropdownRef}>
          <button
            type="button"
            onClick={() => setDropdownOpen((value) => !value)}
            className="flex items-center gap-2 rounded-lg bg-purple-700 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-purple-600"
          >
            + Add Mentee
            <span className="text-xs text-purple-300">{dropdownOpen ? "^" : "v"}</span>
          </button>

          {dropdownOpen ? (
            <div className="absolute right-0 top-full z-20 mt-2 w-64 overflow-hidden rounded-xl border border-gray-200 bg-white shadow-xl dark:border-gray-700 dark:bg-gray-900">
              {unassignedMentees.length === 0 ? (
                <p className="px-4 py-3 text-xs text-gray-500">All mentees already added</p>
              ) : (
                unassignedMentees.map((mentee) => (
                  <button
                    key={mentee.username}
                    type="button"
                    onClick={() => void addMentee(mentee.username)}
                    className="flex w-full items-center gap-3 px-4 py-2.5 text-left text-sm text-gray-700 transition-colors hover:bg-purple-50 dark:text-gray-200 dark:hover:bg-gray-800"
                  >
                    <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-purple-700 text-xs font-bold text-white">
                      {mentee.firstName[0]}
                      {mentee.lastName[0] ?? ""}
                    </div>
                    <span>
                      {mentee.firstName} {mentee.lastName}
                    </span>
                  </button>
                ))
              )}
            </div>
          ) : null}
        </div>
      </div>

      {assignedMentees.length === 0 ? (
        <div className="rounded-xl border border-dashed border-gray-200 py-16 text-center text-sm text-gray-400 dark:border-gray-800 dark:text-gray-600">
          No mentees assigned to {dayLabel} yet. Use Add Mentee to get started.
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {assignedMentees.map((mentee) => (
            <div
              key={mentee.username}
              className="flex items-center justify-between gap-4 rounded-xl border border-gray-200 bg-white px-5 py-4 dark:border-gray-800 dark:bg-gray-900"
            >
              <div className="flex min-w-0 items-center gap-3">
                <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-purple-700 text-sm font-bold text-white">
                  {mentee.firstName[0]}
                  {mentee.lastName[0] ?? ""}
                </div>
                <div className="min-w-0">
                  <p className="truncate text-sm font-semibold text-gray-800 dark:text-gray-100">
                    {mentee.firstName} {mentee.lastName}
                  </p>
                  <p className="text-xs text-gray-500">
                    @{mentee.username}
                    {mentee.assignedSheet ? (
                      <span className="ml-2 text-purple-400">| {mentee.assignedSheet}</span>
                    ) : null}
                  </p>
                </div>
              </div>

              <div className="flex shrink-0 items-center gap-2">
                <button
                  type="button"
                  onClick={() => void openTaskModal(mentee)}
                  className="rounded-lg border border-gray-200 bg-gray-100 px-3 py-1.5 text-xs font-semibold text-gray-700 transition-colors hover:bg-gray-200 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
                >
                  Assigned Tasks
                </button>
                <button
                  type="button"
                  onClick={() =>
                    router.push(
                      `/mentor-dashboard/assign-tasklist/${day}/${mentee.username}/${mentee.assignedSheet ?? "gfg-dsa-360"}`
                    )
                  }
                  className="rounded-lg bg-purple-700 px-3 py-1.5 text-xs font-semibold text-white transition-colors hover:bg-purple-600"
                >
                  Assign Questions
                </button>
                <button
                  type="button"
                  onClick={() => void removeMentee(mentee.username)}
                  className="rounded-lg px-2 py-1.5 text-xs text-red-400 transition-colors hover:bg-red-50 hover:text-red-500 dark:hover:bg-red-900/20"
                  aria-label="Remove mentee"
                >
                  Remove
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {taskModal ? (
        <Modal
          onClose={() => setTaskModal(null)}
          className="flex max-h-[80vh] w-full max-w-lg flex-col rounded-2xl border border-gray-200 bg-white shadow-2xl dark:border-gray-700 dark:bg-gray-900"
        >
          <div className="flex items-center justify-between border-b border-gray-100 px-6 py-4 dark:border-gray-800">
            <div>
              <h2 className="text-base font-bold text-purple-400">Assigned Tasks</h2>
              <p className="mt-0.5 text-xs text-gray-500">
                {taskModal.mentee.firstName} {taskModal.mentee.lastName} | {taskModal.tasks.length} task
                {taskModal.tasks.length !== 1 ? "s" : ""}
              </p>
            </div>
            <button
              type="button"
              onClick={() => setTaskModal(null)}
              className="text-sm font-semibold text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-gray-200"
            >
              Close
            </button>
          </div>

          <div className="flex-1 overflow-y-auto px-6 py-4">
            {taskModal.tasks.length === 0 ? (
              <p className="py-8 text-center text-sm text-gray-500">No tasks assigned yet.</p>
            ) : (
              <div className="flex flex-col gap-3">
                {taskModal.tasks.map((task, index) => (
                  <div
                    key={task.id}
                    className="flex items-start justify-between gap-3 rounded-xl border border-gray-100 bg-gray-50 px-4 py-3 dark:border-gray-800 dark:bg-gray-800/50"
                  >
                    <div className="min-w-0">
                      <p className="truncate text-sm font-medium text-gray-800 dark:text-gray-100">
                        {index + 1}. {task.title}
                      </p>
                      <p className="mt-0.5 text-xs text-gray-500">
                        {task.topic} |{" "}
                        <span className={difficultyColor[task.difficulty]}>{task.difficulty}</span>
                      </p>
                      <div className="mt-1.5 flex flex-col gap-0.5">
                        <span className="text-xs text-gray-400">
                          Assigned on {new Date(task.assignedAt).toLocaleDateString()}
                        </span>
                        {task.completedAt ? (
                          <span className="text-xs text-green-400">
                            Solved on {new Date(task.completedAt).toLocaleDateString()}
                          </span>
                        ) : null}
                      </div>
                    </div>
                    <span
                      className={`shrink-0 rounded-full px-2 py-0.5 text-xs font-semibold ${
                        task.status === "completed"
                          ? "bg-green-100 text-green-600 dark:bg-green-900/40 dark:text-green-400"
                          : "bg-yellow-100 text-yellow-600 dark:bg-yellow-900/30 dark:text-yellow-400"
                      }`}
                    >
                      {task.status}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>
        </Modal>
      ) : null}
    </div>
  );
}
