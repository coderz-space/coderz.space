import React from 'react';
import { View, Text, StyleSheet, Alert } from 'react-native';
import Animated, {
  useSharedValue, useAnimatedStyle, withDelay, withSpring,
} from 'react-native-reanimated';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Button from '../../components/atoms/Button';
import Badge from '../../components/atoms/Badge';
import StatBanner from '../../components/molecules/StatBanner';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useAuthStore } from '../../store/authStore';
import { useMenteeStore } from '../../store/menteeStore';
import { getInitials } from '../../utils/formatters';
import { toast } from '../../utils/toast';

export default function MenteeProfileScreen() {
  const { session, logout } = useAuthStore();
  const { activeAssignments } = useMenteeStore();

  const avatarScale = useSharedValue(0.8);
  const avatarOpacity = useSharedValue(0);

  React.useEffect(() => {
    avatarScale.value = withDelay(100, withSpring(1, { damping: 14, stiffness: 180 }));
    avatarOpacity.value = withDelay(100, withSpring(1, { damping: 20 }));
  }, []);

  const avatarAnim = useAnimatedStyle(() => ({
    transform: [{ scale: avatarScale.value }],
    opacity: avatarOpacity.value,
  }));

  if (!session) return null;

  const totalProblems = activeAssignments.reduce((s, a) => s + a.totalProblems, 0);
  const completedProblems = activeAssignments.reduce((s, a) => s + a.completedProblems, 0);
  const doubts = activeAssignments
    .flatMap((a) => a.problems)
    .filter((p) => p.doubt).length;

  const handleLogout = () => {
    Alert.alert('Sign Out', 'Are you sure?', [
      { text: 'Cancel', style: 'cancel' },
      {
        text: 'Sign Out',
        style: 'destructive',
        onPress: async () => {
          await logout();
          toast.info('Signed out');
        },
      },
    ]);
  };

  return (
    <ScreenWrapper scrollable padded>
      <View style={styles.header}>
        <Text style={styles.title}>Profile</Text>
      </View>

      {/* Avatar */}
      <Animated.View style={[styles.avatarCard, avatarAnim]}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>{getInitials(session.user.name)}</Text>
        </View>
        <Text style={styles.name}>{session.user.name}</Text>
        <Text style={styles.email}>{session.user.email}</Text>
        <View style={styles.roleRow}>
          <Badge label="Mentee" variant="info" />
        </View>
      </Animated.View>

      {/* Stats */}
      <StatBanner
        stats={[
          { label: 'Completed', value: completedProblems, accent: true },
          { label: 'Remaining', value: totalProblems - completedProblems },
          { label: 'Doubts', value: doubts },
        ]}
      />

      {/* Session info */}
      <View style={styles.section}>
        <Text style={styles.sectionLabel}>Session</Text>
        <View style={styles.infoRow}>
          <Text style={styles.infoKey}>Org ID</Text>
          <Text style={styles.infoVal}>{session.activeOrgId}</Text>
        </View>
        <View style={styles.infoRow}>
          <Text style={styles.infoKey}>Bootcamp</Text>
          <Text style={styles.infoVal}>{session.activeBootcampId}</Text>
        </View>
        <View style={styles.infoRow}>
          <Text style={styles.infoKey}>Enrollment</Text>
          <Text style={styles.infoVal}>{session.bootcampEnrollmentId}</Text>
        </View>
      </View>

      <Button
        label="Sign Out"
        onPress={handleLogout}
        variant="danger"
        fullWidth
        size="lg"
        style={styles.logoutBtn}
      />
    </ScreenWrapper>
  );
}

const styles = StyleSheet.create({
  header: {
    paddingTop: Spacing.xl,
    paddingBottom: Spacing.lg,
  },
  title: {
    ...Typography.displayMedium,
    color: Colors.textPrimary,
  },
  avatarCard: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius['2xl'],
    padding: Spacing.xl,
    alignItems: 'center',
    marginBottom: Spacing.xl,
    ...Shadow.lg,
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: Colors.primaryMuted,
    borderWidth: 3,
    borderColor: Colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: Spacing.md,
    shadowColor: Colors.primary,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 12,
    elevation: 6,
  },
  avatarText: {
    ...Typography.headingLarge,
    color: Colors.primary,
    fontWeight: '800',
  },
  name: {
    ...Typography.headingMedium,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  email: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.md,
  },
  roleRow: {
    flexDirection: 'row',
  },
  section: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.xl,
  },
  sectionLabel: {
    ...Typography.label,
    color: Colors.textSecondary,
    marginBottom: Spacing.md,
  },
  infoRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingVertical: Spacing.sm,
    borderBottomWidth: 1,
    borderBottomColor: Colors.divider,
  },
  infoKey: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
  infoVal: {
    ...Typography.bodySmall,
    color: Colors.textPrimary,
    fontFamily: 'Courier New',
    fontSize: 11,
    flex: 1,
    textAlign: 'right',
    marginLeft: Spacing.sm,
  },
  logoutBtn: {
    marginBottom: Spacing['3xl'],
  },
});