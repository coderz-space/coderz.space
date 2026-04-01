import type { LeaderboardEntry } from "@/types";
import { api } from "./api";

export async function getLeaderboard(): Promise<LeaderboardEntry[]> {
  return api.get<LeaderboardEntry[]>("/v1/app/leaderboard");
}
