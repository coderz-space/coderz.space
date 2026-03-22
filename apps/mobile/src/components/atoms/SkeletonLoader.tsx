import React, { useEffect, useRef } from 'react';
import { Animated, StyleSheet, View, ViewStyle } from 'react-native';
import { Colors, BorderRadius } from '../../theme';

interface Props {
  width?: number | string;
  height?: number;
  borderRadius?: number;
  style?: ViewStyle;
}

export default function SkeletonLoader({
  width = '100%',
  height = 16,
  borderRadius = BorderRadius.md,
  style,
}: Props) {
  const shimmer = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    Animated.loop(
      Animated.sequence([
        Animated.timing(shimmer, { toValue: 1, duration: 900, useNativeDriver: false }),
        Animated.timing(shimmer, { toValue: 0, duration: 900, useNativeDriver: false }),
      ]),
    ).start();
  }, []);

  const backgroundColor = shimmer.interpolate({
    inputRange: [0, 1],
    outputRange: [Colors.surfaceElevated, Colors.surfaceBorder],
  });

  return (
    <Animated.View
      style={[
        { width: width as any, height, borderRadius, backgroundColor },
        style,
      ]}
    />
  );
}

// ── Preset skeletons ─────────────────────────────────────────────

export function SkeletonCard() {
  return (
    <View style={skStyles.card}>
      <SkeletonLoader height={14} width="60%" />
      <SkeletonLoader height={10} width="90%" style={skStyles.mt8} />
      <SkeletonLoader height={10} width="75%" style={skStyles.mt6} />
      <View style={skStyles.row}>
        <SkeletonLoader height={20} width={70} borderRadius={4} />
        <SkeletonLoader height={20} width={50} borderRadius={4} style={skStyles.ml8} />
      </View>
    </View>
  );
}

export function SkeletonMenteeRow() {
  return (
    <View style={skStyles.menteeRow}>
      <SkeletonLoader width={40} height={40} borderRadius={20} />
      <View style={skStyles.flex}>
        <SkeletonLoader height={13} width="50%" />
        <SkeletonLoader height={10} width="35%" style={skStyles.mt6} />
      </View>
      <SkeletonLoader height={10} width={50} />
    </View>
  );
}

const skStyles = StyleSheet.create({
  card: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: 16,
    marginBottom: 12,
    gap: 0,
  },
  menteeRow: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    padding: 14,
    marginBottom: 10,
    gap: 12,
  },
  row: {
    flexDirection: 'row',
    marginTop: 12,
  },
  flex: { flex: 1, gap: 0 },
  mt8: { marginTop: 8 },
  mt6: { marginTop: 6 },
  ml8: { marginLeft: 8 },
});