"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import {
  difficultyBadgeClass,
  formatDate,
  progressBadgeClass,
  progressOptions,
} from "@/components/dashboard/constants";
import { getMenteeQuestions, updateQuestionProgress } from "@/services/mentee";
import type { Question, QuestionProgressStatus } from "@/types";

export default function PendingQuestionsPage() {
  const { username } = useParams() as { username: string };
  const [questions, setQuestions] = useState<Question[]>([]);

  useEffect(() => {
    const loadQuestions = async () => {
      const allQuestions = await getMenteeQuestions(username);
      setQuestions(allQuestions.filter((question) => question.status === "pending"));
    };

    void loadQuestions();
  }, [username]);

  const handleProgressChange = async (questionId: string, progressStatus: QuestionProgressStatus) => {
    await updateQuestionProgress(username, questionId, progressStatus);
    const refreshedQuestions = await getMenteeQuestions(username);
    setQuestions(refreshedQuestions.filter((question) => question.status === "pending"));
  };

  return (
    <div className="max-w-3xl">
      <h1 className="mb-6 text-2xl font-bold text-gray-900 dark:text-white">Pending Questions</h1>

      {!questions.length ? (
        <p className="text-sm text-gray-500 dark:text-gray-400">No pending questions. You are all caught up.</p>
      ) : null}

      <div className="flex flex-col gap-4">
        {questions.map((question) => (
          <article
            key={question.id}
            className="rounded-2xl border border-gray-200 bg-white p-5 shadow-sm dark:border-gray-800 dark:bg-gray-900"
          >
            <div className="mb-3 flex items-start justify-between gap-3">
              <div className="min-w-0">
                <h2 className="font-semibold text-gray-900 dark:text-gray-100">{question.title}</h2>
                <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{question.description}</p>
              </div>
              <span
                className={`shrink-0 rounded-full px-2.5 py-1 text-xs font-semibold ${difficultyBadgeClass[question.difficulty]}`}
              >
                {question.difficulty}
              </span>
            </div>

            <div className="mb-4 flex flex-wrap items-center gap-3 text-xs text-gray-500 dark:text-gray-400">
              <span className="rounded-full bg-gray-100 px-2.5 py-1 dark:bg-gray-800">{question.topic}</span>
              <span>Assigned {formatDate(question.assignedAt)}</span>
            </div>

            <select
              value={question.progressStatus}
              onChange={(event) =>
                handleProgressChange(question.id, event.target.value as QuestionProgressStatus)
              }
              className={`rounded-lg border-0 px-3 py-2 text-xs font-semibold outline-none ring-1 ring-inset ring-transparent transition focus:ring-purple-400 ${progressBadgeClass[question.progressStatus]}`}
            >
              {progressOptions.map((option) => (
                <option
                  key={option.value}
                  value={option.value}
                  className="bg-white text-gray-900 dark:bg-gray-900 dark:text-gray-100"
                >
                  {option.label}
                </option>
              ))}
            </select>
          </article>
        ))}
      </div>
    </div>
  );
}
