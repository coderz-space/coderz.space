import type { Question, QuestionProgressStatus } from "@/types";

export const difficultyBadgeClass: Record<Question["difficulty"], string> = {
  easy: "bg-emerald-100 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-300",
  medium: "bg-amber-100 text-amber-700 dark:bg-amber-950/60 dark:text-amber-300",
  hard: "bg-rose-100 text-rose-700 dark:bg-rose-950/60 dark:text-rose-300",
};

export const progressBadgeClass: Record<QuestionProgressStatus, string> = {
  not_started: "bg-slate-200 text-slate-700 dark:bg-slate-800 dark:text-slate-300",
  discussion_needed: "bg-yellow-100 text-yellow-700 dark:bg-yellow-950/60 dark:text-yellow-300",
  revision_needed: "bg-orange-100 text-orange-700 dark:bg-orange-950/60 dark:text-orange-300",
  completed: "bg-emerald-100 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-300",
};

export const progressLabel: Record<QuestionProgressStatus, string> = {
  not_started: "Not Started",
  discussion_needed: "Discussion Needed",
  revision_needed: "Revision Needed",
  completed: "Completed",
};

export const progressOptions: { value: QuestionProgressStatus; label: string }[] = [
  { value: "not_started", label: "Not Started" },
  { value: "discussion_needed", label: "Discussion Needed" },
  { value: "revision_needed", label: "Revision Needed" },
  { value: "completed", label: "Completed" },
];

export function formatDate(value?: string): string {
  if (!value) {
    return "Unknown";
  }

  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) {
    return "Unknown";
  }

  return parsed.toLocaleDateString();
}

export function getInitials(firstName?: string, lastName?: string, fallback = "AB"): string {
  const first = firstName?.trim().charAt(0) ?? "";
  const last = lastName?.trim().charAt(0) ?? "";
  const initials = `${first}${last}`.trim();
  return initials || fallback;
}
