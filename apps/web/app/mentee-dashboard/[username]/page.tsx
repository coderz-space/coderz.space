import { redirect } from "next/navigation";

export default async function MenteeDashboardRoot({
  params,
}: {
  params: Promise<{ username: string }>;
}) {
  const { username } = await params;
  redirect(`/mentee-dashboard/${username}/pending`);
}
