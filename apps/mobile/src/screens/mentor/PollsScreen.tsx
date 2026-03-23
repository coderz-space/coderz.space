import React, { useState } from 'react';
import {
  View, Text, StyleSheet, TouchableOpacity, TextInput, Alert,
} from 'react-native';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import PollCard, { PollData } from '../../components/molecules/PollCard';
import Button from '../../components/atoms/Button';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { toast } from '../../utils/toast';

// Mock polls for now — replace with analyticsStore when backend ready
const MOCK_POLLS: PollData[] = [
  {
    id: 'poll-1',
    question: 'Was LRU Cache difficult?',
    problemTitle: 'LRU Cache',
    createdAt: new Date(Date.now() - 7200000).toISOString(),
    totalVotes: 12,
    results: { easy: 2, medium: 4, hard: 6 },
  },
  {
    id: 'poll-2',
    question: 'How did you find Binary Search?',
    problemTitle: 'Binary Search',
    createdAt: new Date(Date.now() - 3600000).toISOString(),
    totalVotes: 9,
    results: { easy: 7, medium: 2, hard: 0 },
  },
];

export default function PollsScreen() {
  const [polls, setPolls] = useState<PollData[]>(MOCK_POLLS);
  const [showCreate, setShowCreate] = useState(false);
  const [question, setQuestion] = useState('');
  const [problemTitle, setProblemTitle] = useState('');
  const [creating, setCreating] = useState(false);

  const handleCreate = async () => {
    if (!question.trim() || question.length < 10) {
      Alert.alert('Too short', 'Question must be at least 10 characters.');
      return;
    }
    setCreating(true);
    await new Promise((r) => setTimeout(r, 800));
    const newPoll: PollData = {
      id: `poll-${Date.now()}`,
      question,
      problemTitle: problemTitle || undefined,
      createdAt: new Date().toISOString(),
      totalVotes: 0,
      results: { easy: 0, medium: 0, hard: 0 },
    };
    setPolls((prev) => [newPoll, ...prev]);
    setQuestion('');
    setProblemTitle('');
    setShowCreate(false);
    setCreating(false);
    toast.success('Poll created!');
  };

  return (
    <ScreenWrapper scrollable padded>
      {/* Header */}
      <View style={styles.header}>
        <View>
          <Text style={styles.title}>Polls</Text>
          <Text style={styles.subtitle}>Track problem difficulty feedback</Text>
        </View>
        <TouchableOpacity
          style={styles.createBtn}
          onPress={() => setShowCreate((v) => !v)}
        >
          <Text style={styles.createBtnText}>{showCreate ? '✕' : '+ Poll'}</Text>
        </TouchableOpacity>
      </View>

      {/* Create Poll Form */}
      {showCreate && (
        <View style={styles.createForm}>
          <Text style={styles.formTitle}>New Poll</Text>
          <TextInput
            style={styles.input}
            placeholder="Problem title (optional)"
            placeholderTextColor={Colors.textDisabled}
            value={problemTitle}
            onChangeText={setProblemTitle}
          />
          <TextInput
            style={[styles.input, styles.questionInput]}
            placeholder="Your poll question... (min 10 chars)"
            placeholderTextColor={Colors.textDisabled}
            value={question}
            onChangeText={setQuestion}
            multiline
            numberOfLines={3}
          />
          <Button
            label="Create Poll"
            onPress={handleCreate}
            loading={creating}
            fullWidth
            size="md"
          />
        </View>
      )}

      {/* Polls list */}
      <Text style={styles.sectionTitle}>Active Polls</Text>
      {polls.map((poll) => (
        <PollCard
          key={poll.id}
          poll={poll}
          showResults={true}
        />
      ))}
    </ScreenWrapper>
  );
}

const styles = StyleSheet.create({
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    paddingTop: Spacing.xl,
    paddingBottom: Spacing.lg,
  },
  title: {
    ...Typography.displayMedium,
    color: Colors.textPrimary,
  },
  subtitle: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginTop: 2,
  },
  createBtn: {
    backgroundColor: Colors.primary,
    paddingHorizontal: Spacing.md,
    paddingVertical: Spacing.sm,
    borderRadius: BorderRadius.lg,
    shadowColor: Colors.primary,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.35,
    shadowRadius: 8,
    elevation: 5,
  },
  createBtnText: {
    ...Typography.label,
    color: Colors.textInverse,
    fontWeight: '700',
  },
  createForm: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.xl,
    borderWidth: 1,
    borderColor: Colors.primary,
    ...Shadow.md,
  },
  formTitle: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    marginBottom: Spacing.md,
  },
  input: {
    backgroundColor: Colors.surfaceElevated,
    borderRadius: BorderRadius.lg,
    paddingHorizontal: Spacing.base,
    paddingVertical: Spacing.md,
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
    borderWidth: 1,
    borderColor: Colors.surfaceBorder,
    marginBottom: Spacing.md,
  },
  questionInput: {
    minHeight: 80,
    textAlignVertical: 'top',
  },
  sectionTitle: {
    ...Typography.headingSmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.md,
  },
});