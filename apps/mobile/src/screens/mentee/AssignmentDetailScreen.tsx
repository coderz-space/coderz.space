import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, FlatList } from 'react-native';
import { useNavigation, useRoute, RouteProp } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Badge from '../../components/atoms/Badge';
import { SkeletonCard } from '../../components/atoms/SkeletonLoader';
import ProgressRing from '../../components/molecules/ProgressRing';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMenteeStore } from '../../store/menteeStore';
import { useAuthStore } from '../../store/authStore';
import { MenteeStackParamList, AssignmentProblem } from '../../types';
import { TouchableOpacity } from 'react-native';

type Route = RouteProp<MenteeStackParamList, 'AssignmentDetail'>;
type Nav = NativeStackNavigationProp<MenteeStackParamList>;

const DIFFICULTY_COLOR = {
  easy: Colors.success,
  medium: Colors.warning,
  hard: Colors.error,
};

const STATUS_LABEL: Record<string, string> = {
  not_started: 'Not Started',
  discussion_needed: 'Discussion Needed',
  revision_needed: 'Revision Needed',
  completed: 'Completed',
};

export default function AssignmentDetailScreen() {
  const navigation = useNavigation<Nav>();
  const route = useRoute<Route>();
  const { assignmentId } = route.params;
  const { session } = useAuthStore();
  const { activeAssignments, isLoadingAssignments } = useMenteeStore();

  const assignment = activeAssignments.find((a) => a.id === assignmentId);

  if (isLoadingAssignments) {
    return (
      <ScreenWrapper padded>
        <SkeletonCard /><SkeletonCard /><SkeletonCard />
      </ScreenWrapper>
    );
  }

  if (!assignment) {
    return (
      <ScreenWrapper padded>
        <Text style={{ color: Colors.textSecondary, marginTop: Spacing.xl }}>
          Assignment not found.
        </Text>
      </ScreenWrapper>
    );
  }

  const daysLeft = Math.ceil(
    (new Date(assignment.deadlineAt).getTime() - Date.now()) / (1000 * 60 * 60 * 24),
  );

  return (
    <ScreenWrapper scrollable padded>
      {/* Assignment Header */}
      <View style={styles.header}>
        <TouchableOpacity onPress={() => navigation.goBack()} style={styles.back}>
          <Text style={styles.backText}>← Back</Text>
        </TouchableOpacity>
        <Text style={styles.title}>{assignment.assignmentGroup.title}</Text>
        {assignment.assignmentGroup.description ? (
          <Text style={styles.description}>{assignment.assignmentGroup.description}</Text>
        ) : null}
        <View style={styles.metaRow}>
          <Text style={[styles.deadline, daysLeft <= 2 && styles.deadlineUrgent]}>
            {daysLeft > 0 ? `${daysLeft} days left` : 'Overdue'}
          </Text>
          <Text style={styles.progress}>
            {assignment.completedProblems}/{assignment.totalProblems} done
          </Text>
        </View>
      </View>

      {/* Progress Ring */}
      <View style={styles.ringSection}>
        <ProgressRing
          progress={assignment.progressPercent}
          size={100}
          strokeWidth={9}
          label="Progress"
          sublabel={`${assignment.progressPercent}%`}
          showPercentage={false}
        />
      </View>

      {/* Problem List - Left/Total style from wireframe */}
      <Text style={styles.sectionTitle}>
        Problems  •  {assignment.totalProblems - assignment.completedProblems} left
      </Text>

      {assignment.problems.map((ap) => (
        <ProblemRow
          key={ap.id}
          ap={ap}
          onPress={() =>
            navigation.navigate('ProblemDetail', {
              assignmentProblemId: ap.id,
              problemTitle: ap.problem.title,
              isCompleted: ap.status === 'completed',
            })
          }
        />
      ))}
    </ScreenWrapper>
  );
}

function ProblemRow({
  ap,
  onPress,
}: {
  ap: AssignmentProblem;
  onPress: () => void;
}) {
  const diffColor = DIFFICULTY_COLOR[ap.problem.difficulty];

  return (
    <TouchableOpacity style={styles.problemRow} onPress={onPress} activeOpacity={0.8}>
      <View style={styles.problemLeft}>
        <Text style={styles.problemTitle} numberOfLines={1}>
          {ap.problem.title}
        </Text>
        <View style={styles.diffRow}>
          <View style={[styles.diffDot, { backgroundColor: diffColor }]} />
          <Text style={[styles.diffText, { color: diffColor }]}>
            {ap.problem.difficulty}
          </Text>
          {ap.problem.tags.slice(0, 2).map((t) => (
            <Text key={t.id} style={styles.tag}>{t.name}</Text>
          ))}
        </View>
      </View>

      <View style={styles.problemRight}>
        <Badge
          label={STATUS_LABEL[ap.menteeStatus] ?? ap.menteeStatus}
          variant={
            ap.menteeStatus === 'completed' ? 'completed'
              : ap.menteeStatus === 'discussion_needed' ? 'doubt'
              : ap.menteeStatus === 'revision_needed' ? 'in_progress'
              : 'pending'
          }
        />
        <Text style={styles.chevron}>›</Text>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  header: {
    paddingTop: Spacing.lg,
    paddingBottom: Spacing.base,
  },
  back: {
    marginBottom: Spacing.base,
  },
  backText: {
    ...Typography.label,
    color: Colors.primary,
  },
  title: {
    ...Typography.displayMedium,
    color: Colors.textPrimary,
    marginBottom: Spacing.sm,
  },
  description: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.sm,
  },
  metaRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  deadline: {
    ...Typography.label,
    color: Colors.primary,
  },
  deadlineUrgent: {
    color: Colors.error,
  },
  progress: {
    ...Typography.label,
    color: Colors.textSecondary,
  },
  ringSection: {
    alignItems: 'center',
    paddingVertical: Spacing.xl,
  },
  sectionTitle: {
    ...Typography.headingSmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.md,
  },
  problemRow: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    padding: Spacing.base,
    marginBottom: Spacing.sm,
    ...Shadow.sm,
  },
  problemLeft: {
    flex: 1,
    marginRight: Spacing.sm,
  },
  problemTitle: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    marginBottom: 4,
  },
  diffRow: {
    flexDirection: 'row',
    alignItems: 'center',
    flexWrap: 'wrap',
    gap: 4,
  },
  diffDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
  },
  diffText: {
    ...Typography.caption,
    fontWeight: '600',
    textTransform: 'capitalize',
    marginRight: 4,
  },
  tag: {
    ...Typography.caption,
    color: Colors.textDisabled,
    backgroundColor: Colors.surfaceElevated,
    paddingHorizontal: 6,
    paddingVertical: 1,
    borderRadius: 4,
  },
  problemRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: Spacing.sm,
  },
  chevron: {
    ...Typography.headingMedium,
    color: Colors.textDisabled,
  },
});