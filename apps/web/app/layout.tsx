import type { Metadata } from "next";
import "./globals.css";
import ThemeToggle from "@/components/ThemeToggle";

export const metadata: Metadata = {
  title: "Algo Buddy",
  description:
    "A collaborative platform where developers master DSA and tech stacks through peer learning, problem-solving, and real progress tracking.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="h-full antialiased" suppressHydrationWarning>
      <head>
        <script
          dangerouslySetInnerHTML={{
            __html:
              "(function(){try{var t=localStorage.getItem('coderz_theme');var d=t?t==='dark':window.matchMedia('(prefers-color-scheme: dark)').matches;if(d)document.documentElement.classList.add('dark');}catch(e){}})();",
          }}
        />
      </head>
      <body className="flex min-h-full flex-col bg-white text-gray-900 transition-colors dark:bg-gray-950 dark:text-gray-100">
        <ThemeToggle />
        {children}
      </body>
    </html>
  );
}
