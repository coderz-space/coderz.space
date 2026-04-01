import type {
  DayAssignments,
  MenteeRequest,
  Sheet,
  SheetId,
} from "@/types";
import { api } from "./api";

export async function getMenteeRequests(): Promise<MenteeRequest[]> {
  return api.get<MenteeRequest[]>("/v1/app/mentor/mentee-requests");
}

export async function updateMenteeStatus(
  requestId: string,
  status: "approved" | "rejected",
  assignedSheet?: SheetId
): Promise<MenteeRequest> {
  return api.patch<MenteeRequest, { status: "approved" | "rejected"; sheetKey?: SheetId }>(
    `/v1/app/mentor/mentee-requests/${requestId}`,
    {
      status,
      sheetKey: assignedSheet,
    }
  );
}

export async function getSheets(): Promise<Sheet[]> {
  return api.get<Sheet[]>("/v1/app/sheets");
}

export async function getDayAssignments(day: string): Promise<DayAssignments> {
  return api.get<DayAssignments>(`/v1/app/mentor/day-assignments/${day}`);
}

export async function updateDayAssignments(day: string, usernames: string[]): Promise<DayAssignments> {
  return api.put<DayAssignments, { usernames: string[] }>(
    `/v1/app/mentor/day-assignments/${day}`,
    { usernames }
  );
}

export async function createAssignments(input: {
  day?: string;
  menteeUsernames: string[];
  sheetKey: SheetId;
  questionIds: string[];
}): Promise<{ assignmentGroupId: string; assignmentsCount: number }> {
  return api.post<{ assignmentGroupId: string; assignmentsCount: number }, typeof input>(
    "/v1/app/mentor/assignments",
    input
  );
}
