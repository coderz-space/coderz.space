"use client";

import { useEffect, useState, useTransition } from "react";
import type { MentorProfile } from "@/types";
import { formatDate, getInitials } from "./constants";

type ProfileStat = {
  label: string;
  value: string | number;
};

type AccountProfileEditorProps = {
  title: string;
  roleLabel: string;
  profile: MentorProfile | null;
  stats: ProfileStat[];
  onSave: (profile: MentorProfile) => Promise<MentorProfile>;
  onUpdatePassword: (currentPassword: string, newPassword: string) => Promise<{ ok: boolean; error?: string }>;
};

const textInputClass =
  "w-full rounded-xl border border-gray-200 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition focus:border-purple-500 focus:ring-2 focus:ring-purple-200 dark:border-gray-700 dark:bg-gray-950 dark:text-gray-100 dark:focus:border-purple-500 dark:focus:ring-purple-900/60";

export default function AccountProfileEditor({
  title,
  roleLabel,
  profile,
  stats,
  onSave,
  onUpdatePassword,
}: AccountProfileEditorProps) {
  const [form, setForm] = useState<MentorProfile | null>(profile);
  const [editing, setEditing] = useState(false);
  const [saveMessage, setSaveMessage] = useState("");
  const [showPasswordCard, setShowPasswordCard] = useState(false);
  const [passwordForm, setPasswordForm] = useState({ current: "", next: "", confirm: "" });
  const [passwordMessage, setPasswordMessage] = useState("");
  const [savePending, startSaveTransition] = useTransition();
  const [passwordPending, startPasswordTransition] = useTransition();

  useEffect(() => {
    setForm(profile);
  }, [profile]);

  if (!profile || !form) {
    return <p className="text-sm text-gray-500 dark:text-gray-400">Loading profile...</p>;
  }

  const handleSave = () => {
    startSaveTransition(async () => {
      try {
        const updated = await onSave(form);
        setForm(updated);
        setEditing(false);
        setSaveMessage("Profile updated.");
      } catch (error) {
        setSaveMessage(error instanceof Error ? error.message : "Failed to update profile.");
      }
    });
  };

  const handlePasswordUpdate = () => {
    setPasswordMessage("");

    if (!passwordForm.next.trim()) {
      setPasswordMessage("New password cannot be empty.");
      return;
    }
    if (passwordForm.next !== passwordForm.confirm) {
      setPasswordMessage("Passwords do not match.");
      return;
    }

    startPasswordTransition(async () => {
      const result = await onUpdatePassword(passwordForm.current, passwordForm.next);
      if (!result.ok) {
        setPasswordMessage(result.error ?? "Failed to update password.");
        return;
      }

      setPasswordForm({ current: "", next: "", confirm: "" });
      setPasswordMessage("Password updated.");
      setShowPasswordCard(false);
    });
  };

  return (
    <div className="max-w-3xl">
      <h1 className="mb-6 text-2xl font-bold text-gray-900 dark:text-white">{title}</h1>

      <section className="mb-8 rounded-3xl border border-purple-200 bg-white p-6 shadow-sm dark:border-purple-900/60 dark:bg-gray-900">
        <div className="mb-5 flex items-center gap-5">
          <div className="flex h-20 w-20 shrink-0 items-center justify-center rounded-full bg-purple-600 text-3xl font-bold text-white">
            {getInitials(profile.firstName, profile.lastName, "AB")}
          </div>

          <div>
            <p className="text-xl font-bold text-gray-900 dark:text-white">
              {profile.firstName} {profile.lastName}
            </p>
            <p className="text-sm text-gray-500 dark:text-gray-400">
              @{profile.username} · {roleLabel}
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400">Joined {formatDate(profile.joinedAt)}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {stats.map((stat) => (
            <div
              key={stat.label}
              className="rounded-2xl border border-gray-200 bg-gray-50 px-5 py-4 dark:border-gray-800 dark:bg-gray-950/60"
            >
              <p className="text-2xl font-bold text-purple-600 dark:text-purple-300">{stat.value}</p>
              <p className="text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">{stat.label}</p>
            </div>
          ))}
        </div>
      </section>

      <section className="rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-800 dark:bg-gray-900">
        <div className="mb-5 flex items-center justify-between">
          <h2 className="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">Details</h2>
          {editing ? null : (
            <button
              type="button"
              onClick={() => {
                setEditing(true);
                setSaveMessage("");
              }}
              className="text-sm font-semibold text-purple-600 transition-colors hover:text-purple-500 dark:text-purple-300"
            >
              Edit
            </button>
          )}
        </div>

        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          {([
            ["firstName", "First Name"],
            ["lastName", "Last Name"],
            ["username", "Username"],
            ["email", "Email"],
          ] as const).map(([field, label]) => (
            <div key={field}>
              <label className="mb-1 block text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                {label}
              </label>
              {editing ? (
                <input
                  value={form[field]}
                  onChange={(event) => setForm({ ...form, [field]: event.target.value })}
                  className={textInputClass}
                />
              ) : (
                <p className="rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-900 dark:border-gray-800 dark:bg-gray-950/60 dark:text-gray-100">
                  {profile[field] || "-"}
                </p>
              )}
            </div>
          ))}
        </div>

        <div className="mt-4">
          <label className="mb-1 block text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
            Bio
          </label>
          {editing ? (
            <textarea
              value={form.bio ?? ""}
              onChange={(event) => setForm({ ...form, bio: event.target.value })}
              rows={4}
              className={textInputClass}
              placeholder="Tell the cohort a bit about yourself."
            />
          ) : (
            <p className="rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm leading-6 text-gray-900 dark:border-gray-800 dark:bg-gray-950/60 dark:text-gray-100">
              {profile.bio || "No bio yet."}
            </p>
          )}
        </div>

        <div className="mt-4 grid grid-cols-1 gap-4 md:grid-cols-2">
          {([
            ["github", "GitHub", "https://github.com/username"],
            ["linkedin", "LinkedIn", "https://linkedin.com/in/username"],
          ] as const).map(([field, label, placeholder]) => (
            <div key={field}>
              <label className="mb-1 block text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                {label}
              </label>
              {editing ? (
                <input
                  value={form[field] ?? ""}
                  onChange={(event) => setForm({ ...form, [field]: event.target.value })}
                  placeholder={placeholder}
                  className={textInputClass}
                />
              ) : profile[field] ? (
                <a
                  href={profile[field]}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="block rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-purple-600 hover:underline dark:border-gray-800 dark:bg-gray-950/60 dark:text-purple-300"
                >
                  {profile[field]}
                </a>
              ) : (
                <p className="rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-500 dark:border-gray-800 dark:bg-gray-950/60 dark:text-gray-400">
                  -
                </p>
              )}
            </div>
          ))}
        </div>

        {saveMessage ? (
          <p className="mt-4 text-sm text-gray-600 dark:text-gray-300">{saveMessage}</p>
        ) : null}

        {editing ? (
          <div className="mt-6 flex justify-end gap-2">
            <button
              type="button"
              onClick={() => {
                setForm(profile);
                setEditing(false);
                setSaveMessage("");
              }}
              className="rounded-lg px-4 py-2 text-sm font-semibold text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-gray-300 dark:hover:bg-gray-800 dark:hover:text-white"
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={handleSave}
              disabled={savePending}
              className="rounded-lg bg-purple-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-purple-500 disabled:opacity-60"
            >
              {savePending ? "Saving..." : "Save Changes"}
            </button>
          </div>
        ) : null}
      </section>

      <section className="mt-6 rounded-3xl border border-gray-200 bg-white p-6 shadow-sm dark:border-gray-800 dark:bg-gray-900">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
            Password
          </h2>
          <button
            type="button"
            onClick={() => {
              setShowPasswordCard((value) => !value);
              setPasswordForm({ current: "", next: "", confirm: "" });
              setPasswordMessage("");
            }}
            className="text-sm font-semibold text-purple-600 transition-colors hover:text-purple-500 dark:text-purple-300"
          >
            {showPasswordCard ? "Cancel" : "Update Password"}
          </button>
        </div>

        {showPasswordCard ? (
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div>
              <label className="mb-1 block text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                Current Password
              </label>
              <input
                type="password"
                value={passwordForm.current}
                onChange={(event) => setPasswordForm({ ...passwordForm, current: event.target.value })}
                className={textInputClass}
              />
            </div>
            <div>
              <label className="mb-1 block text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                New Password
              </label>
              <input
                type="password"
                value={passwordForm.next}
                onChange={(event) => setPasswordForm({ ...passwordForm, next: event.target.value })}
                className={textInputClass}
              />
            </div>
            <div className="md:col-span-2">
              <label className="mb-1 block text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
                Confirm New Password
              </label>
              <input
                type="password"
                value={passwordForm.confirm}
                onChange={(event) => setPasswordForm({ ...passwordForm, confirm: event.target.value })}
                className={textInputClass}
              />
            </div>

            {passwordMessage ? (
              <p className="md:col-span-2 text-sm text-gray-600 dark:text-gray-300">{passwordMessage}</p>
            ) : null}

            <div className="md:col-span-2 flex justify-end">
              <button
                type="button"
                onClick={handlePasswordUpdate}
                disabled={passwordPending}
                className="rounded-lg bg-purple-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-purple-500 disabled:opacity-60"
              >
                {passwordPending ? "Updating..." : "Update Password"}
              </button>
            </div>
          </div>
        ) : (
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Keep your account secure with a strong password.
          </p>
        )}
      </section>
    </div>
  );
}
