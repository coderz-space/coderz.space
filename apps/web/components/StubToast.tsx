"use client";

import { useState, useEffect, useCallback, createContext, useContext } from "react";

type Toast = {
  id: number;
  message: string;
};

type StubToastContextType = {
  showStubToast: (feature: string) => void;
};

const StubToastContext = createContext<StubToastContextType>({
  showStubToast: () => {},
});

let toastIdCounter = 0;

/**
 * Hook to trigger "not implemented" toasts from any component
 */
export function useStubToast() {
  return useContext(StubToastContext);
}

/**
 * Standalone function to show a stub toast without needing React context.
 * Used in service layer functions that aren't inside React components.
 */
export function showStubNotification(feature: string): void {
  if (typeof window === "undefined") return;
  const event = new CustomEvent("stub-toast", {
    detail: `⚠️ "${feature}" — backend not implemented yet`,
  });
  window.dispatchEvent(event);
}

/**
 * Provider component — wrap your app with this to enable stub toasts.
 * Listens for both context calls and custom DOM events (for service layer usage).
 */
export function StubToastProvider({ children }: { children: React.ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const addToast = useCallback((message: string) => {
    const id = ++toastIdCounter;
    setToasts((prev) => [...prev, { id, message }]);
    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 4000);
  }, []);

  const showStubToast = useCallback(
    (feature: string) => {
      addToast(`⚠️ "${feature}" — backend not implemented yet`);
    },
    [addToast]
  );

  // Listen for events dispatched from service layer (outside React tree)
  useEffect(() => {
    const handler = (e: Event) => {
      const msg = (e as CustomEvent).detail as string;
      addToast(msg);
    };
    window.addEventListener("stub-toast", handler);
    return () => window.removeEventListener("stub-toast", handler);
  }, [addToast]);

  return (
    <StubToastContext.Provider value={{ showStubToast }}>
      {children}

      {/* Toast container — fixed bottom-right */}
      {toasts.length > 0 && (
        <div className="fixed bottom-6 right-6 z-[9999] flex flex-col gap-2 pointer-events-none">
          {toasts.map((t) => (
            <div
              key={t.id}
              className="pointer-events-auto animate-slide-in-right bg-yellow-900/90 border border-yellow-600/60 text-yellow-200 px-5 py-3 rounded-xl shadow-2xl backdrop-blur-sm text-sm font-medium max-w-sm"
              style={{
                animation: "slideInRight 0.3s ease-out",
              }}
            >
              {t.message}
            </div>
          ))}
        </div>
      )}

      {/* Inline keyframes */}
      <style>{`
        @keyframes slideInRight {
          from { opacity: 0; transform: translateX(100px); }
          to   { opacity: 1; transform: translateX(0); }
        }
      `}</style>
    </StubToastContext.Provider>
  );
}
