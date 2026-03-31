"use client";

import { useParams, useRouter } from "next/navigation";
import { useState } from "react";
import { getMenteeProfile, getMenteeQuestions } from "@/services";
import type { Question, QuestionProgressStatus } from "@/types";

const difficultyColor: Record<Question["difficulty"], string> = {
  easy: "text-green-400 bg-green-900/30",
  medium: "text-yellow-400 bg-yellow-900/30",
  hard: "text-red-400 bg-red-900/30",
};

const progressColor: Record<QuestionProgressStatus, string> = {
  not_started:       "text-gray-400 bg-gray-700",
  discussion_needed: "text-yellow-400 bg-yellow-900/40",
  revision_needed:   "text-orange-400 bg-orange-900/40",
  completed:         "text-green-400 bg-green-900/40",
};

const progressLabel: Record<QuestionProgressStatus, string> = {
  not_started:       "Not Started",
  discussion_needed: "Discussion Needed",
  revision_needed:   "Revision Needed",
  completed:         "Completed",
};

export default function MentorMenteeProfilePage() {
  const { profileUsername } = useParams() as { profileUsername: string };
  const router = useRouter();
  // Replace with: GET /api/mentees/:profileUsername/profile
  const profile = getMenteeProfile(profileUsername);
  const [completedQuestions] = useState<Question[]>(() =>
    profile ? getMenteeQuestions(profileUsername).filter((q) => q.status === "completed") : []
  );

  if (!profile) {
    return <p className="text-gray-500 text-sm">Profile not found.</p>;
  }

  return (
    <div className="max-w-2xl">
      <button
        onClick={() => router.back()}
        className="text-sm text-purple-400 hover:text-purple-300 mb-6 flex items-center gap-1"
      >
        ← Back
      </button>

      {/* Profile card */}
      <div className="bg-gray-900 border border-purple-800/50 rounded-2xl p-6 mb-8">
        <div className="w-16 h-16 rounded-full bg-purple-700 flex items-center justify-center text-2xl font-bold text-white mb-4">
          {profile.firstName[0]}{profile.lastName[0]}
        </div>
        <h1 className="text-2xl font-bold text-white">{profile.firstName} {profile.lastName}</h1>
        <p className="text-gray-500 text-sm mb-4">@{profile.username}</p>
        <div className="flex gap-6 text-sm">
          <div>
            <p className="text-purple-400 font-bold text-xl">{profile.solved}</p>
            <p className="text-gray-500 text-xs">Questions Solved</p>
          </div>
          <div>
            <p className="text-gray-300 font-semibold">{new Date(profile.joinedAt).toLocaleDateString()}</p>
            <p className="text-gray-500 text-xs">Joined</p>
          </div>
        </div>

        {/* Bio / links */}
        {(profile.bio || profile.github || profile.linkedin) && (
          <div className="mt-5 pt-5 border-t border-gray-800 flex flex-col gap-3">
            {profile.bio && (
              <p className="text-sm text-gray-300 leading-relaxed">{profile.bio}</p>
            )}
            <div className="flex flex-wrap gap-3">
              {profile.github && (
                <a href={profile.github} target="_blank" rel="noopener noreferrer"
                  className="text-xs text-purple-400 hover:underline flex items-center gap-1">
                  GitHub ↗
                </a>
              )}
              {profile.linkedin && (
                <a href={profile.linkedin} target="_blank" rel="noopener noreferrer"
                  className="text-xs text-purple-400 hover:underline flex items-center gap-1">
                  LinkedIn ↗
                </a>
              )}
            </div>
          </div>
        )}
      </div>

      {/* Solved questions — read-only for mentor */}
      <h2 className="text-lg font-semibold text-purple-400 mb-4">Solved Questions</h2>

      {completedQuestions.length === 0 ? (
        <p className="text-gray-500 text-sm">No solved questions yet.</p>
      ) : (
        <div className="flex flex-col gap-3">
          {completedQuestions.map((q) => (
            <div key={q.id} className="bg-gray-900 border border-gray-800 rounded-xl px-4 py-3">
              <div className="flex items-start justify-between gap-3 mb-2">
                <div className="min-w-0">
                  <p className="text-sm font-medium text-gray-200 truncate">{q.title}</p>
                  <p className="text-xs text-gray-600">{q.topic}</p>
                  <div className="flex gap-3 text-xs text-gray-600 mt-1">
                    <span>Assigned on {new Date(q.assignedAt).toLocaleDateString()}</span>
                    {q.completedAt && (
                      <span>Solved on {new Date(q.completedAt).toLocaleDateString()}</span>
                    )}
                  </div>
                </div>
                <span className={`text-xs font-semibold px-2.5 py-1 rounded-full shrink-0 ${difficultyColor[q.difficulty]}`}>
                  {q.difficulty}
                </span>
              </div>
              <span className={`text-xs font-semibold px-2.5 py-1 rounded-lg ${progressColor[q.progressStatus]}`}>
                {progressLabel[q.progressStatus]}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
