"use client";

import { useState } from "react";
import { registerMentee } from "@/services/auth";

interface MenteeSignUpCardProps {
  role: "mentor" | "mentee";
  onClose: () => void;
  onBackToLogin: () => void;
}

function EyeIcon({ visible }: { visible: boolean }) {
  return visible ? (
    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-5 0-9-4-9-7s4-7 9-7a9.956 9.956 0 016.21 2.16M15 12a3 3 0 11-6 0 3 3 0 016 0zm6 0c0 3-4 7-9 7" />
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3l18 18" />
    </svg>
  ) : (
    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.477 0 8.268 2.943 9.542 7-1.274 4.057-5.065 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
    </svg>
  );
}

const inputClass =
  "w-full rounded-lg border border-purple-300 bg-white px-4 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-purple-500 dark:border-purple-700 dark:bg-gray-800 dark:text-gray-100";

export default function MenteeSignUpCard({ role, onClose, onBackToLogin }: MenteeSignUpCardProps) {
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSignUp = async () => {
    setError("");
    if (role !== "mentee") {
      setError("Mentor accounts are provisioned separately.");
      return;
    }

    if (!firstName.trim() || !username.trim() || !email.trim() || !password) {
      setError("Please fill all required fields.");
      return;
    }
    if (password.length < 8) {
      setError("Password must be at least 8 characters.");
      return;
    }
    if (password !== confirmPassword) {
      setError("Passwords do not match.");
      return;
    }

    setLoading(true);
    try {
      await registerMentee({
        firstName: firstName.trim(),
        lastName: lastName.trim(),
        username: username.trim().toLowerCase(),
        email: email.trim(),
        password,
      });
      setSubmitted(true);
    } catch (error) {
      setError(error instanceof Error ? error.message : "Registration failed.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm" onClick={onClose}>
      <div
        className="mx-4 flex max-h-[90vh] w-full max-w-sm flex-col gap-4 overflow-y-auto rounded-2xl border border-purple-300 bg-white p-8 shadow-xl dark:border-purple-700 dark:bg-gray-900"
        onClick={(event) => event.stopPropagation()}
      >
        <h2 className="text-center text-xl font-semibold text-purple-700 dark:text-purple-400">Mentee Sign Up</h2>

        <div className="flex gap-2">
          <input
            type="text"
            placeholder="First Name"
            value={firstName}
            onChange={(event) => setFirstName(event.target.value)}
            className="w-1/2 rounded-lg border border-purple-300 bg-white px-4 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-purple-500 dark:border-purple-700 dark:bg-gray-800 dark:text-gray-100"
          />
          <input
            type="text"
            placeholder="Last Name"
            value={lastName}
            onChange={(event) => setLastName(event.target.value)}
            className="w-1/2 rounded-lg border border-purple-300 bg-white px-4 py-2 text-gray-900 focus:outline-none focus:ring-2 focus:ring-purple-500 dark:border-purple-700 dark:bg-gray-800 dark:text-gray-100"
          />
        </div>

        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(event) => setUsername(event.target.value)}
          className={inputClass}
        />

        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          className={inputClass}
        />

        <div className="relative">
          <input
            type={showPassword ? "text" : "password"}
            placeholder="Set Password (min 8 chars)"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            className={`${inputClass} pr-10`}
          />
          <button
            type="button"
            onClick={() => setShowPassword((value) => !value)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-purple-600 dark:text-gray-400"
            aria-label={showPassword ? "Hide password" : "Show password"}
          >
            <EyeIcon visible={showPassword} />
          </button>
        </div>

        <div className="relative">
          <input
            type={showConfirm ? "text" : "password"}
            placeholder="Confirm Password"
            value={confirmPassword}
            onChange={(event) => setConfirmPassword(event.target.value)}
            className={`${inputClass} pr-10`}
          />
          <button
            type="button"
            onClick={() => setShowConfirm((value) => !value)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-purple-600 dark:text-gray-400"
            aria-label={showConfirm ? "Hide confirm password" : "Show confirm password"}
          >
            <EyeIcon visible={showConfirm} />
          </button>
        </div>

        {error ? <p className="text-center text-sm text-red-500">{error}</p> : null}
        {submitted ? (
          <p className="text-center text-sm text-green-500">
            Sign-up request submitted. A mentor will approve your account before you can log in.
          </p>
        ) : null}

        <button
          onClick={handleSignUp}
          disabled={loading}
          className="w-full rounded-lg bg-purple-600 py-2 font-semibold text-white hover:bg-purple-700 disabled:opacity-50"
        >
          {loading ? "Signing up..." : "Sign Up"}
        </button>

        <button onClick={onBackToLogin} className="text-center text-sm text-purple-600 hover:underline dark:text-purple-400">
          Already have an account? Log In
        </button>
      </div>
    </div>
  );
}
