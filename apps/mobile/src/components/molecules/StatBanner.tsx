import React from 'react';
import { View, Text, StyleSheet, ViewStyle } from 'react-native';
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withDelay,
  withSpring,
} from 'react-native-reanimated';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';

interface StatItem {
  label: string;
  value: string | number;
  accent?: boolean;
  subValue?: string;
}

interface Props {
  stats: StatItem[];
  style?: ViewStyle;
}

export default function StatBanner({ stats, style }: Props) {
  return (
    <View style={[styles.container, style]}>
      {stats.map((stat, i) => (
        <AnimatedStatItem key={stat.label} stat={stat} index={i} />
      ))}
    </View>
  );
}

function AnimatedStatItem({ stat, index }: { stat: StatItem; index: number }) {
  const scale = useSharedValue(0.85);
  const opacity = useSharedValue(0);

  React.useEffect(() => {
    scale.value = withDelay(index * 80, withSpring(1, { damping: 14, stiffness: 180 }));
    opacity.value = withDelay(index * 80, withSpring(1, { damping: 20 }));
  }, []);

  const animStyle = useAnimatedStyle(() => ({
    transform: [{ scale: scale.value }],
    opacity: opacity.value,
  }));

  return (
    <Animated.View
      style={[
        styles.statCard,
        stat.accent && styles.statCardAccent,
        animStyle,
      ]}
    >
      <Text style={[styles.value, stat.accent && styles.valueAccent]}>
        {stat.value}
      </Text>
      {stat.subValue && (
        <Text style={styles.subValue}>{stat.subValue}</Text>
      )}
      <Text style={styles.label}>{stat.label}</Text>
    </Animated.View>
  );
}

const styles = StyleSheet.create({
  container: {
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
  statCardAccent: {
    backgroundColor: Colors.primaryMuted,
    borderWidth: 1,
    borderColor: Colors.primary,
  },
  value: {
    ...Typography.headingLarge,
    color: Colors.textPrimary,
    fontWeight: '800',
  },
  valueAccent: {
    color: Colors.primary,
  },
  subValue: {
    ...Typography.caption,
    color: Colors.textDisabled,
    marginTop: -2,
  },
  label: {
    ...Typography.caption,
    color: Colors.textSecondary,
    marginTop: 2,
    textAlign: 'center',
  },
});