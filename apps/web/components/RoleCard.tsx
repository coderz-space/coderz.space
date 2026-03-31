import { RoleCardProps } from "@/types";

export default function RoleCard({ onSelectRole, onClose }: RoleCardProps) {
  return (
    <div
      className="fixed inset-0 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      onClick={onClose}
    >
      <div
        className="bg-white border border-purple-300 dark:bg-gray-900 dark:border-purple-700 rounded-2xl p-8 shadow-xl w-full max-w-sm mx-4 flex flex-col gap-4"
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-xl font-semibold text-center text-purple-700 dark:text-purple-400">
          Log in as
        </h2>
        <button
          onClick={() => onSelectRole("mentor")}
          className="bg-purple-600 hover:bg-purple-700 text-white w-full py-2 rounded-lg"
        >
          Mentor
        </button>
        <button
          onClick={() => onSelectRole("mentee")}
          className="bg-purple-600 hover:bg-purple-700 text-white w-full py-2 rounded-lg"
        >
          Mentee
        </button>
      </div>
    </div>
  );
}
