import React, { useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, FlatList } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withDelay,
  withSpring,
} from 'react-native-reanimated';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import { SkeletonMenteeRow } from '../../components/atoms/SkeletonLoader';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMenteeStore } from '../../store/menteeStore';
import { useAuthStore } from '../../store/authStore';
import { LeaderboardEntry, LeaderboardPeriod } from '../../types';

const MEDAL = ['🥇', '🥈', '🥉'];

export default function LeaderboardScreen() {
  const navigation = useNavigation();
  const { session } = useAuthStore();
  const {
    leaderboard,
    isLoadingLeaderboard,
    leaderboardPeriod,
    setLeaderboardPeriod,
    fetchLeaderboard,
  } = useMenteeStore();

  useEffect(() => {
    if (!session) return;
    fetchLeaderboard({
      orgId: session.activeOrgId,
      bootcampId: session.activeBootcampId,
      period: leaderboardPeriod,
    });
  }, [leaderboardPeriod, session]);

  const myEnrollmentId = session?.bootcampEnrollmentId;

  const periodOptions: { label: string; value: LeaderboardPeriod }[] = [
    { label: 'Weekly', value: 'weekly' },
    { label: 'Monthly', value: 'monthly' },
    { label: 'All Time', value: 'allTime' },
  ];

  return (
    <ScreenWrapper scrollable padded>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.title}>Leaderboard</Text>
        <Text style={styles.subtitle}>Bootcamp rankings</Text>
      </View>

      {/* Period Filter */}
      <View style={styles.filterContainer}>
        {periodOptions.map((opt) => (
          <TouchableOpacity
            key={opt.value}
            style={[
              styles.filterTab,
              leaderboardPeriod === opt.value && styles.filterTabActive,
            ]}
            onPress={() => setLeaderboardPeriod(opt.value)}
          >
            <Text
              style={[
                styles.filterText,
                leaderboardPeriod === opt.value && styles.filterTextActive,
              ]}
            >
              {opt.label}
            </Text>
          </TouchableOpacity>
        ))}
      </View>

      {/* Top 3 Podium */}
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
                  <Text
                    style={[
                      styles.podiumAvatarText,
                      isTop && styles.podiumAvatarTextTop,
                    ]}
                  >
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

      {/* Full Rankings List */}
      <Text style={styles.sectionTitle}>All Rankings</Text>
      {isLoadingLeaderboard ? (
        <>
          <SkeletonMenteeRow />
          <SkeletonMenteeRow />
          <SkeletonMenteeRow />
        </>
      ) : (
        <FlatList
          data={leaderboard}
          keyExtractor={(item) => item.bootcampEnrollmentId}
          renderItem={({ item, index }) => (
            <LeaderboardRow
              entry={item}
              index={index}
              isMe={item.bootcampEnrollmentId === myEnrollmentId}
            />
          )}
          scrollEnabled={false}
        />
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
    <Animated.View style={[styles.row, isMe && styles.rowMe, animStyle]}>
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
          {entry.user.name}
          {isMe ? ' (You)' : ''}
        </Text>
        <Text style={styles.meta}>
          {entry.problemsCompleted} done • {entry.streakDays} d streak
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
  filterContainer: {
    flexDirection: 'row',
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    padding: 4,
    marginBottom: Spacing.xl,
  },
  filterTab: {
    flex: 1,
    paddingVertical: Spacing.sm,
    alignItems: 'center',
    borderRadius: BorderRadius.md,
  },
  filterTabActive: {
    backgroundColor: Colors.primary,
  },
  filterText: {
    ...Typography.label,
    color: Colors.textSecondary,
  },
  filterTextActive: {
    color: Colors.textInverse,
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