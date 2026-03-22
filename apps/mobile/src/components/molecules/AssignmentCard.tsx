import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { Assignment } from '../../types';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import ProgressRing from './ProgressRing';
import Badge from '../atoms/Badge';

interface Props {
  assignment: Assignment;
  onPress: () => void;
}

export default function AssignmentCard({ assignment, onPress }: Props) {
  const { assignmentGroup, progressPercent, completedProblems, totalProblems, deadlineAt, status } = assignment;

  const daysLeft = Math.ceil(
    (new Date(deadlineAt).getTime() - Date.now()) / (1000 * 60 * 60 * 24),
  );

  const hasDoubts = assignment.problems.some((p) => p.doubt && !p.doubt.resolved);

  return (
    <TouchableOpacity style={styles.card} onPress={onPress} activeOpacity={0.85}>
      <View style={styles.left}>
        <ProgressRing
          progress={progressPercent}
          size={64}
          strokeWidth={6}
          showPercentage={false}
          label={`${completedProblems}/${totalProblems}`}
        />
      </View>

      <View style={styles.body}>
        <Text style={styles.title} numberOfLines={2}>
          {assignmentGroup.title}
        </Text>
        <Text style={styles.subtitle}>
          {completedProblems} of {totalProblems} problems done
        </Text>

        <View style={styles.footer}>
          {status === 'active' && daysLeft > 0 ? (
            <Text style={[styles.deadline, daysLeft <= 2 && styles.deadlineUrgent]}>
              {daysLeft}d left
            </Text>
          ) : status === 'completed' ? (
            <Badge label="Completed" variant="completed" dot />
          ) : (
            <Badge label="Expired" variant="doubt" dot />
          )}
          {hasDoubts && (
            <Badge label="Doubt" variant="doubt" style={styles.doubtBadge} />
          )}
        </View>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    flexDirection: 'row',
    alignItems: 'center',
    padding: Spacing.base,
    marginBottom: Spacing.md,
    ...Shadow.md,
  },
  left: {
    marginRight: Spacing.base,
  },
  body: {
    flex: 1,
  },
  title: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    marginBottom: 4,
  },
  subtitle: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.sm,
  },
  footer: {
    flexDirection: 'row',
    alignItems: 'center',
    flexWrap: 'wrap',
    gap: Spacing.sm,
  },
  deadline: {
    ...Typography.label,
    color: Colors.primary,
  },
  deadlineUrgent: {
    color: Colors.error,
  },
  doubtBadge: {
    marginLeft: Spacing.sm,
  },
});