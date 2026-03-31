"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";

const navItems = [
  { label: "Approve Mentee", href: "/mentor-dashboard/approve-mentee" },
  { label: "Assign Tasklist", href: "/mentor-dashboard/assign-tasklist" },
  { label: "Master Tasklist", href: "/mentor-dashboard/master-tasklist" },
  { label: "🏆 Leaderboard", href: "/mentor-dashboard/leaderboard" },
  { label: "👤 My Profile", href: "/mentor-dashboard/my-profile" },
];

export default function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();

  const handleLogout = () => {
    router.push("/");
  };

  return (
    <aside className="flex flex-col w-60 h-screen sticky top-0 bg-gray-100 dark:bg-gray-950 border-r border-purple-200 dark:border-purple-900 px-4 py-6 transition-colors">
      <div className="mb-8 px-2">
        <span className="text-purple-400 font-bold text-lg tracking-wide">Algo Buddy</span>
        <p className="text-gray-500 dark:text-gray-500 text-xs mt-1">Mentor Dashboard</p>
      </div>

      <nav className="flex flex-col gap-2 flex-1 overflow-y-auto">
        {navItems.map((item) => {
          const active = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`px-4 py-2.5 rounded-lg text-sm font-medium transition-colors ${
                active
                  ? "bg-purple-700 text-white"
                  : "text-gray-600 dark:text-gray-300 hover:bg-purple-100 dark:hover:bg-purple-900/50 hover:text-purple-900 dark:hover:text-white"
              }`}
            >
              {item.label}
            </Link>
          );
        })}
      </nav>

      <button
        onClick={handleLogout}
        className="mt-auto px-4 py-2.5 rounded-lg text-sm font-medium text-red-400 hover:bg-red-900/30 hover:text-red-300 transition-colors text-left"
      >
        Logout
      </button>
    </aside>
  );
}
