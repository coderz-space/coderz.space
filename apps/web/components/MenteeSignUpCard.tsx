"use client";

import { useState } from "react";
import { registerMentee } from "@/services";

interface MenteeSignUpCardProps {
  role: "mentor" | "mentee";
  onClose: () => void;
  onBackToLogin: () => void;
}

function EyeIcon({ visible }: { visible: boolean }) {
  return visible ? (
    <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-5 0-9-4-9-7s4-7 9-7a9.956 9.956 0 016.21 2.16M15 12a3 3 0 11-6 0 3 3 0 016 0zm6 0c0 3-4 7-9 7" />
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3l18 18" />
    </svg>
  ) : (
    <svg xmlns="http://www.w3.org/2000/svg" className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.477 0 8.268 2.943 9.542 7-1.274 4.057-5.065 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
    </svg>
  );
}

const inputClass =
  "w-full px-4 py-2 rounded-lg border border-purple-300 dark:border-purple-700 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500";

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

  const handleSignUp = () => {
    setError("");
    if (!firstName || !email || !password) { setError("Please fill all required fields."); return; }
    if (role === "mentee" && !username) { setError("Username is required."); return; }
    if (password !== confirmPassword) { setError("Passwords do not match."); return; }

    if (role === "mentee") {
      registerMentee({
        firstName,
        lastName,
        username: username.trim().toLowerCase(),
        email,
        passwordHash: password, // dev stub — backend should hash this
      });
    }

    setSubmitted(true);
  };

  return (
    <div
      className="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      onClick={onClose}
    >
      <div
        className="bg-white border border-purple-300 dark:bg-gray-900 dark:border-purple-700 rounded-2xl p-8 shadow-xl w-full max-w-sm mx-4 flex flex-col gap-4 max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-xl font-semibold text-center text-purple-700 dark:text-purple-400">
          {role === "mentor" ? "Mentor Sign Up" : "Mentee Sign Up"}
        </h2>

        {/* First & Last name */}
        <div className="flex gap-2">
          <input type="text" placeholder="First Name" value={firstName}
            onChange={(e) => setFirstName(e.target.value)}
            className="w-1/2 px-4 py-2 rounded-lg border border-purple-300 dark:border-purple-700 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500" />
          <input type="text" placeholder="Last Name" value={lastName}
            onChange={(e) => setLastName(e.target.value)}
            className="w-1/2 px-4 py-2 rounded-lg border border-purple-300 dark:border-purple-700 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500" />
        </div>

        {/* Username — mentee only */}
        {role === "mentee" && (
          <input type="text" placeholder="Username (used to log in)"
            value={username} onChange={(e) => setUsername(e.target.value)}
            className={inputClass} />
        )}

        <input type="email" placeholder="Email" value={email}
          onChange={(e) => setEmail(e.target.value)} className={inputClass} />

        {/* Password */}
        <div className="relative">
          <input type={showPassword ? "text" : "password"} placeholder="Set Password"
            value={password} onChange={(e) => setPassword(e.target.value)}
            className={`${inputClass} pr-10`} />
          <button type="button" onClick={() => setShowPassword((v) => !v)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 dark:text-gray-400 hover:text-purple-600"
            aria-label={showPassword ? "Hide password" : "Show password"}>
            <EyeIcon visible={showPassword} />
          </button>
        </div>

        {/* Confirm Password */}
        <div className="relative">
          <input type={showConfirm ? "text" : "password"} placeholder="Confirm Password"
            value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)}
            className={`${inputClass} pr-10`} />
          <button type="button" onClick={() => setShowConfirm((v) => !v)}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 dark:text-gray-400 hover:text-purple-600"
            aria-label={showConfirm ? "Hide confirm password" : "Show confirm password"}>
            <EyeIcon visible={showConfirm} />
          </button>
        </div>

        {error && <p className="text-red-500 text-sm text-center">{error}</p>}

        <button onClick={handleSignUp}
          className="bg-purple-600 hover:bg-purple-700 text-white w-full py-2 rounded-lg font-semibold">
          Sign Up
        </button>

        {submitted && (
          <p className="text-green-500 text-sm text-center">
            {role === "mentee"
              ? "Request sent! Await mentor approval."
              : "Account created! You can now log in."}
          </p>
        )}

        <button onClick={onBackToLogin}
          className="text-sm text-purple-600 dark:text-purple-400 text-center hover:underline">
          Already have an account? Log In
        </button>
      </div>
    </div>
  );
}
