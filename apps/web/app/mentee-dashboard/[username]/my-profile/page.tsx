"use client";

import { use, useState, useEffect } from "react";
import { getMenteeRequests, getMenteeQuestions, updateMenteeProfile, updateMenteePassword } from "@/services";
import type { MenteeRequest } from "@/types";

const SHEET_NAMES: Record<string, string> = {
  "gfg-dsa-360": "GFG DSA 360",
  "strivers-dsa-sheet": "Striver's DSA Sheet",
};

export default function MenteeMyProfilePage({
  params,
}: {
  params: Promise<{ username: string }>;
}) {
  const { username } = use(params);

  const [mentee, setMentee] = useState<MenteeRequest | null>(null);
  const [editing, setEditing] = useState(false);
  const [form, setForm] = useState({ bio: "", github: "", linkedin: "" });
  const [saved, setSaved] = useState(false);
  const [showPwCard, setShowPwCard] = useState(false);
  const [pwForm, setPwForm] = useState({ current: "", next: "", confirm: "" });
  const [pwError, setPwError] = useState("");
  const [pwSaved, setPwSaved] = useState(false);
  const [stats, setStats] = useState({ solved: 0, pending: 0, total: 0 });

  useEffect(() => {
    const m = getMenteeRequests().find(
      (r) => r.username === username && r.status === "approved"
    ) ?? null;
    setMentee(m);
    if (m) {
      setForm({ bio: m.bio ?? "", github: m.github ?? "", linkedin: m.linkedin ?? "" });
      const questions = getMenteeQuestions(username);
      setStats({
        solved: questions.filter((q) => q.status === "completed").length,
        pending: questions.filter((q) => q.status === "pending").length,
        total: questions.length,
      });
    }
  }, [username]);

  if (!mentee) return null;

  const initials = mentee.firstName[0] + (mentee.lastName?.[0] ?? "");

  const handleSave = () => {
    updateMenteeProfile(username, form);
    setMentee({ ...mentee, ...form });
    setEditing(false);
    setSaved(true);
    setTimeout(() => setSaved(false), 2000);
  };

  const handlePasswordUpdate = () => {
    setPwError("");
    if (!pwForm.next.trim()) { setPwError("New password cannot be empty."); return; }
    if (pwForm.next !== pwForm.confirm) { setPwError("Passwords do not match."); return; }
    if (pwForm.next.length < 6) { setPwError("Password must be at least 6 characters."); return; }
    const result = updateMenteePassword(username, pwForm.current, pwForm.next);
    if (!result.ok) { setPwError(result.error ?? "Failed to update password."); return; }
    setPwForm({ current: "", next: "", confirm: "" });
    setPwSaved(true);
    setTimeout(() => { setPwSaved(false); setShowPwCard(false); }, 2000);
  };

  return (
    <div className="max-w-xl">
      <h1 className="text-2xl font-bold text-purple-400 mb-6">My Profile</h1>

      {/* Avatar + name */}
      <div className="flex items-center gap-5 mb-8">
        <div className="w-20 h-20 rounded-full bg-purple-700 flex items-center justify-center text-3xl font-bold text-white shrink-0">
          {initials}
        </div>
        <div>
          <p className="text-xl font-bold text-gray-900 dark:text-white">
            {mentee.firstName} {mentee.lastName}
          </p>
          <p className="text-sm text-gray-500">@{mentee.username} · Mentee</p>
          <p className="text-xs text-gray-400 mt-0.5">
            Joined {new Date(mentee.signedUpAt).toLocaleDateString()}
          </p>
          {mentee.assignedSheet && (
            <span className="inline-block mt-1 text-xs font-medium text-purple-400 bg-purple-900/30 px-2 py-0.5 rounded-full">
              {SHEET_NAMES[mentee.assignedSheet] ?? mentee.assignedSheet}
            </span>
          )}
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-3 gap-3 mb-8">
        {[
          { label: "Solved", value: stats.solved, color: "text-green-400" },
          { label: "Pending", value: stats.pending, color: "text-yellow-400" },
          { label: "Total", value: stats.total, color: "text-purple-400" },
        ].map((s) => (
          <div
            key={s.label}
            className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 px-4 py-4 text-center"
          >
            <p className={`text-2xl font-bold ${s.color}`}>{s.value}</p>
            <p className="text-xs text-gray-500 mt-0.5">{s.label}</p>
          </div>
        ))}
      </div>

      {/* Editable details */}
      <div className="rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-6">
        <div className="flex items-center justify-between mb-5">
          <h2 className="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wider">
            Details
          </h2>
          {!editing && (
            <button
              onClick={() => setEditing(true)}
              className="text-xs font-semibold text-purple-400 hover:text-purple-300 transition-colors"
            >
              Edit
            </button>
          )}
        </div>

        <div className="flex flex-col gap-4">
          {/* Read-only fields */}
          {[
            { label: "First Name", value: mentee.firstName },
            { label: "Last Name", value: mentee.lastName },
            { label: "Username", value: `@${mentee.username}` },
            { label: "Email", value: mentee.email },
          ].map((f) => (
            <div key={f.label}>
              <label className="block text-xs font-semibold text-gray-400 mb-1">{f.label}</label>
              <p className="text-sm text-gray-800 dark:text-gray-200">{f.value}</p>
            </div>
          ))}

          {/* Editable fields */}
          <div>
            <label className="block text-xs font-semibold text-gray-400 mb-1">Bio</label>
            {editing ? (
              <textarea
                value={form.bio}
                onChange={(e) => setForm({ ...form, bio: e.target.value })}
                rows={3}
                placeholder="Tell others about yourself..."
                className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500 resize-none"
              />
            ) : (
              <p className="text-sm text-gray-800 dark:text-gray-200">
                {mentee.bio || <span className="text-gray-400">No bio yet</span>}
              </p>
            )}
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-400 mb-1">GitHub</label>
            {editing ? (
              <input
                value={form.github}
                onChange={(e) => setForm({ ...form, github: e.target.value })}
                placeholder="https://github.com/username"
                className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
              />
            ) : (
              <p className="text-sm">
                {mentee.github ? (
                  <a href={mentee.github} target="_blank" rel="noopener noreferrer" className="text-purple-400 hover:underline">
                    {mentee.github}
                  </a>
                ) : (
                  <span className="text-gray-400">—</span>
                )}
              </p>
            )}
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-400 mb-1">LinkedIn</label>
            {editing ? (
              <input
                value={form.linkedin}
                onChange={(e) => setForm({ ...form, linkedin: e.target.value })}
                placeholder="https://linkedin.com/in/username"
                className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
              />
            ) : (
              <p className="text-sm">
                {mentee.linkedin ? (
                  <a href={mentee.linkedin} target="_blank" rel="noopener noreferrer" className="text-purple-400 hover:underline">
                    {mentee.linkedin}
                  </a>
                ) : (
                  <span className="text-gray-400">—</span>
                )}
              </p>
            )}
          </div>
        </div>

        {editing && (
          <div className="flex gap-2 mt-6 justify-end">
            <button
              onClick={() => {
                setForm({ bio: mentee.bio ?? "", github: mentee.github ?? "", linkedin: mentee.linkedin ?? "" });
                setEditing(false);
              }}
              className="px-4 py-2 text-sm text-gray-500 hover:text-gray-700 dark:hover:text-gray-200 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              className="px-4 py-2 text-sm font-semibold bg-purple-600 hover:bg-purple-500 text-white rounded-lg transition-colors"
            >
              {saved ? "Saved ✓" : "Save Changes"}
            </button>
          </div>
        )}
      </div>

      {/* Update Password */}
      <div className="mt-6">
        <button
          onClick={() => { setShowPwCard((v) => !v); setPwError(""); setPwForm({ current: "", next: "", confirm: "" }); }}
          className="text-sm font-semibold text-purple-400 hover:text-purple-300 transition-colors"
        >
          {showPwCard ? "Cancel Password Update" : "Update Password"}
        </button>

        {showPwCard && (
          <div className="rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-6 mt-3">
            <h2 className="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wider mb-5">
              Update Password
            </h2>
            <div className="flex flex-col gap-4">
              {[
                { label: "Current Password", key: "current" },
                { label: "New Password", key: "next" },
                { label: "Confirm New Password", key: "confirm" },
              ].map(({ label, key }) => (
                <div key={key}>
                  <label className="block text-xs font-semibold text-gray-400 mb-1">{label}</label>
                  <input
                    type="password"
                    value={pwForm[key as keyof typeof pwForm]}
                    onChange={(e) => setPwForm({ ...pwForm, [key]: e.target.value })}
                    className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                  />
                </div>
              ))}
              {pwError && <p className="text-xs text-red-400">{pwError}</p>}
            </div>
            <div className="flex justify-end mt-5">
              <button
                onClick={handlePasswordUpdate}
                className={`px-4 py-2 text-sm font-semibold rounded-lg transition-colors ${
                  pwSaved ? "bg-green-600 text-white" : "bg-purple-600 hover:bg-purple-500 text-white"
                }`}
              >
                {pwSaved ? "Updated ✓" : "Update Password"}
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
