"use client";

import { Suspense, useEffect, useState } from "react";
import { useParams, useRouter, useSearchParams } from "next/navigation";
import { getQuestionDetail, updateQuestionDetails } from "@/services/mentee";
import type { Question } from "@/types";

function QuestionDetailContent() {
  const { username, questionId } = useParams() as { username: string; questionId: string };
  const searchParams = useSearchParams();
  const router = useRouter();
  const owner = searchParams.get("owner") ?? username;
  const isReadOnly = owner !== username;

  const [question, setQuestion] = useState<Question | null>(null);
  const [solution, setSolution] = useState("");
  const [resources, setResources] = useState("");
  const [saveMessage, setSaveMessage] = useState("");

  useEffect(() => {
    const loadQuestion = async () => {
      const nextQuestion = await getQuestionDetail(owner, questionId);
      setQuestion(nextQuestion);
      setSolution(nextQuestion?.solution ?? "");
      setResources(nextQuestion?.resources ?? "");
    };

    void loadQuestion();
  }, [owner, questionId]);

  const handleSave = async () => {
    await updateQuestionDetails(username, questionId, { solution, resources });
    setSaveMessage("Saved.");
  };

  if (!question) {
    return <p className="text-sm text-gray-500 dark:text-gray-400">Question not found.</p>;
  }

  return (
    <div className="max-w-3xl">
      <button
        type="button"
        onClick={() => router.back()}
        className="mb-6 text-sm font-medium text-purple-600 transition-colors hover:text-purple-500 dark:text-purple-300 dark:hover:text-purple-200"
      >
        Back
      </button>

      <div className="mb-2 flex items-center gap-3">
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">{question.title}</h1>
        {isReadOnly ? (
          <span className="rounded-full bg-gray-100 px-2.5 py-1 text-xs font-semibold text-gray-600 dark:bg-gray-800 dark:text-gray-300">
            @{owner} read-only
          </span>
        ) : null}
      </div>
      <p className="mb-6 text-sm text-gray-500 dark:text-gray-400">{question.description}</p>

      <div className="mb-5">
        <label className="mb-2 block text-sm font-semibold text-gray-700 dark:text-gray-300">Solution</label>
        <textarea
          value={solution}
          onChange={(event) => {
            if (!isReadOnly) {
              setSolution(event.target.value);
            }
          }}
          readOnly={isReadOnly}
          rows={8}
          placeholder={isReadOnly ? "No solution written yet." : "Write your solution, approach, or notes here."}
          className={`w-full rounded-2xl border p-4 text-sm outline-none ${
            isReadOnly
              ? "cursor-default border-gray-200 bg-gray-50 text-gray-500 dark:border-gray-800 dark:bg-gray-950 dark:text-gray-400"
              : "border-gray-200 bg-white text-gray-900 focus:border-purple-500 focus:ring-2 focus:ring-purple-200 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-100 dark:focus:border-purple-500 dark:focus:ring-purple-900/60"
          }`}
        />
      </div>

      <div className="mb-6">
        <label className="mb-2 block text-sm font-semibold text-gray-700 dark:text-gray-300">Resources</label>
        <textarea
          value={resources}
          onChange={(event) => {
            if (!isReadOnly) {
              setResources(event.target.value);
            }
          }}
          readOnly={isReadOnly}
          rows={4}
          placeholder={isReadOnly ? "No resources added yet." : "Paste links, references, or notes here."}
          className={`w-full rounded-2xl border p-4 text-sm outline-none ${
            isReadOnly
              ? "cursor-default border-gray-200 bg-gray-50 text-gray-500 dark:border-gray-800 dark:bg-gray-950 dark:text-gray-400"
              : "border-gray-200 bg-white text-gray-900 focus:border-purple-500 focus:ring-2 focus:ring-purple-200 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-100 dark:focus:border-purple-500 dark:focus:ring-purple-900/60"
          }`}
        />
      </div>

      {!isReadOnly ? (
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={handleSave}
            className="rounded-lg bg-purple-600 px-5 py-2 text-sm font-semibold text-white transition-colors hover:bg-purple-500"
          >
            Save
          </button>
          {saveMessage ? <p className="text-sm text-gray-600 dark:text-gray-300">{saveMessage}</p> : null}
        </div>
      ) : null}
    </div>
  );
}

export default function QuestionDetailPage() {
  return (
    <Suspense fallback={<p className="text-sm text-gray-500 dark:text-gray-400">Loading...</p>}>
      <QuestionDetailContent />
    </Suspense>
  );
}
