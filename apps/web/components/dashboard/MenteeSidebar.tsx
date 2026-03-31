"use client";

import Link from "next/link";
import { usePathname, useRouter, useParams } from "next/navigation";

const navItems = (username: string) => [
  { label: "Pending Questions", href: `/mentee-dashboard/${username}/pending` },
  { label: "Completed Questions", href: `/mentee-dashboard/${username}/completed` },
  { label: "🏆 Leaderboard", href: `/mentee-dashboard/${username}/leaderboard` },
  { label: "👤 My Profile", href: `/mentee-dashboard/${username}/my-profile` },
];

export default function MenteeSidebar() {
  const pathname = usePathname();
  const router = useRouter();
  const params = useParams();
  const username = params?.username as string;

  return (
    <aside className="flex flex-col w-64 h-screen sticky top-0 bg-gray-100 dark:bg-gray-950 border-r border-purple-200 dark:border-purple-900 px-4 py-6 transition-colors">
      <div className="mb-8 px-2">
        <span className="text-purple-400 font-bold text-lg tracking-wide">Algo Buddy</span>
        <p className="text-gray-500 dark:text-gray-500 text-xs mt-1">@{username}</p>
      </div>

      <nav className="flex flex-col gap-2 flex-1 overflow-y-auto">
        {navItems(username).map((item) => {
          const active = pathname.startsWith(item.href);
          return (
            <Link key={item.href} href={item.href}
              className={`px-4 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                active
                  ? "bg-purple-700 text-white"
                  : "text-gray-600 dark:text-gray-300 hover:bg-purple-100 dark:hover:bg-purple-900/50 hover:text-purple-900 dark:hover:text-white"
              }`}>
              {item.label}
            </Link>
          );
        })}
      </nav>

      <button
        onClick={() => router.push("/")}
        className="mt-auto px-4 py-2.5 rounded-lg text-sm font-medium text-red-400 hover:bg-red-900/30 hover:text-red-300 transition-colors text-left"
      >
        Logout
      </button>
    </aside>
  );
}
