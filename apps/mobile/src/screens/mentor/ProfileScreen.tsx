import React from 'react';
import {
  View, Text, StyleSheet, TouchableOpacity, Alert,
} from 'react-native';
import Animated, {
  useSharedValue, useAnimatedStyle, withSpring,
} from 'react-native-reanimated';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Button from '../../components/atoms/Button';
import Badge from '../../components/atoms/Badge';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useAuthStore } from '../../store/authStore';
import { getInitials, formatDate } from '../../utils/formatters';
import { toast } from '../../utils/toast';
import ChangePasswordForm from '../../components/molecules/ChangePasswordForm';

export default function MentorProfileScreen() {
  const { session, logout } = useAuthStore();
  const scale = useSharedValue(1);

  if (!session) return null;

  const animStyle = useAnimatedStyle(() => ({
    transform: [{ scale: scale.value }],
  }));

  const handleLogout = () => {
    Alert.alert(
      'Sign Out',
      'Are you sure you want to sign out?',
      [
        { text: 'Cancel', style: 'cancel' },
        {
          text: 'Sign Out',
          style: 'destructive',
          onPress: async () => {
            await logout();
            toast.info('Signed out');
          },
        },
      ],
    );
  };

  const initials = getInitials(session.user.name);

  return (
    <ScreenWrapper scrollable padded>
      <View style={styles.header}>
        <Text style={styles.title}>Profile</Text>
      </View>

      {/* Avatar card */}
      <Animated.View style={[styles.avatarCard, animStyle]}>
        <View style={styles.avatar}>
          <Text style={styles.avatarText}>{initials}</Text>
        </View>
        <Text style={styles.name}>{session.user.name}</Text>
        <Text style={styles.email}>{session.user.email}</Text>
        <View style={styles.roleRow}>
          <Badge label={session.orgRole} variant="pending" />
          {session.bootcampRole && session.bootcampRole !== session.orgRole && (
            <Badge
              label={`${session.bootcampRole} in bootcamp`}
              variant="info"
              style={styles.ml8}
            />
          )}
        </View>
      </Animated.View>

      {/* Details */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Session Info</Text>
        <DetailRow label="Org ID" value={session.activeOrgId} />
        <DetailRow label="Bootcamp ID" value={session.activeBootcampId} />
        <DetailRow label="Member ID" value={session.orgMemberId} />
      </View>

      <ChangePasswordForm />

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

function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <View style={styles.detailRow}>
      <Text style={styles.detailLabel}>{label}</Text>
      <Text style={styles.detailValue} numberOfLines={1}>{value}</Text>
    </View>
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
    alignItems: 'center',
  },
  ml8: {
    marginLeft: Spacing.sm,
  },
  section: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.xl,
  },
  sectionTitle: {
    ...Typography.label,
    color: Colors.textSecondary,
    marginBottom: Spacing.md,
  },
  detailRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingVertical: Spacing.sm,
    borderBottomWidth: 1,
    borderBottomColor: Colors.divider,
  },
  detailLabel: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
  detailValue: {
    ...Typography.bodySmall,
    color: Colors.textPrimary,
    fontFamily: 'Courier New',
    flex: 1,
    textAlign: 'right',
    marginLeft: Spacing.sm,
  },
  logoutBtn: {
    marginBottom: Spacing['3xl'],
  },
});