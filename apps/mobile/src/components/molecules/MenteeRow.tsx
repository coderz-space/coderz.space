import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { OrgMember } from '../../types';
import { Colors, Typography, Spacing, BorderRadius } from '../../theme';

interface Props {
  member: OrgMember;
  completedCount?: number;
  totalCount?: number;
  hasDoubt?: boolean;
  onPress: () => void;
}

export default function MenteeRow({
  member,
  completedCount = 0,
  totalCount = 0,
  hasDoubt = false,
  onPress,
}: Props) {
  const initials = member.user.name
    .split(' ')
    .map((n) => n[0])
    .slice(0, 2)
    .join('')
    .toUpperCase();

  return (
    <TouchableOpacity style={styles.row} onPress={onPress} activeOpacity={0.8}>
      {/* Avatar */}
      <View style={styles.avatar}>
        <Text style={styles.avatarText}>{initials}</Text>
        {hasDoubt && <View style={styles.doubtDot} />}
      </View>

      {/* Info */}
      <View style={styles.info}>
        <Text style={styles.name}>{member.user.name}</Text>
        <Text style={styles.email}>{member.user.email}</Text>
      </View>

      {/* Progress */}
      <View style={styles.progress}>
        <Text style={styles.progressText}>
          {completedCount}
          <Text style={styles.progressTotal}>/{totalCount}</Text>
        </Text>
        <Text style={styles.progressLabel}>done</Text>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    padding: Spacing.base,
    marginBottom: Spacing.sm,
  },
  avatar: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: Colors.primaryMuted,
    borderWidth: 1.5,
    borderColor: Colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: Spacing.md,
    position: 'relative',
  },
  avatarText: {
    ...Typography.label,
    color: Colors.primary,
    fontWeight: '700',
  },
  doubtDot: {
    position: 'absolute',
    top: 0,
    right: 0,
    width: 10,
    height: 10,
    borderRadius: 5,
    backgroundColor: Colors.error,
    borderWidth: 1.5,
    borderColor: Colors.surface,
  },
  info: {
    flex: 1,
  },
  name: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    marginBottom: 2,
  },
  email: {
    ...Typography.caption,
    color: Colors.textSecondary,
  },
  progress: {
    alignItems: 'flex-end',
  },
  progressText: {
    ...Typography.headingSmall,
    color: Colors.primary,
    fontWeight: '700',
  },
  progressTotal: {
    color: Colors.textSecondary,
    fontWeight: '400',
  },
  progressLabel: {
    ...Typography.caption,
    color: Colors.textSecondary,
  },
});