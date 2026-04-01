import type { Question, QuestionProgressStatus } from "@/types";
import { api } from "./api";

export async function getMenteeQuestions(username: string): Promise<Question[]> {
  return api.get<Question[]>(`/v1/app/mentees/${username}/questions`);
}

export async function getQuestionDetail(username: string, questionId: string): Promise<Question | null> {
  try {
    return await api.get<Question>(`/v1/app/mentees/${username}/questions/${questionId}`);
  } catch {
    return null;
  }
}

export async function updateQuestionProgress(
  username: string,
  questionId: string,
  progressStatus: QuestionProgressStatus
): Promise<Question> {
  return api.patch<Question, { progressStatus: QuestionProgressStatus }>(
    `/v1/app/mentees/${username}/questions/${questionId}`,
    { progressStatus }
  );
}

export async function updateQuestionDetails(
  username: string,
  questionId: string,
  details: { solution?: string; resources?: string }
): Promise<Question> {
  return api.patch<Question, { solution?: string; resources?: string }>(
    `/v1/app/mentees/${username}/questions/${questionId}`,
    details
  );
}
