import React, { useEffect } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import Animated, {
  useSharedValue, useAnimatedStyle, withDelay, withSpring,
} from 'react-native-reanimated';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import { SkeletonMenteeRow } from '../../components/atoms/SkeletonLoader';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMenteeStore } from '../../store/menteeStore';
import { useAuthStore } from '../../store/authStore';
import { LeaderboardEntry } from '../../types';

const MEDAL = ['🥇', '🥈', '🥉'];

export default function LeaderboardScreen() {
  const { session } = useAuthStore();
  const { leaderboard, isLoadingLeaderboard, fetchLeaderboard } = useMenteeStore();

  useEffect(() => {
    if (!session) return;
    fetchLeaderboard({
      orgId: session.activeOrgId,
      bootcampId: session.activeBootcampId,
    });
  }, []);

  const myEnrollmentId = session?.bootcampEnrollmentId;

  return (
    <ScreenWrapper scrollable padded>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.title}>Leaderboard</Text>
        <Text style={styles.subtitle}>Bootcamp rankings</Text>
      </View>

      {/* Top 3 podium */}
      {!isLoadingLeaderboard && leaderboard.length >= 3 && (
        <View style={styles.podium}>
          {[leaderboard[1], leaderboard[0], leaderboard[2]].map((entry, i) => {
            const rank = i === 1 ? 1 : i === 0 ? 2 : 3;
            const isTop = rank === 1;
            return (
              <View
                key={entry.bootcampEnrollmentId}
                style={[styles.podiumItem, isTop && styles.podiumItemTop]}
              >
                <Text style={styles.medal}>{MEDAL[rank - 1]}</Text>
                <View style={[styles.podiumAvatar, isTop && styles.podiumAvatarTop]}>
                  <Text style={[styles.podiumAvatarText, isTop && styles.podiumAvatarTextTop]}>
                    {entry.user.name.slice(0, 2).toUpperCase()}
                  </Text>
                </View>
                <Text style={styles.podiumName} numberOfLines={1}>
                  {entry.user.name.split(' ')[0]}
                </Text>
                <Text style={styles.podiumScore}>{entry.score}pts</Text>
              </View>
            );
          })}
        </View>
      )}

      {/* Full list */}
      <Text style={styles.sectionTitle}>All Rankings</Text>

      {isLoadingLeaderboard ? (
        <><SkeletonMenteeRow /><SkeletonMenteeRow /><SkeletonMenteeRow /></>
      ) : (
        leaderboard.map((entry, i) => (
          <LeaderboardRow
            key={entry.bootcampEnrollmentId}
            entry={entry}
            index={i}
            isMe={entry.bootcampEnrollmentId === myEnrollmentId}
          />
        ))
      )}
    </ScreenWrapper>
  );
}

function LeaderboardRow({
  entry,
  index,
  isMe,
}: {
  entry: LeaderboardEntry;
  index: number;
  isMe: boolean;
}) {
  const scale = useSharedValue(0.96);
  const opacity = useSharedValue(0);

  useEffect(() => {
    scale.value = withDelay(index * 60, withSpring(1, { damping: 14, stiffness: 180 }));
    opacity.value = withDelay(index * 60, withSpring(1, { damping: 20 }));
  }, []);

  const animStyle = useAnimatedStyle(() => ({
    transform: [{ scale: scale.value }],
    opacity: opacity.value,
  }));

  return (
    <Animated.View
      style={[
        styles.row,
        isMe && styles.rowMe,
        animStyle,
      ]}
    >
      {/* Rank */}
      <View style={[styles.rankBox, entry.rank <= 3 && styles.rankBoxTop]}>
        <Text style={[styles.rank, entry.rank <= 3 && styles.rankTop]}>
          {entry.rank <= 3 ? MEDAL[entry.rank - 1] : `#${entry.rank}`}
        </Text>
      </View>

      {/* Avatar */}
      <View style={[styles.avatar, isMe && styles.avatarMe]}>
        <Text style={[styles.avatarText, isMe && styles.avatarTextMe]}>
          {entry.user.name.slice(0, 2).toUpperCase()}
        </Text>
      </View>

      {/* Info */}
      <View style={styles.info}>
        <Text style={[styles.name, isMe && styles.nameMe]}>
          {entry.user.name}{isMe ? ' (You)' : ''}
        </Text>
        <Text style={styles.meta}>
          {entry.problemsCompleted} done · {entry.streakDays}d streak
        </Text>
      </View>

      {/* Score */}
      <View style={styles.scoreBox}>
        <Text style={[styles.score, isMe && styles.scoreMe]}>{entry.score}</Text>
        <Text style={styles.scoreLabel}>pts</Text>
      </View>
    </Animated.View>
  );
}

const styles = StyleSheet.create({
  header: {
    paddingTop: Spacing.xl,
    paddingBottom: Spacing.md,
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
  podium: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'flex-end',
    marginBottom: Spacing.xl,
    paddingVertical: Spacing.base,
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    gap: Spacing.md,
    ...Shadow.md,
  },
  podiumItem: {
    alignItems: 'center',
    flex: 1,
  },
  podiumItemTop: {
    marginBottom: -Spacing.base,
  },
  medal: {
    fontSize: 20,
    marginBottom: Spacing.xs,
  },
  podiumAvatar: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: Colors.surfaceElevated,
    borderWidth: 2,
    borderColor: Colors.surfaceBorder,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: Spacing.xs,
  },
  podiumAvatarTop: {
    width: 56,
    height: 56,
    borderRadius: 28,
    borderColor: Colors.primary,
    backgroundColor: Colors.primaryMuted,
  },
  podiumAvatarText: {
    ...Typography.label,
    color: Colors.textSecondary,
    fontWeight: '700',
    fontSize: 12,
  },
  podiumAvatarTextTop: {
    color: Colors.primary,
    fontSize: 16,
  },
  podiumName: {
    ...Typography.caption,
    color: Colors.textSecondary,
    fontWeight: '600',
  },
  podiumScore: {
    ...Typography.caption,
    color: Colors.primary,
    fontWeight: '700',
  },
  sectionTitle: {
    ...Typography.headingSmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.md,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    padding: Spacing.md,
    marginBottom: Spacing.sm,
    gap: Spacing.md,
    ...Shadow.sm,
  },
  rowMe: {
    backgroundColor: Colors.primaryMuted,
    borderWidth: 1.5,
    borderColor: Colors.primary,
  },
  rankBox: {
    width: 36,
    alignItems: 'center',
  },
  rankBoxTop: {},
  rank: {
    ...Typography.label,
    color: Colors.textSecondary,
    fontSize: 12,
  },
  rankTop: {
    fontSize: 18,
  },
  avatar: {
    width: 36,
    height: 36,
    borderRadius: 18,
    backgroundColor: Colors.surfaceElevated,
    borderWidth: 1,
    borderColor: Colors.surfaceBorder,
    alignItems: 'center',
    justifyContent: 'center',
  },
  avatarMe: {
    borderColor: Colors.primary,
    backgroundColor: Colors.primaryMuted,
  },
  avatarText: {
    ...Typography.caption,
    color: Colors.textSecondary,
    fontWeight: '700',
    fontSize: 11,
  },
  avatarTextMe: {
    color: Colors.primary,
  },
  info: {
    flex: 1,
  },
  name: {
    ...Typography.label,
    color: Colors.textPrimary,
    fontWeight: '600',
  },
  nameMe: {
    color: Colors.primary,
    fontWeight: '800',
  },
  meta: {
    ...Typography.caption,
    color: Colors.textSecondary,
    marginTop: 2,
  },
  scoreBox: {
    alignItems: 'flex-end',
  },
  score: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    fontWeight: '800',
  },
  scoreMe: {
    color: Colors.primary,
  },
  scoreLabel: {
    ...Typography.caption,
    color: Colors.textSecondary,
  },
});