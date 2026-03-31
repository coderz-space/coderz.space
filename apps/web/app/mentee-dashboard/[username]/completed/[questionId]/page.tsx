"use client";

import { useParams, useRouter, useSearchParams } from "next/navigation";
import { useState, Suspense } from "react";
import { getQuestionDetail, updateQuestionDetails } from "@/services";
import type { Question } from "@/types";

function QuestionDetailContent() {
  const { username, questionId } = useParams() as { username: string; questionId: string };
  const searchParams = useSearchParams();
  const router = useRouter();

  // owner = whose notes we're viewing; defaults to the logged-in user (username)
  const owner = searchParams.get("owner") ?? username;
  const isReadOnly = owner !== username;

  const [question] = useState<Question | null>(() =>
    // Replace with: GET /api/mentees/:owner/questions/:questionId
    getQuestionDetail(owner, questionId)
  );
  const [solution, setSolution] = useState(() => question?.solution ?? "");
  const [resources, setResources] = useState(() => question?.resources ?? "");
  const [saved, setSaved] = useState(false);

  const handleSave = () => {
    // Replace with: PATCH /api/mentees/:username/questions/:questionId  { solution, resources }
    updateQuestionDetails(username, questionId, { solution, resources });
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  };

  if (!question) {
    return <p className="text-gray-500 text-sm">Question not found.</p>;
  }

  return (
    <div className="max-w-2xl">
      {/* Back */}
      <button
        onClick={() => router.back()}
        className="text-sm text-purple-400 hover:text-purple-300 mb-6 flex items-center gap-1"
      >
        ← Back
      </button>

      {/* Question header */}
      <div className="flex items-center gap-3 mb-1">
        <h1 className="text-2xl font-bold text-purple-400">{question.title}</h1>
        {isReadOnly && (
          <span className="text-xs font-semibold px-2.5 py-1 rounded-full bg-gray-700 text-gray-400">
            @{owner} · read-only
          </span>
        )}
      </div>
      <p className="text-sm text-gray-400 mb-6">{question.description}</p>

      {/* Solution */}
      <div className="mb-5">
        <label className="block text-sm font-semibold text-gray-300 mb-2">Solution</label>
        <textarea
          value={solution}
          onChange={(e) => !isReadOnly && setSolution(e.target.value)}
          readOnly={isReadOnly}
          rows={8}
          placeholder={isReadOnly ? "No solution written yet." : "Write your solution approach, code, or notes here..."}
          className={`w-full border rounded-xl p-4 text-sm text-gray-100 placeholder-gray-600 focus:outline-none resize-y ${
            isReadOnly
              ? "bg-gray-900 border-gray-800 cursor-default text-gray-400"
              : "bg-gray-800 border-gray-700 focus:ring-2 focus:ring-purple-500"
          }`}
        />
      </div>

      {/* Resources */}
      <div className="mb-6">
        <label className="block text-sm font-semibold text-gray-300 mb-2">Resources</label>
        <textarea
          value={resources}
          onChange={(e) => !isReadOnly && setResources(e.target.value)}
          readOnly={isReadOnly}
          rows={4}
          placeholder={isReadOnly ? "No resources added yet." : "Paste links, references, or notes on resources you used..."}
          className={`w-full border rounded-xl p-4 text-sm text-gray-100 placeholder-gray-600 focus:outline-none resize-y ${
            isReadOnly
              ? "bg-gray-900 border-gray-800 cursor-default text-gray-400"
              : "bg-gray-800 border-gray-700 focus:ring-2 focus:ring-purple-500"
          }`}
        />
      </div>

      {/* Save — only shown to the owner */}
      {!isReadOnly && (
        <button
          onClick={handleSave}
          className="px-5 py-2 bg-purple-700 hover:bg-purple-600 text-white text-sm font-semibold rounded-lg transition-colors"
        >
          {saved ? "Saved ✓" : "Save"}
        </button>
      )}
    </div>
  );
}

export default function QuestionDetailPage() {
  return (
    <Suspense fallback={<p className="text-gray-500 text-sm">Loading...</p>}>
      <QuestionDetailContent />
    </Suspense>
  );
}
