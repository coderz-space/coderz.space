import { redirect } from "next/navigation";

// Redirect /mentor-dashboard → /mentor-dashboard/approve-mentee by default
export default function MentorDashboardPage() {
  redirect("/mentor-dashboard/approve-mentee");
}
