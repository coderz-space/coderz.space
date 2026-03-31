import React from 'react';
import { View, Text, StyleSheet, ViewStyle } from 'react-native';
import { Colors, Typography, Spacing, BorderRadius } from '../../theme';
import { TaskStatus } from '../../types';

type BadgeVariant = TaskStatus | 'doubt' | 'info' | 'custom';

interface Props {
  label: string;
  variant?: BadgeVariant;
  size?: 'sm' | 'md';
  style?: ViewStyle;
  dot?: boolean;
}

const VARIANT_CONFIG: Record<
  BadgeVariant,
  { bg: string; text: string; border: string }
> = {
  pending: {
    bg: Colors.primaryMuted,
    text: Colors.primary,
    border: Colors.primary,
  },
  in_progress: {
    bg: Colors.warningMuted,
    text: Colors.warning,
    border: Colors.warning,
  },
  completed: {
    bg: Colors.successMuted,
    text: Colors.success,
    border: Colors.success,
  },
  review: {
    bg: `${Colors.info}20`,
    text: Colors.info,
    border: Colors.info,
  },
  doubt: {
    bg: Colors.errorMuted,
    text: Colors.error,
    border: Colors.error,
  },
  info: {
    bg: `${Colors.info}20`,
    text: Colors.info,
    border: Colors.info,
  },
  custom: {
    bg: Colors.surfaceElevated,
    text: Colors.textSecondary,
    border: Colors.surfaceBorder,
  },
};

export default function Badge({
  label,
  variant = 'custom',
  size = 'sm',
  style,
  dot = false,
}: Props) {
  const config = VARIANT_CONFIG[variant] ?? VARIANT_CONFIG.custom;

  return (
    <View
      style={[
        styles.base,
        size === 'sm' ? styles.sm : styles.md,
        {
          backgroundColor: config.bg,
          borderColor: config.border,
        },
        style,
      ]}
    >
      {dot && <View style={[styles.dot, { backgroundColor: config.text }]} />}
      <Text
        style={[
          size === 'sm' ? styles.labelSm : styles.labelMd,
          { color: config.text },
        ]}
      >
        {label.toUpperCase()}
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  base: {
    flexDirection: 'row',
    alignItems: 'center',
    alignSelf: 'flex-start',
    borderWidth: 1,
    borderRadius: BorderRadius.sm,
  },
  sm: {
    paddingHorizontal: Spacing.sm,
    paddingVertical: 3,
  },
  md: {
    paddingHorizontal: Spacing.md,
    paddingVertical: Spacing.xs,
  },
  dot: {
    width: 5,
    height: 5,
    borderRadius: 9999,
    marginRight: 5,
  },
  labelSm: {
    ...Typography.labelSmall,
    fontSize: 9,
  },
  labelMd: {
    ...Typography.labelSmall,
    fontSize: 11,
  },
});
