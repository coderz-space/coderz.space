"use client";

import { useState, useEffect } from "react";
import { getMenteeRequests, updateMenteeStatus } from "@/services";
import { deleteMentee } from "@/services";
import Modal from "@/components/Modal";
import type { MenteeRequest, SheetId } from "@/types";

const SHEETS: { id: SheetId; name: string; desc: string }[] = [
  { id: "gfg-dsa-360", name: "GFG DSA 360", desc: "360 curated problems by GeeksForGeeks" },
  { id: "strivers-dsa-sheet", name: "Striver's DSA Sheet", desc: "Top DSA problems by Striver" },
];

export default function ApproveMenteePage() {
  const [requests, setRequests] = useState<MenteeRequest[]>([]);
  // pending approval: { menteeId, action }
  const [pendingApproval, setPendingApproval] = useState<{ id: string } | null>(null);

  useEffect(() => {
    const loadRequests = async () => {
      const reqs = await getMenteeRequests();
      setRequests(reqs);
    };
    loadRequests();
  }, []);

  const refresh = async () => {
    const reqs = await getMenteeRequests();
    setRequests(reqs);
  };

  const handleApprove = (id: string) => {
    setPendingApproval({ id });
  };

  const confirmApprove = async (sheet: SheetId) => {
    if (!pendingApproval) return;
    await updateMenteeStatus(pendingApproval.id, "approved", sheet);
    setPendingApproval(null);
    await refresh();
  };

  const handleReject = async (id: string) => {
    await updateMenteeStatus(id, "rejected");
    await refresh();
  };

  const handleDelete = async (id: string) => {
    await deleteMentee(id);
    await refresh();
  };

  const pending = requests.filter((r) => r.status === "pending");
  const handled = requests.filter((r) => r.status !== "pending");

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold text-purple-400 mb-6">Approve Mentee</h1>

      {pending.length === 0 && (
        <p className="text-gray-500 text-sm">No pending requests.</p>
      )}

      <div className="flex flex-col gap-4">
        {pending.map((req) => (
          <div
            key={req.id}
            className="bg-gray-100 dark:bg-gray-800 border border-purple-200 dark:border-purple-800 rounded-xl p-5 flex items-center justify-between gap-4"
          >
            <div>
              <p className="font-semibold text-gray-900 dark:text-white">
                {req.firstName} {req.lastName}
                <span className="ml-2 text-purple-500 text-sm font-normal">@{req.username}</span>
              </p>
              <p className="text-sm text-gray-500 dark:text-gray-400">{req.email}</p>
              <p className="text-xs text-gray-400 dark:text-gray-600 mt-1">
                Signed up {new Date(req.signedUpAt).toLocaleString()}
              </p>
            </div>
            <div className="flex gap-2 shrink-0">
              <button
                onClick={() => handleApprove(req.id)}
                className="px-4 py-1.5 rounded-lg bg-purple-600 hover:bg-purple-700 text-white text-sm font-medium transition-colors"
              >
                Approve
              </button>
              <button
                onClick={() => handleReject(req.id)}
                className="px-4 py-1.5 rounded-lg bg-gray-200 dark:bg-gray-700 hover:bg-red-100 dark:hover:bg-red-900/60 text-gray-600 dark:text-gray-300 hover:text-red-600 dark:hover:text-red-300 text-sm font-medium transition-colors"
              >
                Reject
              </button>
            </div>
          </div>
        ))}
      </div>

      {handled.length > 0 && (
        <div className="mt-10">
          <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-3">
            Handled
          </h2>
          <div className="flex flex-col gap-3">
            {handled.map((req) => (
              <div
                key={req.id}
                className="bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-xl p-4 flex items-center justify-between"
              >
                <div>
                  <p className="font-medium text-gray-700 dark:text-gray-300">
                    {req.firstName} {req.lastName}
                    <span className="ml-2 text-purple-500 text-sm font-normal">@{req.username}</span>
                  </p>
                  <p className="text-sm text-gray-500">{req.email}</p>
                  {req.assignedSheet && (
                    <p className="text-xs text-purple-400 mt-0.5">
                      Sheet: {SHEETS.find((s) => s.id === req.assignedSheet)?.name}
                    </p>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  <span
                    className={`text-xs font-semibold px-3 py-1 rounded-full ${
                      req.status === "approved"
                        ? "bg-green-100 dark:bg-green-900/50 text-green-600 dark:text-green-400"
                        : "bg-red-100 dark:bg-red-900/50 text-red-600 dark:text-red-400"
                    }`}
                  >
                    {req.status.charAt(0).toUpperCase() + req.status.slice(1)}
                  </span>
                  <button
                    onClick={() => handleDelete(req.id)}
                    className="px-3 py-1 rounded-lg bg-red-100 dark:bg-red-900/40 hover:bg-red-200 dark:hover:bg-red-800 text-red-500 dark:text-red-400 text-xs font-medium transition-colors"
                  >
                    Delete
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Sheet picker modal */}
      {pendingApproval && (
        <Modal onClose={() => setPendingApproval(null)} className="bg-white dark:bg-gray-900 border border-purple-200 dark:border-purple-800 rounded-2xl p-6 w-full max-w-md shadow-2xl">
            <h2 className="text-lg font-bold text-purple-400 mb-1">Assign a Sheet</h2>
            <p className="text-sm text-gray-500 mb-6">
              Choose the question sheet for this mentee. They will follow this sheet throughout.
            </p>

            <div className="flex flex-col gap-3 mb-6">
              {SHEETS.map((sheet) => (
                <button
                  key={sheet.id}
                  onClick={() => confirmApprove(sheet.id)}
                  className="flex flex-col items-start text-left px-5 py-4 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-purple-500 dark:hover:border-purple-500 hover:bg-purple-50 dark:hover:bg-purple-950/30 transition-colors group"
                >
                  <span className="font-semibold text-gray-800 dark:text-gray-100 group-hover:text-purple-600 dark:group-hover:text-purple-400 transition-colors">
                    {sheet.name}
                  </span>
                  <span className="text-xs text-gray-500 mt-0.5">{sheet.desc}</span>
                </button>
              ))}
            </div>

            <button
              onClick={() => setPendingApproval(null)}
              className="w-full text-sm text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors py-2"
            >
              Cancel
            </button>
        </Modal>
      )}
    </div>
  );
}
