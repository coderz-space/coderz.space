"use client";

import { useState, useEffect } from "react";
import { getMentorProfile, saveMentorProfile, getLeaderboard, updateMentorPassword } from "@/services";
import type { MentorProfile } from "@/types";

export default function MentorMyProfilePage() {
  const [profile, setProfile] = useState<MentorProfile | null>(null);
  const [editing, setEditing] = useState(false);
  const [form, setForm] = useState<MentorProfile | null>(null);
  const [saved, setSaved] = useState(false);
  const [showPwCard, setShowPwCard] = useState(false);
  const [pwForm, setPwForm] = useState({ current: "", next: "", confirm: "" });
  const [pwError, setPwError] = useState("");
  const [pwSaved, setPwSaved] = useState(false);
  const [totalMentees, setTotalMentees] = useState(0);

  useEffect(() => {
    // Load mentor profile
    const loadProfile = async () => {
      const p = await getMentorProfile();
      setProfile(p);
      setForm(p);
    };
    loadProfile();

    // Load mentees count
    const loadMentees = async () => {
      const mentees = await getLeaderboard();
      setTotalMentees(mentees.length);
    };
    loadMentees();
  }, []);

  if (!profile || !form) return null;

  const initials =
    (profile.firstName?.[0] ?? "M") + (profile.lastName?.[0] ?? "");

  const handleSave = async () => {
    if (!form) return;
    try {
      await saveMentorProfile(form);
      setProfile(form);
      setEditing(false);
      setSaved(true);
      setTimeout(() => setSaved(false), 2000);
    } catch (error) {
      console.error("Failed to save profile:", error);
    }
  };

  const handlePasswordUpdate = async () => {
    setPwError("");
    if (!pwForm.next.trim()) { setPwError("New password cannot be empty."); return; }
    if (pwForm.next !== pwForm.confirm) { setPwError("Passwords do not match."); return; }
    if (pwForm.next.length < 6) { setPwError("Password must be at least 6 characters."); return; }
    
    try {
      const result = await updateMentorPassword(pwForm.current, pwForm.next);
      if (!result.ok) { setPwError(result.error ?? "Failed to update password."); return; }
      setPwForm({ current: "", next: "", confirm: "" });
      setPwSaved(true);
      setTimeout(() => { setPwSaved(false); setShowPwCard(false); }, 2000);
    } catch (error) {
      setPwError("An error occurred while updating password");
      console.error("Password update error:", error);
    }
  };

  return (
    <div className="max-w-xl">
      <h1 className="text-2xl font-bold text-purple-400 mb-6">My Profile</h1>

      {/* Avatar + name */}
      <div className="flex items-center gap-5 mb-8">
        <div className="w-20 h-20 rounded-full bg-purple-700 flex items-center justify-center text-3xl font-bold text-white shrink-0">
          {initials || "M"}
        </div>
        <div>
          <p className="text-xl font-bold text-gray-900 dark:text-white">
            {profile.firstName} {profile.lastName}
          </p>
          <p className="text-sm text-gray-500">@{profile.username} · Mentor</p>
          <p className="text-xs text-gray-400 mt-0.5">
            {profile.joinedAt ? `Joined ${new Date(profile.joinedAt).toLocaleDateString()}` : ""}
          </p>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 gap-4 mb-8">
        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 px-5 py-4">
          <p className="text-2xl font-bold text-purple-400">{totalMentees}</p>
          <p className="text-xs text-gray-500 mt-0.5">Active Mentees</p>
        </div>
        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 px-5 py-4">
          <p className="text-2xl font-bold text-purple-400">Mentor</p>
          <p className="text-xs text-gray-500 mt-0.5">Role</p>
        </div>
      </div>

      {/* Profile fields */}
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
          {(["firstName", "lastName", "username", "email"] as const).map((field) => (
            <div key={field}>
              <label className="block text-xs font-semibold text-gray-400 mb-1 capitalize">
                {field === "firstName" ? "First Name" : field === "lastName" ? "Last Name" : field}
              </label>
              {editing ? (
                <input
                  value={form[field] ?? ""}
                  onChange={(e) => setForm({ ...form, [field]: e.target.value })}
                  className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                />
              ) : (
                <p className="text-sm text-gray-800 dark:text-gray-200">
                  {profile[field] || <span className="text-gray-400">—</span>}
                </p>
              )}
            </div>
          ))}

          <div>
            <label className="block text-xs font-semibold text-gray-400 mb-1">Bio</label>
            {editing ? (
              <textarea
                value={form.bio ?? ""}
                onChange={(e) => setForm({ ...form, bio: e.target.value })}
                rows={3}
                placeholder="Tell your mentees about yourself..."
                className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500 resize-none"
              />
            ) : (
              <p className="text-sm text-gray-800 dark:text-gray-200">
                {profile.bio || <span className="text-gray-400">No bio yet</span>}
              </p>
            )}
          </div>

          <div>
            <label className="block text-xs font-semibold text-gray-400 mb-1">GitHub</label>
            {editing ? (
              <input
                value={form.github ?? ""}
                onChange={(e) => setForm({ ...form, github: e.target.value })}
                placeholder="https://github.com/username"
                className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
              />
            ) : (
              <p className="text-sm">
                {profile.github ? (
                  <a href={profile.github} target="_blank" rel="noopener noreferrer" className="text-purple-400 hover:underline">
                    {profile.github}
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
                value={form.linkedin ?? ""}
                onChange={(e) => setForm({ ...form, linkedin: e.target.value })}
                placeholder="https://linkedin.com/in/username"
                className="w-full bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-800 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
              />
            ) : (
              <p className="text-sm">
                {profile.linkedin ? (
                  <a href={profile.linkedin} target="_blank" rel="noopener noreferrer" className="text-purple-400 hover:underline">
                    {profile.linkedin}
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
              onClick={() => { setForm(profile); setEditing(false); }}
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
