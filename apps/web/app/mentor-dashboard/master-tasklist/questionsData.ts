export interface SheetQuestion {
  id: string;
  title: string;
  topic: string;
  difficulty: "easy" | "medium" | "hard";
}

export const SHEET_QUESTIONS: Record<string, SheetQuestion[]> = {
  "gfg-dsa-360": [
    { id: "gfg-1", title: "Array Rotation", topic: "Arrays", difficulty: "easy" },
    { id: "gfg-2", title: "Kadane's Algorithm", topic: "Arrays", difficulty: "medium" },
    { id: "gfg-3", title: "Stock Buy and Sell", topic: "Arrays", difficulty: "easy" },
    { id: "gfg-4", title: "Trapping Rain Water", topic: "Arrays", difficulty: "hard" },
    { id: "gfg-5", title: "Reverse a Linked List", topic: "Linked List", difficulty: "easy" },
    { id: "gfg-6", title: "Detect Loop in Linked List", topic: "Linked List", difficulty: "medium" },
    { id: "gfg-7", title: "Merge Two Sorted Lists", topic: "Linked List", difficulty: "easy" },
    { id: "gfg-8", title: "Binary Search", topic: "Binary Search", difficulty: "easy" },
    { id: "gfg-9", title: "Search in Rotated Array", topic: "Binary Search", difficulty: "medium" },
    { id: "gfg-10", title: "Balanced Parentheses", topic: "Stack", difficulty: "easy" },
    { id: "gfg-11", title: "Next Greater Element", topic: "Stack", difficulty: "medium" },
    { id: "gfg-12", title: "Level Order Traversal", topic: "Trees", difficulty: "easy" },
    { id: "gfg-13", title: "Height of Binary Tree", topic: "Trees", difficulty: "easy" },
    { id: "gfg-14", title: "Lowest Common Ancestor", topic: "Trees", difficulty: "medium" },
    { id: "gfg-15", title: "Dijkstra's Algorithm", topic: "Graphs", difficulty: "hard" },
  ],
  "strivers-dsa-sheet": [
    { id: "stv-1", title: "Set Matrix Zeroes", topic: "Arrays", difficulty: "medium" },
    { id: "stv-2", title: "Pascal's Triangle", topic: "Arrays", difficulty: "easy" },
    { id: "stv-3", title: "Next Permutation", topic: "Arrays", difficulty: "medium" },
    { id: "stv-4", title: "Maximum Subarray", topic: "Arrays", difficulty: "medium" },
    { id: "stv-5", title: "Sort Colors", topic: "Arrays", difficulty: "medium" },
    { id: "stv-6", title: "Two Sum", topic: "Arrays", difficulty: "easy" },
    { id: "stv-7", title: "Reverse Linked List", topic: "Linked List", difficulty: "easy" },
    { id: "stv-8", title: "Middle of Linked List", topic: "Linked List", difficulty: "easy" },
    { id: "stv-9", title: "Merge Sort", topic: "Sorting", difficulty: "medium" },
    { id: "stv-10", title: "Quick Sort", topic: "Sorting", difficulty: "medium" },
    { id: "stv-11", title: "Implement Stack using Queue", topic: "Stack/Queue", difficulty: "easy" },
    { id: "stv-12", title: "Sliding Window Maximum", topic: "Sliding Window", difficulty: "hard" },
    { id: "stv-13", title: "Inorder Traversal", topic: "Trees", difficulty: "easy" },
    { id: "stv-14", title: "Diameter of Binary Tree", topic: "Trees", difficulty: "medium" },
    { id: "stv-15", title: "Number of Islands", topic: "Graphs", difficulty: "medium" },
  ],
};

export const SHEET_NAMES: Record<string, string> = {
  "gfg-dsa-360": "GFG DSA 360",
  "strivers-dsa-sheet": "Striver's DSA Sheet",
};
