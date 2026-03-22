import React, { useEffect, useState } from 'react';
import {
  View, Text, StyleSheet, TouchableOpacity, Alert,
} from 'react-native';
import { useNavigation, useRoute, RouteProp } from '@react-navigation/native';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Button from '../../components/atoms/Button';
import Badge from '../../components/atoms/Badge';
import { SkeletonCard } from '../../components/atoms/SkeletonLoader';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMentorStore } from '../../store/mentorStore';
import { useAuthStore } from '../../store/authStore';
import { MentorStackParamList, AssignmentGroup } from '../../types';

type Route = RouteProp<MentorStackParamList, 'AssignTask'>;

export default function AssignTaskScreen() {
  const navigation = useNavigation();
  const route = useRoute<Route>();
  const { menteeEnrollmentId } = route.params ?? {};
  const { session } = useAuthStore();
  const {
    assignmentGroups, mentees, isLoadingGroups,
    fetchAssignmentGroups, assignToMentee,
  } = useMentorStore();

  const [selectedGroupId, setSelectedGroupId] = useState<string | null>(null);
  const [selectedMenteeId, setSelectedMenteeId] = useState<string | null>(
    menteeEnrollmentId ?? null,
  );
  const [isAssigning, setIsAssigning] = useState(false);

  useEffect(() => {
    if (!session) return;
    fetchAssignmentGroups({ orgId: session.activeOrgId, bootcampId: session.activeBootcampId });
  }, []);

  const selectedMentee = mentees.find((m) => m.id === selectedMenteeId);
  const selectedGroup = assignmentGroups.find((g) => g.id === selectedGroupId);

  const handleAssign = async () => {
    if (!selectedGroupId || !selectedMenteeId || !session) {
      Alert.alert('Incomplete', 'Select both a mentee and an assignment group.');
      return;
    }
    setIsAssigning(true);
    try {
      const deadlineAt = new Date(
        Date.now() + (selectedGroup?.deadlineDays ?? 7) * 24 * 60 * 60 * 1000,
      ).toISOString();

      await assignToMentee({
        orgId: session.activeOrgId,
        bootcampId: session.activeBootcampId,
        assignmentGroupId: selectedGroupId,
        bootcampEnrollmentId: selectedMenteeId,
        deadlineAt,
      });

      Alert.alert(
        'Assigned! ✅',
        `"${selectedGroup?.title}" assigned to ${selectedMentee?.user.name}`,
        [{ text: 'Done', onPress: () => navigation.goBack() }],
      );
    } catch (e: any) {
      Alert.alert('Error', e.message);
    } finally {
      setIsAssigning(false);
    }
  };

  return (
    <ScreenWrapper scrollable padded>
      <TouchableOpacity onPress={() => navigation.goBack()} style={styles.back}>
        <Text style={styles.backText}>← Back</Text>
      </TouchableOpacity>
      <Text style={styles.title}>Assign Task</Text>

      {/* Step 1: Select Mentee */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>1. Select Mentee</Text>
        {mentees.map((m) => (
          <TouchableOpacity
            key={m.id}
            style={[
              styles.selectionRow,
              selectedMenteeId === m.id && styles.selectionRowActive,
            ]}
            onPress={() => setSelectedMenteeId(m.id)}
          >
            <View style={styles.selectionAvatar}>
              <Text style={styles.selectionAvatarText}>
                {m.user.name.slice(0, 2).toUpperCase()}
              </Text>
            </View>
            <Text style={styles.selectionLabel}>{m.user.name}</Text>
            {selectedMenteeId === m.id && (
              <Text style={styles.checkmark}>✓</Text>
            )}
          </TouchableOpacity>
        ))}
      </View>

      {/* Step 2: Select Assignment Group (Master Tasklist from wireframe) */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>2. Select from Master Tasklist</Text>
        <Text style={styles.sectionHint}>These are reusable assignment groups in your org</Text>

        {isLoadingGroups ? (
          <><SkeletonCard /><SkeletonCard /></>
        ) : (
          assignmentGroups.map((g) => (
            <AssignmentGroupCard
              key={g.id}
              group={g}
              selected={selectedGroupId === g.id}
              onSelect={() => setSelectedGroupId(g.id)}
            />
          ))
        )}
      </View>

      {/* Summary */}
      {selectedGroup && selectedMentee && (
        <View style={styles.summary}>
          <Text style={styles.summaryTitle}>Assignment Preview</Text>
          <Text style={styles.summaryText}>
            Assigning <Text style={styles.summaryHighlight}>{selectedGroup.title}</Text>
            {' '}to{' '}
            <Text style={styles.summaryHighlight}>{selectedMentee.user.name}</Text>
          </Text>
          <Text style={styles.summaryDeadline}>
            Deadline: {selectedGroup.deadlineDays} days from today
          </Text>
        </View>
      )}

      <Button
        label="Assign Now"
        onPress={handleAssign}
        loading={isAssigning}
        disabled={!selectedGroupId || !selectedMenteeId}
        fullWidth
        size="lg"
        style={styles.assignBtn}
      />
    </ScreenWrapper>
  );
}

function AssignmentGroupCard({
  group,
  selected,
  onSelect,
}: {
  group: AssignmentGroup;
  selected: boolean;
  onSelect: () => void;
}) {
  return (
    <TouchableOpacity
      style={[styles.groupCard, selected && styles.groupCardActive]}
      onPress={onSelect}
      activeOpacity={0.85}
    >
      <View style={styles.groupHeader}>
        <Text style={styles.groupTitle}>{group.title}</Text>
        {selected && <Text style={styles.checkmark}>✓</Text>}
      </View>
      {group.description && (
        <Text style={styles.groupDesc} numberOfLines={2}>{group.description}</Text>
      )}
      <View style={styles.groupMeta}>
        <Badge label={`${group.problems?.length ?? 0} problems`} variant="info" />
        <Text style={styles.groupDeadline}>{group.deadlineDays}d deadline</Text>
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  back: { paddingTop: Spacing.lg, marginBottom: Spacing.md },
  backText: { ...Typography.label, color: Colors.primary },
  title: { ...Typography.displayMedium, color: Colors.textPrimary, marginBottom: Spacing.xl },
  section: { marginBottom: Spacing.xl },
  sectionTitle: {
    ...Typography.headingSmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.xs,
  },
  sectionHint: {
    ...Typography.caption,
    color: Colors.textDisabled,
    marginBottom: Spacing.md,
  },
  selectionRow: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    padding: Spacing.md,
    marginBottom: Spacing.sm,
    borderWidth: 1.5,
    borderColor: Colors.surfaceBorder,
  },
  selectionRowActive: {
    borderColor: Colors.primary,
    backgroundColor: Colors.primaryMuted,
  },
  selectionAvatar: {
    width: 36,
    height: 36,
    borderRadius: 18,
    backgroundColor: Colors.primaryMuted,
    borderWidth: 1,
    borderColor: Colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: Spacing.md,
  },
  selectionAvatarText: {
    ...Typography.label,
    color: Colors.primary,
    fontWeight: '700',
    fontSize: 12,
  },
  selectionLabel: {
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
    flex: 1,
  },
  checkmark: {
    color: Colors.primary,
    fontWeight: '700',
    fontSize: 16,
  },
  groupCard: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.sm,
    borderWidth: 1.5,
    borderColor: Colors.surfaceBorder,
    ...Shadow.sm,
  },
  groupCardActive: {
    borderColor: Colors.primary,
    backgroundColor: Colors.primaryMuted,
  },
  groupHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: Spacing.xs,
  },
  groupTitle: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    flex: 1,
    marginRight: Spacing.sm,
  },
  groupDesc: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.sm,
  },
  groupMeta: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: Spacing.sm,
  },
  groupDeadline: {
    ...Typography.caption,
    color: Colors.textSecondary,
  },
  summary: {
    backgroundColor: Colors.primaryMuted,
    borderRadius: BorderRadius.lg,
    padding: Spacing.base,
    borderWidth: 1,
    borderColor: Colors.primary,
    marginBottom: Spacing.lg,
  },
  summaryTitle: {
    ...Typography.label,
    color: Colors.primary,
    marginBottom: Spacing.xs,
  },
  summaryText: {
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
  },
  summaryHighlight: {
    color: Colors.primary,
    fontWeight: '700',
  },
  summaryDeadline: {
    ...Typography.caption,
    color: Colors.textSecondary,
    marginTop: Spacing.xs,
  },
  assignBtn: {
    marginBottom: Spacing['3xl'],
  },
});