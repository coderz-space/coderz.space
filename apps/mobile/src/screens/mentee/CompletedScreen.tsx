import React, { useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Badge from '../../components/atoms/Badge';
import { SkeletonCard } from '../../components/atoms/SkeletonLoader';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMenteeStore } from '../../store/menteeStore';
import { useAuthStore } from '../../store/authStore';
import { Assignment } from '../../types';

export default function CompletedScreen() {
  const navigation = useNavigation();
  const { session } = useAuthStore();
  const { completedAssignments, isLoadingCompleted, fetchCompletedAssignments } = useMenteeStore();

  useEffect(() => {
    if (!session) return;
    fetchCompletedAssignments({
      orgId: session.activeOrgId,
      bootcampId: session.activeBootcampId,
      enrollmentId: session.bootcampEnrollmentId,
    });
  }, []);

  return (
    <ScreenWrapper scrollable padded>
      <TouchableOpacity onPress={() => navigation.goBack()} style={styles.back}>
        <Text style={styles.backText}>← Back</Text>
      </TouchableOpacity>
      <Text style={styles.title}>Completed Assignments</Text>

      {isLoadingCompleted ? (
        <><SkeletonCard /><SkeletonCard /></>
      ) : completedAssignments.length === 0 ? (
        <Text style={styles.empty}>No completed assignments yet.</Text>
      ) : (
        completedAssignments.map((a) => <CompletedAssignmentCard key={a.id} assignment={a} />)
      )}
    </ScreenWrapper>
  );
}

function CompletedAssignmentCard({ assignment }: { assignment: Assignment }) {
  return (
    <View style={styles.card}>
      <View style={styles.cardHeader}>
        <Text style={styles.cardTitle}>{assignment.assignmentGroup.title}</Text>
        <Badge label="Completed" variant="completed" dot />
      </View>

      {/* Problem rows - matches bottom part of wireframe */}
      {assignment.problems.map((ap) => (
        <View key={ap.id} style={styles.problemRow}>
          <Text style={styles.problemName} numberOfLines={1}>
            {ap.problem.title}
          </Text>
          <Badge
            label={ap.menteeStatus === 'revision_needed' ? 'Revision' : 'Completed'}
            variant={ap.menteeStatus === 'revision_needed' ? 'in_progress' : 'completed'}
          />
        </View>
      ))}

      {/* Special sections from wireframe */}
      {assignment.problems.some((p) => p.menteeStatus === 'revision_needed') && (
        <View style={styles.alertSection}>
          <Text style={styles.alertLabel}>🔄 Revision Needed</Text>
        </View>
      )}
      {assignment.problems.some((p) => p.doubt && !p.doubt.resolved) && (
        <View style={[styles.alertSection, styles.alertDoubt]}>
          <Text style={[styles.alertLabel, { color: Colors.error }]}>
            💬 Discussion Needed
          </Text>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  back: { paddingTop: Spacing.lg, marginBottom: Spacing.md },
  backText: { ...Typography.label, color: Colors.primary },
  title: {
    ...Typography.displayMedium,
    color: Colors.textPrimary,
    marginBottom: Spacing.xl,
  },
  empty: {
    ...Typography.bodyMedium,
    color: Colors.textSecondary,
    textAlign: 'center',
    marginTop: Spacing['3xl'],
  },
  card: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.md,
    ...Shadow.md,
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: Spacing.md,
    paddingBottom: Spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: Colors.divider,
  },
  cardTitle: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    flex: 1,
    marginRight: Spacing.sm,
  },
  problemRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: Spacing.sm,
    borderBottomWidth: 1,
    borderBottomColor: Colors.divider,
  },
  problemName: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    flex: 1,
    marginRight: Spacing.sm,
  },
  alertSection: {
    marginTop: Spacing.sm,
    paddingTop: Spacing.sm,
    borderTopWidth: 1,
    borderTopColor: Colors.warningMuted,
  },
  alertDoubt: {
    borderTopColor: Colors.errorMuted,
  },
  alertLabel: {
    ...Typography.label,
    color: Colors.warning,
  },
});