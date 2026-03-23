import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { formatRelativeTime } from '../../utils/formatters';

export interface PollData {
  id: string;
  question: string;
  problemTitle?: string;
  createdAt: string;
  myVote?: 'easy' | 'medium' | 'hard';
  totalVotes?: number;
  results?: { easy: number; medium: number; hard: number };
}

interface Props {
  poll: PollData;
  onVote?: (pollId: string, vote: 'easy' | 'medium' | 'hard') => void;
  showResults?: boolean;
}

const VOTE_CONFIG = {
  easy: { color: Colors.success, label: 'Easy' },
  medium: { color: Colors.warning, label: 'Medium' },
  hard: { color: Colors.error, label: 'Hard' },
};

export default function PollCard({ poll, onVote, showResults = false }: Props) {
  const total = poll.results
    ? (poll.results.easy + poll.results.medium + poll.results.hard) || 1
    : 1;

  return (
    <View style={styles.card}>
      {poll.problemTitle && (
        <Text style={styles.problemLabel}>{poll.problemTitle}</Text>
      )}
      <Text style={styles.question}>{poll.question}</Text>

      <View style={styles.voteRow}>
        {(['easy', 'medium', 'hard'] as const).map((v) => {
          const cfg = VOTE_CONFIG[v];
          const isSelected = poll.myVote === v;
          const pct = showResults && poll.results
            ? Math.round((poll.results[v] / total) * 100)
            : null;

          return (
            <TouchableOpacity
              key={v}
              style={[
                styles.voteBtn,
                { borderColor: cfg.color },
                isSelected && { backgroundColor: `${cfg.color}20` },
              ]}
              onPress={() => onVote?.(poll.id, v)}
              disabled={!onVote || poll.myVote !== undefined}
            >
              {showResults && pct !== null && (
                <View
                  style={[
                    styles.voteBar,
                    { width: `${pct}%` as any, backgroundColor: `${cfg.color}15` },
                  ]}
                />
              )}
              <Text style={[styles.voteBtnText, { color: cfg.color }]}>
                {cfg.label}
                {showResults && pct !== null ? `  ${pct}%` : ''}
              </Text>
              {isSelected && <Text style={[styles.checkmark, { color: cfg.color }]}>✓</Text>}
            </TouchableOpacity>
          );
        })}
      </View>

      <View style={styles.footer}>
        <Text style={styles.time}>{formatRelativeTime(poll.createdAt)}</Text>
        {poll.totalVotes !== undefined && (
          <Text style={styles.votes}>{poll.totalVotes} votes</Text>
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.md,
    ...Shadow.md,
  },
  problemLabel: {
    ...Typography.labelSmall,
    color: Colors.primary,
    marginBottom: Spacing.xs,
  },
  question: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    marginBottom: Spacing.md,
  },
  voteRow: {
    flexDirection: 'row',
    gap: Spacing.sm,
    marginBottom: Spacing.md,
  },
  voteBtn: {
    flex: 1,
    borderWidth: 1.5,
    borderRadius: BorderRadius.lg,
    paddingVertical: Spacing.sm,
    alignItems: 'center',
    overflow: 'hidden',
    position: 'relative',
    minHeight: 40,
    justifyContent: 'center',
  },
  voteBar: {
    position: 'absolute',
    top: 0,
    left: 0,
    bottom: 0,
    borderRadius: BorderRadius.lg,
  },
  voteBtnText: {
    ...Typography.label,
    fontSize: 12,
    fontWeight: '700',
  },
  checkmark: {
    fontSize: 10,
    fontWeight: '800',
    position: 'absolute',
    top: 4,
    right: 6,
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  time: {
    ...Typography.caption,
    color: Colors.textDisabled,
  },
  votes: {
    ...Typography.caption,
    color: Colors.textSecondary,
  },
});