import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import DoubtsCard from '../../components/molecules/DoubtsCard';
import { SkeletonCard } from '../../components/atoms/SkeletonLoader';
import { Colors, Typography, Spacing, BorderRadius } from '../../theme';
import { useMentorStore } from '../../store/mentorStore';
import { useAuthStore } from '../../store/authStore';
import { toast } from '../../utils/toast';

type FilterType = 'pending' | 'all';

export default function DoubtsScreen() {
  const { session } = useAuthStore();
  const { pendingDoubts, isLoadingDoubts, fetchPendingDoubts, resolveDoubt } = useMentorStore();
  const [filter, setFilter] = useState<FilterType>('pending');

  useEffect(() => {
    if (!session) return;
    fetchPendingDoubts({
      orgId: session.activeOrgId,
      bootcampId: session.activeBootcampId,
    });
  }, []);

  const handleResolve = async (doubtId: string) => {
    if (!session) return;
    try {
      await resolveDoubt({
        orgId: session.activeOrgId,
        bootcampId: session.activeBootcampId,
        doubtId,
      });
      toast.success('Doubt resolved ✓');
    } catch {
      toast.error('Failed to resolve doubt');
    }
  };

  const displayed = filter === 'pending'
    ? pendingDoubts.filter((d) => !d.resolved)
    : pendingDoubts;

  return (
    <ScreenWrapper scrollable padded>
      {/* Header */}
      <View style={styles.header}>
        <View>
          <Text style={styles.title}>Doubts</Text>
          <Text style={styles.subtitle}>
            {pendingDoubts.filter((d) => !d.resolved).length} unresolved
          </Text>
        </View>
      </View>

      {/* Filter tabs */}
      <View style={styles.filterRow}>
        {(['pending', 'all'] as FilterType[]).map((f) => (
          <TouchableOpacity
            key={f}
            style={[styles.filterTab, filter === f && styles.filterTabActive]}
            onPress={() => setFilter(f)}
          >
            <Text style={[styles.filterTabText, filter === f && styles.filterTabTextActive]}>
              {f === 'pending' ? 'Unresolved' : 'All'}
            </Text>
          </TouchableOpacity>
        ))}
      </View>

      {isLoadingDoubts ? (
        <><SkeletonCard /><SkeletonCard /></>
      ) : displayed.length === 0 ? (
        <View style={styles.empty}>
          <Text style={styles.emptyIcon}>🎉</Text>
          <Text style={styles.emptyTitle}>No pending doubts!</Text>
          <Text style={styles.emptyText}>All caught up. Your mentees are on track.</Text>
        </View>
      ) : (
        displayed.map((d) => (
          <DoubtsCard
            key={d.id}
            doubt={d}
            onResolve={!d.resolved ? handleResolve : undefined}
          />
        ))
      )}
    </ScreenWrapper>
  );
}

const styles = StyleSheet.create({
  header: {
    paddingTop: Spacing.xl,
    paddingBottom: Spacing.md,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
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
  filterRow: {
    flexDirection: 'row',
    gap: Spacing.sm,
    marginBottom: Spacing.xl,
  },
  filterTab: {
    paddingHorizontal: Spacing.base,
    paddingVertical: Spacing.sm,
    borderRadius: BorderRadius.lg,
    backgroundColor: Colors.surface,
    borderWidth: 1,
    borderColor: Colors.surfaceBorder,
  },
  filterTabActive: {
    backgroundColor: Colors.primaryMuted,
    borderColor: Colors.primary,
  },
  filterTabText: {
    ...Typography.label,
    color: Colors.textSecondary,
    fontSize: 13,
  },
  filterTabTextActive: {
    color: Colors.primary,
    fontWeight: '700',
  },
  empty: {
    alignItems: 'center',
    paddingVertical: Spacing['4xl'],
  },
  emptyIcon: {
    fontSize: 40,
    marginBottom: Spacing.md,
  },
  emptyTitle: {
    ...Typography.headingMedium,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  emptyText: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    textAlign: 'center',
  },
});