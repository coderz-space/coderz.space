"use client";

import type { Profile, Question } from "@/types";
import {
  difficultyBadgeClass,
  formatDate,
  getInitials,
  progressBadgeClass,
  progressLabel,
} from "./constants";

type PublicMenteeProfileProps = {
  viewerUsername?: string;
  profile: Profile;
  completedQuestions: Question[];
  onBack: () => void;
  onOpenQuestion?: (questionId: string) => void;
};

export default function PublicMenteeProfile({
  viewerUsername,
  profile,
  completedQuestions,
  onBack,
  onOpenQuestion,
}: PublicMenteeProfileProps) {
  const isOwnProfile = viewerUsername === profile.username;

  return (
    <div className="max-w-3xl">
      <button
        type="button"
        onClick={onBack}
        className="mb-6 text-sm font-medium text-purple-600 transition-colors hover:text-purple-500 dark:text-purple-300 dark:hover:text-purple-200"
      >
        Back
      </button>

      <section className="mb-8 rounded-3xl border border-purple-200 bg-white p-6 shadow-sm dark:border-purple-900/60 dark:bg-gray-900">
        <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-purple-600 text-2xl font-bold text-white">
          {getInitials(profile.firstName, profile.lastName, "MB")}
        </div>

        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          {profile.firstName} {profile.lastName}
          {isOwnProfile ? (
            <span className="ml-2 text-sm font-normal text-purple-600 dark:text-purple-300">(you)</span>
          ) : null}
        </h1>
        <p className="mb-4 text-sm text-gray-500 dark:text-gray-400">@{profile.username}</p>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="rounded-2xl border border-gray-200 bg-gray-50 px-4 py-3 dark:border-gray-800 dark:bg-gray-950/60">
            <p className="text-xl font-bold text-purple-600 dark:text-purple-300">{profile.solved}</p>
            <p className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">Questions Solved</p>
          </div>
          <div className="rounded-2xl border border-gray-200 bg-gray-50 px-4 py-3 dark:border-gray-800 dark:bg-gray-950/60">
            <p className="text-sm font-semibold text-gray-900 dark:text-gray-100">{formatDate(profile.joinedAt)}</p>
            <p className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">Joined</p>
          </div>
        </div>

        {profile.bio || profile.github || profile.linkedin ? (
          <div className="mt-6 border-t border-gray-200 pt-5 dark:border-gray-800">
            {profile.bio ? (
              <p className="mb-4 text-sm leading-6 text-gray-700 dark:text-gray-300">{profile.bio}</p>
            ) : null}

            <div className="flex flex-wrap gap-3">
              {profile.github ? (
                <a
                  href={profile.github}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm font-medium text-purple-600 hover:underline dark:text-purple-300"
                >
                  GitHub
                </a>
              ) : null}
              {profile.linkedin ? (
                <a
                  href={profile.linkedin}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm font-medium text-purple-600 hover:underline dark:text-purple-300"
                >
                  LinkedIn
                </a>
              ) : null}
            </div>
          </div>
        ) : null}
      </section>

      <section>
        <h2 className="mb-4 text-lg font-semibold text-gray-900 dark:text-white">Solved Questions</h2>

        {!completedQuestions.length ? (
          <p className="text-sm text-gray-500 dark:text-gray-400">No solved questions yet.</p>
        ) : (
          <div className="flex flex-col gap-3">
            {completedQuestions.map((question) => (
              <article
                key={question.id}
                className="rounded-2xl border border-gray-200 bg-white px-4 py-4 dark:border-gray-800 dark:bg-gray-900"
              >
                <div className="mb-2 flex items-start justify-between gap-3">
                  <div className="min-w-0">
                    <p className="truncate font-semibold text-gray-900 dark:text-gray-100">{question.title}</p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">{question.topic}</p>
                  </div>
                  <span
                    className={`shrink-0 rounded-full px-2.5 py-1 text-xs font-semibold ${difficultyBadgeClass[question.difficulty]}`}
                  >
                    {question.difficulty}
                  </span>
                </div>

                <div className="mb-3 flex flex-wrap gap-3 text-xs text-gray-500 dark:text-gray-400">
                  <span>Assigned {formatDate(question.assignedAt)}</span>
                  {question.completedAt ? <span>Solved {formatDate(question.completedAt)}</span> : null}
                </div>

                <div className="flex items-center justify-between gap-3">
                  <span
                    className={`rounded-lg px-2.5 py-1 text-xs font-semibold ${progressBadgeClass[question.progressStatus]}`}
                  >
                    {progressLabel[question.progressStatus]}
                  </span>

                  {onOpenQuestion ? (
                    <button
                      type="button"
                      onClick={() => onOpenQuestion(question.id)}
                      className="rounded-lg bg-purple-600 px-3 py-1.5 text-xs font-semibold text-white transition-colors hover:bg-purple-500"
                    >
                      Detail
                    </button>
                  ) : null}
                </div>
              </article>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
