import React, { useEffect, useState } from 'react';
import {
  View, Text, StyleSheet, TouchableOpacity, TextInput, FlatList,
} from 'react-native';
import { useNavigation } from '@react-navigation/native';
import Animated, {
  useSharedValue, useAnimatedStyle, withDelay, withSpring,
} from 'react-native-reanimated';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Badge from '../../components/atoms/Badge';
import { SkeletonCard } from '../../components/atoms/SkeletonLoader';
import Button from '../../components/atoms/Button';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMentorStore } from '../../store/mentorStore';
import { useAuthStore } from '../../store/authStore';
import { Problem, Difficulty } from '../../types';
import { truncate } from '../../utils/formatters';

const DIFFICULTY_OPTIONS: Array<{ label: string; value: Difficulty | 'all' }> = [
  { label: 'All', value: 'all' },
  { label: 'Easy', value: 'easy' },
  { label: 'Medium', value: 'medium' },
  { label: 'Hard', value: 'hard' },
];

export default function QuestionBankScreen() {
  const navigation = useNavigation();
  const { session } = useAuthStore();
  const { problems, fetchProblems } = useMentorStore();

  const [search, setSearch] = useState('');
  const [difficulty, setDifficulty] = useState<Difficulty | 'all'>('all');
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadProblems();
  }, [difficulty]);

  const loadProblems = async () => {
    if (!session) return;
    setIsLoading(true);
    await fetchProblems({
      orgId: session.activeOrgId,
      search: search || undefined,
      difficulty: difficulty === 'all' ? undefined : difficulty,
    });
    setIsLoading(false);
  };

  const filtered = problems.filter((p) =>
    search
      ? p.title.toLowerCase().includes(search.toLowerCase())
      : true,
  );

  return (
    <View style={styles.root}>
      <View style={styles.topBar}>
        <TouchableOpacity onPress={() => navigation.goBack()} style={styles.back}>
          <Text style={styles.backText}>←</Text>
        </TouchableOpacity>
        <Text style={styles.title}>Question Bank</Text>
        <Text style={styles.count}>{problems.length} problems</Text>
      </View>

      {/* Search */}
      <View style={styles.searchWrap}>
        <TextInput
          style={styles.searchInput}
          placeholder="Search problems..."
          placeholderTextColor={Colors.textDisabled}
          value={search}
          onChangeText={setSearch}
          onSubmitEditing={loadProblems}
          returnKeyType="search"
        />
      </View>

      {/* Difficulty Filter */}
      <View style={styles.filterRow}>
        {DIFFICULTY_OPTIONS.map((opt) => (
          <TouchableOpacity
            key={opt.value}
            style={[
              styles.filterChip,
              difficulty === opt.value && styles.filterChipActive,
            ]}
            onPress={() => setDifficulty(opt.value)}
          >
            <Text
              style={[
                styles.filterChipText,
                difficulty === opt.value && styles.filterChipTextActive,
              ]}
            >
              {opt.label}
            </Text>
          </TouchableOpacity>
        ))}
      </View>

      {/* List */}
      {isLoading ? (
        <View style={styles.listContainer}>
          <SkeletonCard /><SkeletonCard /><SkeletonCard />
        </View>
      ) : (
        <FlatList
          data={filtered}
          keyExtractor={(p) => p.id}
          contentContainerStyle={styles.listContainer}
          showsVerticalScrollIndicator={false}
          renderItem={({ item, index }) => (
            <ProblemCard problem={item} index={index} />
          )}
          ListEmptyComponent={
            <Text style={styles.empty}>No problems found.</Text>
          }
        />
      )}
    </View>
  );
}

function ProblemCard({ problem, index }: { problem: Problem; index: number }) {
  const scale = useSharedValue(0.96);
  const opacity = useSharedValue(0);

  useEffect(() => {
    scale.value = withDelay(index * 50, withSpring(1, { damping: 14, stiffness: 180 }));
    opacity.value = withDelay(index * 50, withSpring(1, { damping: 20 }));
  }, []);

  const animStyle = useAnimatedStyle(() => ({
    transform: [{ scale: scale.value }],
    opacity: opacity.value,
  }));

  const diffVariant =
    problem.difficulty === 'easy' ? 'completed'
      : problem.difficulty === 'medium' ? 'in_progress'
      : 'doubt';

  return (
    <Animated.View style={[styles.problemCard, animStyle]}>
      <View style={styles.problemHeader}>
        <Text style={styles.problemTitle} numberOfLines={1}>{problem.title}</Text>
        <Badge label={problem.difficulty} variant={diffVariant} />
      </View>

      {problem.description ? (
        <Text style={styles.problemDesc} numberOfLines={2}>
          {truncate(problem.description, 80)}
        </Text>
      ) : null}

      <View style={styles.problemMeta}>
        <View style={styles.tagRow}>
          {problem.tags.slice(0, 3).map((t) => (
            <Text key={t.id} style={styles.tag}>{t.name}</Text>
          ))}
        </View>
        {problem.externalLink && (
          <Text style={styles.lcLink}>↗ LC</Text>
        )}
      </View>
    </Animated.View>
  );
}

const styles = StyleSheet.create({
  root: {
    flex: 1,
    backgroundColor: Colors.background,
  },
  topBar: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingTop: 56,
    paddingHorizontal: Spacing.base,
    paddingBottom: Spacing.md,
    backgroundColor: Colors.surface,
    borderBottomWidth: 1,
    borderBottomColor: Colors.surfaceBorder,
    gap: Spacing.md,
  },
  back: {
    padding: Spacing.xs,
  },
  backText: {
    ...Typography.headingMedium,
    color: Colors.primary,
  },
  title: {
    ...Typography.headingMedium,
    color: Colors.textPrimary,
    flex: 1,
  },
  count: {
    ...Typography.caption,
    color: Colors.textSecondary,
  },
  searchWrap: {
    paddingHorizontal: Spacing.base,
    paddingVertical: Spacing.md,
    backgroundColor: Colors.surface,
  },
  searchInput: {
    backgroundColor: Colors.surfaceElevated,
    borderRadius: BorderRadius.lg,
    paddingHorizontal: Spacing.base,
    paddingVertical: Spacing.md,
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
    borderWidth: 1,
    borderColor: Colors.surfaceBorder,
  },
  filterRow: {
    flexDirection: 'row',
    paddingHorizontal: Spacing.base,
    paddingBottom: Spacing.md,
    backgroundColor: Colors.surface,
    gap: Spacing.sm,
    borderBottomWidth: 1,
    borderBottomColor: Colors.divider,
  },
  filterChip: {
    paddingHorizontal: Spacing.md,
    paddingVertical: Spacing.xs,
    borderRadius: BorderRadius.full,
    borderWidth: 1,
    borderColor: Colors.surfaceBorder,
  },
  filterChipActive: {
    borderColor: Colors.primary,
    backgroundColor: Colors.primaryMuted,
  },
  filterChipText: {
    ...Typography.label,
    color: Colors.textSecondary,
    fontSize: 12,
  },
  filterChipTextActive: {
    color: Colors.primary,
  },
  listContainer: {
    padding: Spacing.base,
  },
  problemCard: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.md,
    ...Shadow.md,
  },
  problemHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: Spacing.sm,
  },
  problemTitle: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    flex: 1,
    marginRight: Spacing.sm,
  },
  problemDesc: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.sm,
  },
  problemMeta: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  tagRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 4,
    flex: 1,
  },
  tag: {
    ...Typography.caption,
    color: Colors.textDisabled,
    backgroundColor: Colors.surfaceElevated,
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 4,
  },
  lcLink: {
    ...Typography.label,
    color: Colors.primary,
    fontSize: 11,
  },
  empty: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    textAlign: 'center',
    paddingVertical: Spacing['3xl'],
  },
});