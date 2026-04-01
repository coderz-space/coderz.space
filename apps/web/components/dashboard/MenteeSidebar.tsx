"use client";

import Link from "next/link";
import { useParams, usePathname, useRouter } from "next/navigation";
import { logout } from "@/services/auth";
import { clearSelectedRole } from "@/services/roleService";

const navItems = (username: string) => [
  { label: "Pending Questions", href: `/mentee-dashboard/${username}/pending` },
  { label: "Completed Questions", href: `/mentee-dashboard/${username}/completed` },
  { label: "Leaderboard", href: `/mentee-dashboard/${username}/leaderboard` },
  { label: "My Profile", href: `/mentee-dashboard/${username}/my-profile` },
];

export default function MenteeSidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const params = useParams();
  const username = params?.username as string;

  const handleLogout = async () => {
    await logout();
    clearSelectedRole();
    router.push("/");
  };

  return (
    <aside className="sticky top-0 flex h-screen w-64 flex-col border-r border-purple-200 bg-gray-100 px-4 py-6 transition-colors dark:border-purple-900 dark:bg-gray-950">
      <div className="mb-8 px-2">
        <span className="text-lg font-bold tracking-wide text-purple-400">Algo Buddy</span>
        <p className="mt-1 text-xs text-gray-500 dark:text-gray-500">@{username}</p>
      </div>

      <nav className="flex flex-1 flex-col gap-2 overflow-y-auto">
        {navItems(username).map((item) => {
          const active = pathname.startsWith(item.href);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`rounded-lg px-4 py-2.5 text-sm font-medium transition-colors ${
                active
                  ? "bg-purple-700 text-white"
                  : "text-gray-600 hover:bg-purple-100 hover:text-purple-900 dark:text-gray-300 dark:hover:bg-purple-900/50 dark:hover:text-white"
              }`}
            >
              {item.label}
            </Link>
          );
        })}
      </nav>

      <button
        type="button"
        onClick={handleLogout}
        className="mt-auto rounded-lg px-4 py-2.5 text-left text-sm font-medium text-red-400 transition-colors hover:bg-red-900/30 hover:text-red-300"
      >
        Logout
      </button>
    </aside>
  );
}
