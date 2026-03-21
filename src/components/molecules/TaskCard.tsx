import React, { useRef } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Animated,
} from 'react-native';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { Task } from '../../types';
import Badge from '../atoms/Badge';

interface Props {
  task: Task;
  onPress?: (task: Task) => void;
  style?: object;
}

const DIFFICULTY_COLORS = {
  easy: Colors.success,
  medium: Colors.warning,
  hard: Colors.error,
};

export default function TaskCard({ task, onPress, style }: Props) {
  const scaleAnim = useRef(new Animated.Value(1)).current;

  const handlePressIn = () =>
    Animated.spring(scaleAnim, {
      toValue: 0.98,
      useNativeDriver: true,
      speed: 50,
    }).start();

  const handlePressOut = () =>
    Animated.spring(scaleAnim, {
      toValue: 1,
      useNativeDriver: true,
      speed: 50,
    }).start();

  const difficultyColor = DIFFICULTY_COLORS[task.difficulty];
  const isCompleted = task.status === 'completed';

  return (
    <Animated.View style={[{ transform: [{ scale: scaleAnim }] }, style]}>
      <TouchableOpacity
        onPress={() => onPress?.(task)}
        onPressIn={handlePressIn}
        onPressOut={handlePressOut}
        activeOpacity={1}
        style={[styles.card, isCompleted && styles.cardCompleted]}
      >
        {/* Left accent bar */}
        <View
          style={[
            styles.accentBar,
            { backgroundColor: isCompleted ? Colors.success : Colors.primary },
          ]}
        />

        <View style={styles.body}>
          {/* Top row: title + points */}
          <View style={styles.topRow}>
            <Text
              style={[
                styles.title,
                isCompleted && styles.titleCompleted,
              ]}
              numberOfLines={2}
            >
              {task.title}
            </Text>
            <View style={styles.pointsBadge}>
              <Text style={styles.pointsText}>{task.points}pts</Text>
            </View>
          </View>

          {/* Description */}
          <Text style={styles.description} numberOfLines={2}>
            {task.description}
          </Text>

          {/* Bottom row: badges + difficulty */}
          <View style={styles.bottomRow}>
            <Badge
              label={task.status.replace('_', ' ')}
              variant={task.status}
              dot
            />
            {task.hasDoubt && (
              <Badge label="Doubt" variant="doubt" style={styles.badgeSpacing} />
            )}
            <View style={styles.flex} />
            <View
              style={[styles.difficultyDot, { backgroundColor: difficultyColor }]}
            />
            <Text style={[styles.difficultyText, { color: difficultyColor }]}>
              {task.difficulty}
            </Text>
          </View>
        </View>
      </TouchableOpacity>
    </Animated.View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    flexDirection: 'row',
    overflow: 'hidden',
    ...Shadow.md,
    marginBottom: Spacing.md,
  },
  cardCompleted: {
    opacity: 0.75,
  },
  accentBar: {
    width: 4,
    borderTopLeftRadius: BorderRadius.xl,
    borderBottomLeftRadius: BorderRadius.xl,
  },
  body: {
    flex: 1,
    padding: Spacing.base,
  },
  topRow: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    marginBottom: Spacing.xs,
  },
  title: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    flex: 1,
    marginRight: Spacing.sm,
  },
  titleCompleted: {
    textDecorationLine: 'line-through',
    color: Colors.textSecondary,
  },
  pointsBadge: {
    backgroundColor: Colors.primaryMuted,
    paddingHorizontal: Spacing.sm,
    paddingVertical: 2,
    borderRadius: BorderRadius.sm,
  },
  pointsText: {
    ...Typography.labelSmall,
    color: Colors.primary,
    fontSize: 10,
  },
  description: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.md,
  },
  bottomRow: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  badgeSpacing: {
    marginLeft: Spacing.sm,
  },
  flex: {
    flex: 1,
  },
  difficultyDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
    marginRight: 5,
  },
  difficultyText: {
    ...Typography.caption,
    fontWeight: '600',
    textTransform: 'capitalize',
  },
});