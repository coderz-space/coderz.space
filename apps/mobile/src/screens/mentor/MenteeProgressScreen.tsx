import React, { useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ScrollView,
} from 'react-native';
import { useNavigation, useRoute, RouteProp } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withDelay,
  withSpring,
} from 'react-native-reanimated';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Badge from '../../components/atoms/Badge';
import { SkeletonCard } from '../../components/atoms/SkeletonLoader';
import StatBanner from '../../components/molecules/StatBanner';
import DoubtsCard from '../../components/molecules/DoubtsCard';
import ProgressRing from '../../components/molecules/ProgressRing';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMentorStore } from '../../store/mentorStore';
import { useAuthStore } from '../../store/authStore';
import { MentorStackParamList, Assignment } from '../../types';
import { formatDeadline, formatDate } from '../../utils/formatters';
import { toast } from '../../utils/toast';

type Route = RouteProp<MentorStackParamList, 'MenteeProgress'>;
type Nav = NativeStackNavigationProp<MentorStackParamList>;

export default function MenteeProgressScreen() {
  const navigation = useNavigation<Nav>();
  const route = useRoute<Route>();
  const { enrollmentId, menteeName } = route.params;
  const { session } = useAuthStore();
  const { menteeProgress, pendingDoubts, fetchMenteeProgress, resolveDoubt } =
    useMentorStore();

  useEffect(() => {
    if (!session) return;
    fetchMenteeProgress({
      orgId: session.activeOrgId,
      bootcampId: session.activeBootcampId,
      enrollmentId,
    });
  }, [enrollmentId]);

  const assignments = menteeProgress[enrollmentId] ?? [];
  const loading = assignments.length === 0;

  const totalProblems = assignments.reduce((s, a) => s + a.totalProblems, 0);
  const completedProblems = assignments.reduce(
    (s, a) => s + a.completedProblems,
    0,
  );
  const overallProgress =
    totalProblems > 0
      ? Math.round((completedProblems / totalProblems) * 100)
      : 0;

  const menteeDoubts = pendingDoubts.filter(d =>
    assignments.flatMap(a => a.problems).some(p => p.doubt?.id === d.id),
  );

  const handleResolve = async (doubtId: string) => {
    if (!session) return;
    try {
      await resolveDoubt({
        orgId: session.activeOrgId,
        bootcampId: session.activeBootcampId,
        doubtId,
      });
      toast.success('Doubt marked as resolved');
    } catch {
      toast.error('Failed to resolve doubt');
    }
  };

  return (
    <ScreenWrapper scrollable padded>
      {/* Back + Header */}
      <TouchableOpacity onPress={() => navigation.goBack()} style={styles.back}>
        <Text style={styles.backText}>← Back</Text>
      </TouchableOpacity>

      <View style={styles.header}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>
            {menteeName
              .split(' ')
              .map((n: string) => n[0])
              .slice(0, 2)
              .join('')
              .toUpperCase()}
          </Text>
        </View>
        <View>
          <Text style={styles.name}>{menteeName}</Text>
          <Text style={styles.role}>Mentee</Text>
        </View>
        <View style={styles.ringWrap}>
          <ProgressRing
            progress={overallProgress}
            size={64}
            strokeWidth={6}
            showPercentage={false}
            label={`${overallProgress}%`}
          />
        </View>
      </View>

      {/* Stats */}
      <StatBanner
        stats={[
          { label: 'Completed', value: completedProblems, accent: true },
          { label: 'Total', value: totalProblems },
          { label: 'Doubts', value: menteeDoubts.length },
          { label: 'Assignments', value: assignments.length },
        ]}
      />

      {/* Doubts section */}
      {menteeDoubts.length > 0 && (
        <>
          <Text style={styles.sectionTitle}>Open Doubts</Text>
          {menteeDoubts.map(d => (
            <DoubtsCard
              key={d.id}
              doubt={d}
              menteeName={menteeName}
              onResolve={handleResolve}
            />
          ))}
        </>
      )}

      {/* Assignments */}
      <Text style={styles.sectionTitle}>Assignments</Text>

      {loading ? (
        <>
          <SkeletonCard />
          <SkeletonCard />
        </>
      ) : assignments.length === 0 ? (
        <Text style={styles.empty}>No assignments yet.</Text>
      ) : (
        assignments.map((a, i) => (
          <AssignmentProgressCard key={a.id} assignment={a} index={i} />
        ))
      )}
    </ScreenWrapper>
  );
}

function AssignmentProgressCard({
  assignment,
  index,
}: {
  assignment: Assignment;
  index: number;
}) {
  const scale = useSharedValue(0.95);
  const opacity = useSharedValue(0);

  useEffect(() => {
    scale.value = withDelay(
      index * 80,
      withSpring(1, { damping: 14, stiffness: 180 }),
    );
    opacity.value = withDelay(index * 80, withSpring(1, { damping: 20 }));
  }, []);

  const animStyle = useAnimatedStyle(() => ({
    transform: [{ scale: scale.value }],
    opacity: opacity.value,
  }));

  const hasDoubts = assignment.problems.some(p => p.doubt && !p.doubt.resolved);

  return (
    <Animated.View style={[styles.assignCard, animStyle]}>
      <View style={styles.assignHeader}>
        <Text style={styles.assignTitle} numberOfLines={1}>
          {assignment.assignmentGroup.title}
        </Text>
        <Badge
          label={assignment.status}
          variant={
            assignment.status === 'completed'
              ? 'completed'
              : assignment.status === 'expired'
              ? 'doubt'
              : 'pending'
          }
        />
      </View>

      {/* Progress bar */}
      <View style={styles.progressBarWrap}>
        <View style={styles.progressBarTrack}>
          <View
            style={[
              styles.progressBarFill,
              { width: `${assignment.progressPercent}%` as any },
            ]}
          />
        </View>
        <Text style={styles.progressText}>
          {assignment.completedProblems}/{assignment.totalProblems}
        </Text>
      </View>

      {/* Problem list */}
      {assignment.problems.map(ap => (
        <View key={ap.id} style={styles.problemRow}>
          <Text style={styles.problemName} numberOfLines={1}>
            {ap.problem.title}
          </Text>
          <Badge
            label={
              ap.menteeStatus === 'completed'
                ? 'Done'
                : ap.menteeStatus === 'discussion_needed'
                ? 'Doubt'
                : ap.menteeStatus === 'revision_needed'
                ? 'Revision'
                : 'Pending'
            }
            variant={
              ap.menteeStatus === 'completed'
                ? 'completed'
                : ap.menteeStatus === 'discussion_needed'
                ? 'doubt'
                : ap.menteeStatus === 'revision_needed'
                ? 'in_progress'
                : 'pending'
            }
          />
        </View>
      ))}

      <Text style={styles.deadline}>
        {formatDeadline(assignment.deadlineAt)}
      </Text>
    </Animated.View>
  );
}

const styles = StyleSheet.create({
  back: { paddingTop: Spacing.lg, marginBottom: Spacing.md },
  backText: { ...Typography.label, color: Colors.primary },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: Spacing.xl,
    gap: Spacing.md,
  },
  avatar: {
    width: 52,
    height: 52,
    borderRadius: 26,
    backgroundColor: Colors.primaryMuted,
    borderWidth: 2,
    borderColor: Colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  avatarText: {
    ...Typography.headingSmall,
    color: Colors.primary,
    fontWeight: '800',
  },
  name: {
    ...Typography.headingMedium,
    color: Colors.textPrimary,
  },
  role: {
    ...Typography.caption,
    color: Colors.textSecondary,
  },
  ringWrap: {
    marginLeft: 'auto',
  },
  sectionTitle: {
    ...Typography.headingSmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.md,
    marginTop: Spacing.sm,
  },
  empty: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    textAlign: 'center',
    paddingVertical: Spacing['2xl'],
  },
  assignCard: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.md,
    ...Shadow.md,
  },
  assignHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: Spacing.md,
  },
  assignTitle: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    flex: 1,
    marginRight: Spacing.sm,
  },
  progressBarWrap: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: Spacing.sm,
    marginBottom: Spacing.md,
  },
  progressBarTrack: {
    flex: 1,
    height: 6,
    backgroundColor: Colors.surfaceElevated,
    borderRadius: 3,
    overflow: 'hidden',
  },
  progressBarFill: {
    height: '100%',
    backgroundColor: Colors.primary,
    borderRadius: 3,
  },
  progressText: {
    ...Typography.caption,
    color: Colors.textSecondary,
    width: 36,
    textAlign: 'right',
  },
  problemRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: Spacing.xs,
    borderTopWidth: 1,
    borderTopColor: Colors.divider,
  },
  problemName: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    flex: 1,
    marginRight: Spacing.sm,
  },
  deadline: {
    ...Typography.caption,
    color: Colors.textDisabled,
    marginTop: Spacing.sm,
  },
});
