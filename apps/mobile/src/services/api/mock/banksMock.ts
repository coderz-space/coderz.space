import { Problem } from '../../../types';
import { MOCK_PROBLEMS } from './_mockData';

export interface QuestionBank {
  id: string;
  name: string;
}

const banks: QuestionBank[] = [
  { id: 'gfg', name: 'GFG (GeeksforGeeks)' },
  { id: 'striver', name: 'Striver A2Z Sheet' },
  { id: 'leetcode', name: 'LeetCode Top 100' },
];

// Map bank id to problem ids (mock)
const bankProblemsMap: Record<string, string[]> = {
  gfg: ['prob-1', 'prob-2', 'prob-3'],
  striver: ['prob-1', 'prob-4', 'prob-5'],
  leetcode: ['prob-2', 'prob-3', 'prob-4'],
};

export const getQuestionBanks = (): QuestionBank[] => {
  return banks;
};

export const getProblemsByBank = async (bankId: string): Promise<Problem[]> => {
  // Simulate network delay
  await new Promise<void>((r) => setTimeout(() => r()));
  const problemIds = bankProblemsMap[bankId] || [];
  return MOCK_PROBLEMS.filter(p => problemIds.includes(p.id));
};