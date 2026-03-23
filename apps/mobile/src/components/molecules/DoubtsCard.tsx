import React from 'react';
import {
  View, Text, StyleSheet, TouchableOpacity,
} from 'react-native';
import { Doubt } from '../../types';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import Badge from '../atoms/Badge';
import { formatRelativeTime } from '../../utils/formatters';

interface Props {
  doubt: Doubt;
  menteeName?: string;
  problemTitle?: string;
  onResolve?: (doubtId: string) => void;
}

export default function DoubtsCard({ doubt, menteeName, problemTitle, onResolve }: Props) {
  return (
    <View style={[styles.card, doubt.resolved && styles.cardResolved]}>
      {/* Left accent */}
      <View style={[styles.accent, { backgroundColor: doubt.resolved ? Colors.success : Colors.error }]} />

      <View style={styles.body}>
        {/* Header */}
        <View style={styles.header}>
          <View style={styles.headerLeft}>
            {menteeName && (
              <Text style={styles.menteeName}>{menteeName}</Text>
            )}
            {problemTitle && (
              <Text style={styles.problemTitle} numberOfLines={1}>
                {problemTitle}
              </Text>
            )}
          </View>
          <View style={styles.headerRight}>
            <Badge
              label={doubt.resolved ? 'Resolved' : 'Pending'}
              variant={doubt.resolved ? 'completed' : 'doubt'}
              dot
            />
          </View>
        </View>

        {/* Message */}
        <Text style={styles.message}>{doubt.message}</Text>

        {/* Footer */}
        <View style={styles.footer}>
          <Text style={styles.time}>{formatRelativeTime(doubt.createdAt)}</Text>
          {!doubt.resolved && onResolve && (
            <TouchableOpacity
              style={styles.resolveBtn}
              onPress={() => onResolve(doubt.id)}
            >
              <Text style={styles.resolveBtnText}>Mark Resolved</Text>
            </TouchableOpacity>
          )}
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    flexDirection: 'row',
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    overflow: 'hidden',
    marginBottom: Spacing.md,
    ...Shadow.md,
  },
  cardResolved: {
    opacity: 0.65,
  },
  accent: {
    width: 4,
  },
  body: {
    flex: 1,
    padding: Spacing.base,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: Spacing.sm,
  },
  headerLeft: {
    flex: 1,
    marginRight: Spacing.sm,
  },
  headerRight: {},
  menteeName: {
    ...Typography.label,
    color: Colors.primary,
    fontWeight: '700',
    marginBottom: 2,
  },
  problemTitle: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
  message: {
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
    marginBottom: Spacing.md,
    lineHeight: 22,
  },
  footer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  time: {
    ...Typography.caption,
    color: Colors.textDisabled,
  },
  resolveBtn: {
    backgroundColor: Colors.successMuted,
    borderWidth: 1,
    borderColor: Colors.success,
    paddingHorizontal: Spacing.md,
    paddingVertical: Spacing.xs,
    borderRadius: BorderRadius.md,
  },
  resolveBtnText: {
    ...Typography.labelSmall,
    color: Colors.success,
    fontSize: 10,
  },
});