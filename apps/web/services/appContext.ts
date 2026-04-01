import type { AppContext } from "@/types";
import { api } from "./api";

export async function getAppContext(): Promise<AppContext> {
  return api.get<AppContext>("/v1/app/context");
}
