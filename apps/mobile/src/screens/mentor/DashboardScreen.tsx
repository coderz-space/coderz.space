import React, { useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import MenteeRow from '../../components/molecules/MenteeRow';
import { SkeletonMenteeRow } from '../../components/atoms/SkeletonLoader';
import Badge from '../../components/atoms/Badge';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMentorStore } from '../../store/mentorStore';
import { useAuthStore } from '../../store/authStore';
import { MentorStackParamList } from '../../types';

type Nav = NativeStackNavigationProp<MentorStackParamList>;

export default function MentorDashboardScreen() {
  const navigation = useNavigation<Nav>();
  const { session } = useAuthStore();
  const {
    mentees, pendingDoubts,
    isLoadingMentees, isLoadingDoubts,
    fetchMentees, fetchPendingDoubts,
  } = useMentorStore();

  useEffect(() => {
    if (!session) return;
    fetchMentees({ orgId: session.activeOrgId, bootcampId: session.activeBootcampId });
    fetchPendingDoubts({ orgId: session.activeOrgId, bootcampId: session.activeBootcampId });
  }, []);

  const firstName = session?.user.name.split(' ')[0] ?? '';

  return (
    <ScreenWrapper scrollable padded>
      {/* Header */}
      <View style={styles.header}>
        <View>
          <Text style={styles.greeting}>Hello, {firstName} 👨‍💻</Text>
          <Text style={styles.subGreeting}>Mentor Dashboard</Text>
        </View>
        <TouchableOpacity
          style={styles.assignBtn}
          onPress={() => navigation.navigate('AssignTask', {})}
        >
          <Text style={styles.assignBtnText}>+ Assign</Text>
        </TouchableOpacity>
      </View>

      {/* Pending Doubts Alert Banner — from wireframe */}
      {pendingDoubts.length > 0 && (
        <View style={styles.alertBanner}>
          <Text style={styles.alertIcon}>💬</Text>
          <View style={styles.alertBody}>
            <Text style={styles.alertTitle}>
              {pendingDoubts.length} Doubt{pendingDoubts.length > 1 ? 's' : ''} Pending
            </Text>
            <Text style={styles.alertSub}>
              {pendingDoubts[0].message.slice(0, 60)}...
            </Text>
          </View>
          <Badge label={`${pendingDoubts.length}`} variant="doubt" />
        </View>
      )}

      {/* Stats Row */}
      <View style={styles.statsRow}>
        <View style={styles.statCard}>
          <Text style={styles.statValue}>{mentees.length}</Text>
          <Text style={styles.statLabel}>Mentees</Text>
        </View>
        <View style={[styles.statCard, styles.statMiddle]}>
          <Text style={styles.statValue}>{pendingDoubts.length}</Text>
          <Text style={styles.statLabel}>Doubts</Text>
        </View>
        <TouchableOpacity
          style={[styles.statCard, styles.statAction]}
          onPress={() => navigation.navigate('QuestionBank')}
        >
          <Text style={styles.statValueOrange}>QB</Text>
          <Text style={styles.statLabel}>Bank</Text>
        </TouchableOpacity>
      </View>

      {/* Mentee List — "Who are my mentees" from wireframe */}
      <View style={styles.sectionHeader}>
        <Text style={styles.sectionTitle}>My Mentees</Text>
        <Text style={styles.sectionSub}>{mentees.length} enrolled</Text>
      </View>

      {isLoadingMentees ? (
        <><SkeletonMenteeRow /><SkeletonMenteeRow /><SkeletonMenteeRow /></>
      ) : mentees.length === 0 ? (
        <Text style={styles.empty}>No mentees enrolled yet.</Text>
      ) : (
        mentees.map((m) => {
          const hasDoubt = pendingDoubts.some((d) =>
            d.raisedBy === m.id,
          );
          return (
            <MenteeRow
              key={m.id}
              member={m}
              completedCount={Math.floor(Math.random() * 8)}
              totalCount={10}
              hasDoubt={hasDoubt}
              onPress={() =>
                navigation.navigate('MenteeProgress', {
                  enrollmentId: m.id,
                  menteeName: m.user.name,
                })
              }
            />
          );
        })
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
  assignBtn: {
    backgroundColor: Colors.primary,
    paddingHorizontal: Spacing.base,
    paddingVertical: Spacing.sm,
    borderRadius: BorderRadius.lg,
    shadowColor: Colors.primary,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.4,
    shadowRadius: 8,
    elevation: 6,
  },
  assignBtnText: {
    ...Typography.label,
    color: Colors.textInverse,
    fontWeight: '700',
  },
  alertBanner: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: Colors.errorMuted,
    borderWidth: 1,
    borderColor: Colors.error,
    borderRadius: BorderRadius.lg,
    padding: Spacing.md,
    marginBottom: Spacing.base,
    gap: Spacing.sm,
  },
  alertIcon: {
    fontSize: 20,
  },
  alertBody: {
    flex: 1,
  },
  alertTitle: {
    ...Typography.label,
    color: Colors.error,
    fontWeight: '700',
  },
  alertSub: {
    ...Typography.caption,
    color: Colors.textSecondary,
    marginTop: 2,
  },
  statsRow: {
    flexDirection: 'row',
    gap: Spacing.sm,
    marginBottom: Spacing.xl,
  },
  statCard: {
    flex: 1,
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    padding: Spacing.md,
    alignItems: 'center',
    ...Shadow.sm,
  },
  statMiddle: {
    borderWidth: 1,
    borderColor: Colors.surfaceBorder,
  },
  statAction: {
    borderWidth: 1.5,
    borderColor: Colors.primary,
    backgroundColor: Colors.primaryMuted,
  },
  statValue: {
    ...Typography.headingLarge,
    color: Colors.textPrimary,
    fontWeight: '800',
  },
  statValueOrange: {
    ...Typography.headingLarge,
    color: Colors.primary,
    fontWeight: '800',
  },
  statLabel: {
    ...Typography.caption,
    color: Colors.textSecondary,
    marginTop: 2,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: Spacing.md,
  },
  sectionTitle: {
    ...Typography.headingSmall,
    color: Colors.textSecondary,
  },
  sectionSub: {
    ...Typography.caption,
    color: Colors.textDisabled,
  },
  empty: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    textAlign: 'center',
    paddingVertical: Spacing['2xl'],
  },
});