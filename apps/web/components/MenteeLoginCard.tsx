"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { loginMenteeByEmail } from "@/services/auth";

interface MenteeLoginCardProps {
  role: "mentor" | "mentee";
  onClose: () => void;
  onSignUp: () => void;
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

export default function MenteeLoginCard({ role, onClose, onSignUp }: MenteeLoginCardProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleLogin = async () => {
    setError("");
    if (!email.trim() || !password) {
      setError("Please fill in all fields.");
      return;
    }

    setLoading(true);
    try {
      const result = await loginMenteeByEmail(email.trim(), password);
      if (result.context.accountStatus !== "approved") {
        setError("Your account is still pending mentor approval.");
        return;
      }

      if (result.context.role === "mentor") {
        router.push("/mentor-dashboard");
        return;
      }

      if (result.context.role === "mentee") {
        router.push(`/mentee-dashboard/${result.context.user.username}`);
        return;
      }

      setError("Your account is not linked to an active Algo Buddy role.");
    } catch (error) {
      setError(error instanceof Error ? error.message : "Invalid credentials.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 flex flex-col items-center justify-center bg-black/50 backdrop-blur-sm" onClick={onClose}>
      <div
        className="mx-4 flex w-full max-w-sm flex-col gap-4 rounded-2xl border border-purple-300 bg-white p-8 shadow-xl dark:border-purple-700 dark:bg-gray-900"
        onClick={(event) => event.stopPropagation()}
      >
        <h2 className="text-center text-xl font-semibold text-purple-700 dark:text-purple-400">
          {role === "mentor" ? "Mentor Login" : "Mentee Login"}
        </h2>

        <input
          type="email"
          placeholder="Email address"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          className={inputClass}
        />

        <div className="relative">
          <input
            type={showPassword ? "text" : "password"}
            placeholder="Password"
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

        <button className="self-end text-right text-sm text-purple-600 hover:underline dark:text-purple-400">
          Forgot Password?
        </button>

        {error ? <p className="text-center text-sm text-red-500">{error}</p> : null}

        <button
          onClick={handleLogin}
          disabled={loading}
          className="w-full rounded-lg bg-purple-600 py-2 font-semibold text-white hover:bg-purple-700 disabled:opacity-50"
        >
          {loading ? "Logging in..." : "Log In"}
        </button>

        <button className="flex w-full items-center justify-center gap-2 rounded-lg border border-purple-300 py-2 text-gray-700 hover:bg-purple-50 dark:border-purple-700 dark:text-gray-200 dark:hover:bg-gray-800">
          <svg className="h-5 w-5" viewBox="0 0 24 24" aria-hidden="true">
            <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
            <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
            <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" />
            <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
          </svg>
          Login with Google
        </button>
      </div>

      {role === "mentee" ? (
        <button
          onClick={(event) => {
            event.stopPropagation();
            onSignUp();
          }}
          className="mt-4 text-sm text-purple-600 hover:underline dark:text-purple-400"
        >
          New user? Sign Up
        </button>
      ) : null}
    </div>
  );
}
