"use client";

import { useState } from "react";
import type { Role } from "@/types";
import HeroSection from "@/components/HeroSection";
import RoleCard from "@/components/RoleCard";
import MenteeLoginCard from "@/components/MenteeLoginCard";
import MenteeSignUpCard from "@/components/MenteeSignUpCard";
import { selectRole } from "@/services/roleService";

type View = "none" | "roleSelect" | "login" | "signUp";

export default function LandingPage() {
  const [view, setView] = useState<View>("none");
  const [activeRole, setActiveRole] = useState<Role | null>(null);

  const handleSelectRole = async (role: Role) => {
    await selectRole(role);
    setActiveRole(role);
    setView("login");
  };

  return (
    <main className="min-h-screen dark:bg-black bg-white">
      <HeroSection onGetStarted={() => setView("roleSelect")} />

      {view === "roleSelect" && (
        <RoleCard onSelectRole={handleSelectRole} onClose={() => setView("none")} />
      )}

      {view === "login" && activeRole && (
        <MenteeLoginCard
          role={activeRole}
          onClose={() => setView("none")}
          onSignUp={() => setView("signUp")}
        />
      )}

      {view === "signUp" && activeRole && (
        <MenteeSignUpCard
          role={activeRole}
          onClose={() => setView("none")}
          onBackToLogin={() => setView("login")}
        />
      )}
    </main>
  );
}
