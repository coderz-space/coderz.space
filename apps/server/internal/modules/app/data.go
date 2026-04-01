package app

import "strings"

type sheetQuestion struct {
	ID          string
	Title       string
	Topic       string
	Difficulty  string
	Description string
}

type sheetCatalog struct {
	Key       string
	Name      string
	Questions []sheetQuestion
}

var catalogs = map[string]sheetCatalog{
	"gfg-dsa-360": {
		Key:  "gfg-dsa-360",
		Name: "GFG DSA 360",
		Questions: []sheetQuestion{
			{ID: "gfg-1", Title: "Array Rotation", Topic: "Arrays", Difficulty: "easy", Description: "Practice array rotation techniques and in-place updates."},
			{ID: "gfg-2", Title: "Kadane's Algorithm", Topic: "Arrays", Difficulty: "medium", Description: "Find the maximum subarray sum using dynamic running totals."},
			{ID: "gfg-3", Title: "Stock Buy and Sell", Topic: "Arrays", Difficulty: "easy", Description: "Track the best buy and sell window for maximum profit."},
			{ID: "gfg-4", Title: "Trapping Rain Water", Topic: "Arrays", Difficulty: "hard", Description: "Compute trapped water using prefix/suffix or two-pointer logic."},
			{ID: "gfg-5", Title: "Reverse a Linked List", Topic: "Linked List", Difficulty: "easy", Description: "Reverse a singly linked list iteratively or recursively."},
			{ID: "gfg-6", Title: "Detect Loop in Linked List", Topic: "Linked List", Difficulty: "medium", Description: "Use fast and slow pointers to detect a cycle."},
			{ID: "gfg-7", Title: "Merge Two Sorted Lists", Topic: "Linked List", Difficulty: "easy", Description: "Merge two sorted linked lists while preserving order."},
			{ID: "gfg-8", Title: "Binary Search", Topic: "Binary Search", Difficulty: "easy", Description: "Implement binary search on a sorted collection."},
			{ID: "gfg-9", Title: "Search in Rotated Array", Topic: "Binary Search", Difficulty: "medium", Description: "Find a target in a rotated sorted array."},
			{ID: "gfg-10", Title: "Balanced Parentheses", Topic: "Stack", Difficulty: "easy", Description: "Validate bracket matching using a stack."},
			{ID: "gfg-11", Title: "Next Greater Element", Topic: "Stack", Difficulty: "medium", Description: "Use a monotonic stack to find next greater values."},
			{ID: "gfg-12", Title: "Level Order Traversal", Topic: "Trees", Difficulty: "easy", Description: "Traverse a binary tree level by level using a queue."},
			{ID: "gfg-13", Title: "Height of Binary Tree", Topic: "Trees", Difficulty: "easy", Description: "Compute binary tree depth using DFS or BFS."},
			{ID: "gfg-14", Title: "Lowest Common Ancestor", Topic: "Trees", Difficulty: "medium", Description: "Find the lowest common ancestor of two nodes."},
			{ID: "gfg-15", Title: "Dijkstra's Algorithm", Topic: "Graphs", Difficulty: "hard", Description: "Compute shortest paths in a weighted graph."},
		},
	},
	"strivers-dsa-sheet": {
		Key:  "strivers-dsa-sheet",
		Name: "Striver's DSA Sheet",
		Questions: []sheetQuestion{
			{ID: "stv-1", Title: "Set Matrix Zeroes", Topic: "Arrays", Difficulty: "medium", Description: "Zero matrix rows and columns in-place with minimal extra space."},
			{ID: "stv-2", Title: "Pascal's Triangle", Topic: "Arrays", Difficulty: "easy", Description: "Generate rows of Pascal's triangle."},
			{ID: "stv-3", Title: "Next Permutation", Topic: "Arrays", Difficulty: "medium", Description: "Produce the next lexicographical permutation in-place."},
			{ID: "stv-4", Title: "Maximum Subarray", Topic: "Arrays", Difficulty: "medium", Description: "Find the maximum contiguous subarray sum."},
			{ID: "stv-5", Title: "Sort Colors", Topic: "Arrays", Difficulty: "medium", Description: "Sort three values using the Dutch national flag pattern."},
			{ID: "stv-6", Title: "Two Sum", Topic: "Arrays", Difficulty: "easy", Description: "Return indices of the two numbers that add to the target."},
			{ID: "stv-7", Title: "Reverse Linked List", Topic: "Linked List", Difficulty: "easy", Description: "Reverse a singly linked list."},
			{ID: "stv-8", Title: "Middle of Linked List", Topic: "Linked List", Difficulty: "easy", Description: "Find the middle node with fast and slow pointers."},
			{ID: "stv-9", Title: "Merge Sort", Topic: "Sorting", Difficulty: "medium", Description: "Implement divide-and-conquer merge sort."},
			{ID: "stv-10", Title: "Quick Sort", Topic: "Sorting", Difficulty: "medium", Description: "Partition and sort recursively using quick sort."},
			{ID: "stv-11", Title: "Implement Stack using Queue", Topic: "Stack/Queue", Difficulty: "easy", Description: "Simulate stack operations with queue primitives."},
			{ID: "stv-12", Title: "Sliding Window Maximum", Topic: "Sliding Window", Difficulty: "hard", Description: "Track maximum values inside a moving window."},
			{ID: "stv-13", Title: "Inorder Traversal", Topic: "Trees", Difficulty: "easy", Description: "Traverse a binary tree in inorder sequence."},
			{ID: "stv-14", Title: "Diameter of Binary Tree", Topic: "Trees", Difficulty: "medium", Description: "Compute the longest path through a binary tree."},
			{ID: "stv-15", Title: "Number of Islands", Topic: "Graphs", Difficulty: "medium", Description: "Count connected land components in a grid."},
		},
	},
}

var orderedSheetKeys = []string{
	"gfg-dsa-360",
	"strivers-dsa-sheet",
}

func listSheets() []SheetData {
	sheets := make([]SheetData, 0, len(orderedSheetKeys))
	for _, key := range orderedSheetKeys {
		sheets = append(sheets, sheetToData(catalogs[key]))
	}
	return sheets
}

func sheetToData(catalog sheetCatalog) SheetData {
	questions := make([]SheetQuestionData, 0, len(catalog.Questions))
	for _, question := range catalog.Questions {
		questions = append(questions, SheetQuestionData{
			ID:         question.ID,
			Title:      question.Title,
			Topic:      question.Topic,
			Difficulty: question.Difficulty,
		})
	}

	return SheetData{
		Key:       catalog.Key,
		Name:      catalog.Name,
		Questions: questions,
	}
}

func findSheet(key string) (sheetCatalog, bool) {
	catalog, ok := catalogs[key]
	return catalog, ok
}

func findSheetQuestion(sheetKey, questionID string) (sheetQuestion, bool) {
	catalog, ok := catalogs[sheetKey]
	if !ok {
		return sheetQuestion{}, false
	}

	for _, question := range catalog.Questions {
		if question.ID == questionID {
			return question, true
		}
	}

	return sheetQuestion{}, false
}

func catalogLink(sheetKey, questionID string) string {
	return "app-sheet:" + sheetKey + ":" + questionID
}

func findSheetQuestionByLink(link string) (sheetQuestion, bool) {
	if !strings.HasPrefix(link, "app-sheet:") {
		return sheetQuestion{}, false
	}

	parts := strings.Split(link, ":")
	if len(parts) != 3 {
		return sheetQuestion{}, false
	}

	return findSheetQuestion(parts[1], parts[2])
}
