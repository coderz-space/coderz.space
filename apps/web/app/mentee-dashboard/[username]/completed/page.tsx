"use client";

import { useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { getMenteeQuestions, updateQuestionProgress } from "@/services";
import type { Question, QuestionProgressStatus } from "@/types";

const difficultyColor: Record<Question["difficulty"], string> = {
  easy: "text-green-400 bg-green-900/30",
  medium: "text-yellow-400 bg-yellow-900/30",
  hard: "text-red-400 bg-red-900/30",
};

const progressOptions: { value: QuestionProgressStatus; label: string }[] = [
  { value: "not_started",       label: "Not Started" },
  { value: "discussion_needed", label: "Discussion Needed" },
  { value: "revision_needed",   label: "Revision Needed" },
  { value: "completed",         label: "Completed" },
];

const progressColor: Record<QuestionProgressStatus, string> = {
  not_started:       "text-gray-400 bg-gray-700",
  discussion_needed: "text-yellow-400 bg-yellow-900/40",
  revision_needed:   "text-orange-400 bg-orange-900/40",
  completed:         "text-green-400 bg-green-900/40",
};

export default function CompletedQuestionsPage() {
  const { username } = useParams() as { username: string };
  const router = useRouter();
  const [questions, setQuestions] = useState<Question[]>([]);

  useEffect(() => {
    const loadQuestions = async () => {
      const allQuestions = await getMenteeQuestions(username);
      setQuestions(allQuestions.filter((q) => q.status === "completed"));
    };
    loadQuestions();
  }, [username]);

  const handleProgressChange = async (questionId: string, value: QuestionProgressStatus) => {
    // Replace with: PATCH /api/mentees/:username/questions/:questionId  { progressStatus }
    await updateQuestionProgress(username, questionId, value);
    // Re-filter: if moved away from completed it drops off this list
    const allQuestions = await getMenteeQuestions(username);
    setQuestions(allQuestions.filter((q) => q.status === "completed"));
  };

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-purple-400 mb-6">Completed Questions</h1>

      {questions.length === 0 && (
        <p className="text-gray-500 text-sm">No completed questions yet. Keep going!</p>
      )}

      <div className="flex flex-col gap-4">
        {questions.map((q) => (
          <div key={q.id} className="bg-gray-800/60 border border-gray-700 rounded-xl p-5">
            <div className="flex items-start justify-between gap-3 mb-2">
              <h2 className="font-semibold text-gray-200">{q.title}</h2>
              <span className={`text-xs font-semibold px-2.5 py-1 rounded-full shrink-0 ${difficultyColor[q.difficulty]}`}>
                {q.difficulty}
              </span>
            </div>

            <p className="text-sm text-gray-500 mb-3">{q.description}</p>

            <div className="flex items-center justify-between gap-3 flex-wrap">
              <div className="flex items-center gap-3 text-xs text-gray-600">
                <span className="bg-gray-700 px-2 py-0.5 rounded">{q.topic}</span>
              </div>
              <div className="flex gap-4 text-xs text-gray-500 mt-1">
                <span>Assigned on {new Date(q.assignedAt).toLocaleDateString()}</span>
                {q.completedAt && (
                  <span>Solved on {new Date(q.completedAt).toLocaleDateString()}</span>
                )}
              </div>

              {/* Progress dropdown + Detail button */}
              <div className="flex items-center gap-2 flex-wrap">
                <button
                  onClick={() => router.push(`/mentee-dashboard/${username}/completed/${q.id}`)}
                  className="text-xs font-semibold px-3 py-1.5 rounded-lg bg-purple-700 hover:bg-purple-600 text-white transition-colors"
                >
                  Detail
                </button>
                <select
                  value={q.progressStatus}
                  onChange={(e) => handleProgressChange(q.id, e.target.value as QuestionProgressStatus)}
                  className={`text-xs font-semibold px-3 py-1.5 rounded-lg border-0 cursor-pointer focus:outline-none focus:ring-2 focus:ring-purple-500 ${progressColor[q.progressStatus]}`}
                >
                  {progressOptions.map((opt) => (
                    <option key={opt.value} value={opt.value}
                      className="bg-gray-800 text-gray-100">
                      {opt.label}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
