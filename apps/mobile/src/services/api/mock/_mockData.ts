import {
  User, Assignment, AssignmentProblem, Problem, OrgMember,
  LeaderboardEntry, Doubt, AssignmentGroup, Tag,
} from '../../../types';

export const MOCK_USER_MENTEE: User = {
  id: 'user-mentee-1',
  name: 'Arjun Sharma',
  email: 'arjun@coderz.space',
  role: 'mentee',
  emailVerified: true,
  createdAt: '2024-01-15T00:00:00Z',
};

export const MOCK_USER_MENTOR: User = {
  id: 'user-mentor-1',
  name: 'Priya Singh',
  email: 'priya@coderz.space',
  role: 'mentor',
  emailVerified: true,
  createdAt: '2023-11-01T00:00:00Z',
};

const MOCK_TAGS: Tag[] = [
  { id: 'tag-1', name: 'arrays', organizationId: 'org-1' },
  { id: 'tag-2', name: 'dynamic-programming', organizationId: 'org-1' },
  { id: 'tag-3', name: 'graphs', organizationId: 'org-1' },
  { id: 'tag-4', name: 'binary-search', organizationId: 'org-1' },
  { id: 'tag-5', name: 'sliding-window', organizationId: 'org-1' },
];

export const MOCK_PROBLEMS: Problem[] = [
  {
    id: 'prob-1',
    organizationId: 'org-1',
    createdBy: 'orgmember-mentor-1',
    title: 'Two Sum',
    description: 'Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.',
    difficulty: 'easy',
    externalLink: 'https://leetcode.com/problems/two-sum/',
    tags: [MOCK_TAGS[0]],
    resources: [{ id: 'res-1', problemId: 'prob-1', title: 'Array Hash Map Approach', url: 'https://youtube.com' }],
    createdAt: '2024-01-10T00:00:00Z',
  },
  {
    id: 'prob-2',
    organizationId: 'org-1',
    createdBy: 'orgmember-mentor-1',
    title: 'Longest Substring Without Repeating Characters',
    description: 'Given a string s, find the length of the longest substring without repeating characters.',
    difficulty: 'medium',
    externalLink: 'https://leetcode.com/problems/longest-substring-without-repeating-characters/',
    tags: [MOCK_TAGS[4]],
    resources: [],
    createdAt: '2024-01-10T00:00:00Z',
  },
  {
    id: 'prob-3',
    organizationId: 'org-1',
    createdBy: 'orgmember-mentor-1',
    title: 'LRU Cache',
    description: 'Design a data structure that follows the constraints of a Least Recently Used (LRU) cache.',
    difficulty: 'hard',
    externalLink: 'https://leetcode.com/problems/lru-cache/',
    tags: [MOCK_TAGS[0]],
    resources: [],
    createdAt: '2024-01-10T00:00:00Z',
  },
  {
    id: 'prob-4',
    organizationId: 'org-1',
    createdBy: 'orgmember-mentor-1',
    title: 'Number of Islands',
    description: 'Given an m x n 2D binary grid, return the number of islands.',
    difficulty: 'medium',
    externalLink: 'https://leetcode.com/problems/number-of-islands/',
    tags: [MOCK_TAGS[2]],
    resources: [],
    createdAt: '2024-01-10T00:00:00Z',
  },
  {
    id: 'prob-5',
    organizationId: 'org-1',
    createdBy: 'orgmember-mentor-1',
    title: 'Binary Search',
    description: 'Given an array of integers nums which is sorted in ascending order, and an integer target, write a function to search target in nums.',
    difficulty: 'easy',
    externalLink: 'https://leetcode.com/problems/binary-search/',
    tags: [MOCK_TAGS[3]],
    resources: [],
    createdAt: '2024-01-10T00:00:00Z',
  },
];

const makeAssignmentProblems = (assignmentId: string): AssignmentProblem[] => [
  {
    id: 'ap-1',
    assignmentId,
    problemId: 'prob-1',
    problem: MOCK_PROBLEMS[0],
    status: 'completed',
    menteeStatus: 'completed',
    solutionLink: 'https://github.com/arjun/two-sum',
    notes: 'Used hashmap approach O(n)',
    remarkForSelf: 'Need to practice more hash problems',
    remarkForMentor: '',
    completedAt: new Date(Date.now() - 86400000).toISOString(),
  },
  {
    id: 'ap-2',
    assignmentId,
    problemId: 'prob-2',
    problem: MOCK_PROBLEMS[1],
    status: 'attempted',
    menteeStatus: 'discussion_needed',
    solutionLink: '',
    notes: '',
    remarkForSelf: '',
    remarkForMentor: 'Not sure about the window shrink condition',
    doubt: {
      id: 'doubt-1',
      assignmentProblemId: 'ap-2',
      raisedBy: 'orgmember-mentee-1',
      message: 'Not sure about the window shrink condition',
      resolved: false,
      createdAt: new Date().toISOString(),
    },
  },
  {
    id: 'ap-3',
    assignmentId,
    problemId: 'prob-3',
    problem: MOCK_PROBLEMS[2],
    status: 'pending',
    menteeStatus: 'not_started',
    solutionLink: '',
    notes: '',
    remarkForSelf: '',
    remarkForMentor: '',
  },
];

export const MOCK_ASSIGNMENT_GROUP: AssignmentGroup = {
  id: 'ag-1',
  bootcampId: 'bootcamp-1',
  createdBy: 'orgmember-mentor-1',
  title: 'Week 1 — Arrays & Sliding Window',
  description: 'Foundation problems covering arrays, hashmaps, and the sliding window pattern.',
  deadlineDays: 7,
  createdAt: '2024-03-01T00:00:00Z',
  problems: MOCK_PROBLEMS.slice(0, 3),
};

export const MOCK_ASSIGNMENT_GROUP_2: AssignmentGroup = {
  id: 'ag-2',
  bootcampId: 'bootcamp-1',
  createdBy: 'orgmember-mentor-1',
  title: 'Week 2 — Graphs & BFS/DFS',
  description: 'Graph traversal fundamentals.',
  deadlineDays: 7,
  createdAt: '2024-03-08T00:00:00Z',
  problems: MOCK_PROBLEMS.slice(3),
};

export const MOCK_ACTIVE_ASSIGNMENT: Assignment = {
  id: 'assign-1',
  assignmentGroupId: 'ag-1',
  bootcampEnrollmentId: 'enrollment-mentee-1',
  assignedBy: 'orgmember-mentor-1',
  assignedAt: new Date(Date.now() - 86400000 * 2).toISOString(),
  deadlineAt: new Date(Date.now() + 86400000 * 5).toISOString(),
  status: 'active',
  assignmentGroup: MOCK_ASSIGNMENT_GROUP,
  problems: makeAssignmentProblems('assign-1'),
  totalProblems: 3,
  completedProblems: 1,
  progressPercent: 33,
};

export const MOCK_COMPLETED_ASSIGNMENT: Assignment = {
  id: 'assign-0',
  assignmentGroupId: 'ag-0',
  bootcampEnrollmentId: 'enrollment-mentee-1',
  assignedBy: 'orgmember-mentor-1',
  assignedAt: new Date(Date.now() - 86400000 * 14).toISOString(),
  deadlineAt: new Date(Date.now() - 86400000 * 7).toISOString(),
  status: 'completed',
  assignmentGroup: {
    ...MOCK_ASSIGNMENT_GROUP,
    id: 'ag-0',
    title: 'Week 0 — Warm Up',
  },
  problems: [
    {
      id: 'ap-c-1',
      assignmentId: 'assign-0',
      problemId: 'prob-5',
      problem: MOCK_PROBLEMS[4],
      status: 'completed',
      menteeStatus: 'completed',
      solutionLink: 'https://github.com/arjun/binary-search',
      notes: '',
      remarkForSelf: '',
      remarkForMentor: '',
      completedAt: new Date(Date.now() - 86400000 * 8).toISOString(),
    },
  ],
  totalProblems: 1,
  completedProblems: 1,
  progressPercent: 100,
};

export const MOCK_MENTEES: OrgMember[] = [
  {
    id: 'orgmember-mentee-1',
    userId: 'user-mentee-1',
    organizationId: 'org-1',
    role: 'mentee',
    joinedAt: '2024-01-15T00:00:00Z',
    user: { id: 'user-mentee-1', name: 'Arjun Sharma', email: 'arjun@coderz.space', avatarUrl: undefined },
  },
  {
    id: 'orgmember-mentee-2',
    userId: 'user-mentee-2',
    organizationId: 'org-1',
    role: 'mentee',
    joinedAt: '2024-01-16T00:00:00Z',
    user: { id: 'user-mentee-2', name: 'Rahul Verma', email: 'rahul@coderz.space', avatarUrl: undefined },
  },
  {
    id: 'orgmember-mentee-3',
    userId: 'user-mentee-3',
    organizationId: 'org-1',
    role: 'mentee',
    joinedAt: '2024-01-17T00:00:00Z',
    user: { id: 'user-mentee-3', name: 'Sneha Patel', email: 'sneha@coderz.space', avatarUrl: undefined },
  },
];

export const MOCK_LEADERBOARD: LeaderboardEntry[] = [
  { id: 'lb-1', bootcampId: 'bootcamp-1', bootcampEnrollmentId: 'enrollment-mentee-2', problemsCompleted: 28, problemsAttempted: 30, completionRate: 93, streakDays: 12, score: 1240, rank: 1, calculatedAt: new Date().toISOString(), user: { id: 'user-mentee-2', name: 'Rahul Verma' } },
  { id: 'lb-2', bootcampId: 'bootcamp-1', bootcampEnrollmentId: 'enrollment-mentee-1', problemsCompleted: 22, problemsAttempted: 27, completionRate: 81, streakDays: 7, score: 980, rank: 2, calculatedAt: new Date().toISOString(), user: { id: 'user-mentee-1', name: 'Arjun Sharma' } },
  { id: 'lb-3', bootcampId: 'bootcamp-1', bootcampEnrollmentId: 'enrollment-mentee-3', problemsCompleted: 18, problemsAttempted: 22, completionRate: 81, streakDays: 4, score: 760, rank: 3, calculatedAt: new Date().toISOString(), user: { id: 'user-mentee-3', name: 'Sneha Patel' } },
];