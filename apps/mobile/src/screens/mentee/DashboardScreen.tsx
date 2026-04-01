import React, { useEffect } from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import AssignmentCard from '../../components/molecules/AssignmentCard';
import { SkeletonCard } from '../../components/atoms/SkeletonLoader';
import { Colors, Typography, Spacing, BorderRadius } from '../../theme';
import { useMenteeStore } from '../../store/menteeStore';
import { useAuthStore } from '../../store/authStore';
import { MenteeStackParamList } from '../../types';

type Nav = NativeStackNavigationProp<MenteeStackParamList>;

export default function MenteeDashboardScreen() {
  const navigation = useNavigation<Nav>();
  const { session } = useAuthStore();
  const { activeAssignments, isLoadingAssignments, fetchMyAssignments } = useMenteeStore();

  useEffect(() => {
    if (!session) return;
    fetchMyAssignments({
      orgId: session.activeOrgId,
      bootcampId: session.activeBootcampId,
      enrollmentId: session.bootcampEnrollmentId,
    });
  }, []);

  const firstName = session?.user.name.split(' ')[0] ?? '';
  const totalProblems = activeAssignments.reduce((s, a) => s + a.totalProblems, 0);
  const completedProblems = activeAssignments.reduce((s, a) => s + a.completedProblems, 0);
  const pendingDoubts = activeAssignments
    .flatMap((a) => a.problems)
    .filter((p) => p.doubt && !p.doubt.resolved).length;

  return (
    <ScreenWrapper scrollable padded>
      {/* Header */}
      <View style={styles.header}>
        <View>
          <Text style={styles.greeting}>Hey, {firstName} 👋</Text>
          <Text style={styles.subGreeting}>Keep pushing. You're doing great.</Text>
        </View>
        <TouchableOpacity
          style={styles.completedBtn}
          onPress={() => navigation.navigate('CompletedProblems')}
        >
          <Text style={styles.completedBtnText}>Completed</Text>
        </TouchableOpacity>
      </View>

      {/* Stats Row */}
      <View style={styles.statsRow}>
        <View style={styles.statCard}>
          <Text style={styles.statValue}>{completedProblems}</Text>
          <Text style={styles.statLabel}>Done</Text>
        </View>
        <View style={[styles.statCard, styles.statCardMiddle]}>
          <Text style={styles.statValue}>{totalProblems - completedProblems}</Text>
          <Text style={styles.statLabel}>Remaining</Text>
        </View>
        <View style={[styles.statCard, pendingDoubts > 0 && styles.statCardAlert]}>
          <Text style={[styles.statValue, pendingDoubts > 0 && styles.statValueAlert]}>
            {pendingDoubts}
          </Text>
          <Text style={[styles.statLabel, pendingDoubts > 0 && styles.statLabelAlert]}>
            Doubts
          </Text>
        </View>
      </View>

      {/* Assignments */}
      <Text style={styles.sectionTitle}>Active Assignments</Text>

      {isLoadingAssignments ? (
        <>
          <SkeletonCard />
          <SkeletonCard />
        </>
      ) : activeAssignments.length === 0 ? (
        <View style={styles.empty}>
          <Text style={styles.emptyText}>No active assignments yet.</Text>
          <Text style={styles.emptySubText}>Your mentor will assign tasks soon!</Text>
        </View>
      ) : (
        activeAssignments.map((assignment) => (
          <AssignmentCard
            key={assignment.id}
            assignment={assignment}
            onPress={() =>
              navigation.navigate('AssignmentDetail', { assignmentId: assignment.id })
            }
          />
        ))
      )}
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
  greeting: {
    ...Typography.displayMedium,
    color: Colors.textPrimary,
  },
  subGreeting: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginTop: 2,
  },
  completedBtn: {
    backgroundColor: Colors.surfaceElevated,
    paddingHorizontal: Spacing.md,
    paddingVertical: Spacing.sm,
    borderRadius: BorderRadius.lg,
    borderWidth: 1,
    borderColor: Colors.surfaceBorder,
  },
  completedBtnText: {
    ...Typography.label,
    color: Colors.textSecondary,
  },
  statsRow: {
    flexDirection: 'row',
    marginBottom: Spacing.xl,
    gap: Spacing.sm,
  },
  statCard: {
    flex: 1,
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    padding: Spacing.md,
    alignItems: 'center',
  },
  statCardMiddle: {
    borderWidth: 1,
    borderColor: Colors.surfaceBorder,
  },
  statCardAlert: {
    backgroundColor: Colors.errorMuted,
    borderWidth: 1,
    borderColor: Colors.error,
  },
  statValue: {
    ...Typography.headingLarge,
    color: Colors.primary,
    fontWeight: '800',
  },
  statValueAlert: {
    color: Colors.error,
  },
  statLabel: {
    ...Typography.caption,
    color: Colors.textSecondary,
    marginTop: 2,
  },
  statLabelAlert: {
    color: Colors.error,
  },
  sectionTitle: {
    ...Typography.headingSmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.md,
  },
  empty: {
    alignItems: 'center',
    paddingVertical: Spacing['3xl'],
  },
  emptyText: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  emptySubText: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
});