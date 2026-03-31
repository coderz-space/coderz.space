import React, { useEffect, useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  Alert,
  ScrollView,
} from 'react-native';
import { useNavigation, useRoute, RouteProp } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Button from '../../components/atoms/Button';
import Badge from '../../components/atoms/Badge';
import Dropdown from '../../components/atoms/DropDown';
import QuestionBankSelector from '../../components/molecules/QuestionBankSelector';
import { SkeletonCard } from '../../components/atoms/SkeletonLoader';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMentorStore } from '../../store/mentorStore';
import { useAuthStore } from '../../store/authStore';
import { MentorStackParamList, AssignmentGroup, Problem } from '../../types';
import { toast } from '../../utils/toast';

type Route = RouteProp<MentorStackParamList, 'AssignTask'>;
type Nav = NativeStackNavigationProp<MentorStackParamList>;

export default function AssignTaskScreen() {
  const navigation = useNavigation<Nav>();
  const route = useRoute<Route>();
  const { menteeEnrollmentId } = route.params ?? {};
  const { session } = useAuthStore();
  const {
    assignmentGroups,
    mentees,
    isLoadingGroups,
    fetchAssignmentGroups,
    assignToMentee,
    assignProblemsToMentee, // we need to add this to the store
  } = useMentorStore();

  const [selectedGroupId, setSelectedGroupId] = useState<string | null>(null);
  const [selectedMenteeId, setSelectedMenteeId] = useState<string | null>(
    menteeEnrollmentId ?? null
  );
  const [isAssigning, setIsAssigning] = useState(false);
  const [selectedProblems, setSelectedProblems] = useState<Problem[]>([]);
  const [showBankSelector, setShowBankSelector] = useState(false);

  useEffect(() => {
    if (!session) return;
    fetchAssignmentGroups({
      orgId: session.activeOrgId,
      bootcampId: session.activeBootcampId,
    });
  }, []);

  const selectedMentee = mentees.find((m) => m.id === selectedMenteeId);
  const selectedGroup = assignmentGroups.find((g) => g.id === selectedGroupId);

  const handleAssign = async () => {
    if (!session) return;
    if (!selectedMenteeId) {
      Alert.alert('Incomplete', 'Please select a mentee.');
      return;
    }
    if (selectedProblems.length === 0 && !selectedGroupId) {
      Alert.alert('Incomplete', 'Please select an assignment group or choose custom problems.');
      return;
    }

    setIsAssigning(true);
    try {
      const deadlineAt = new Date(
        Date.now() + (selectedGroup?.deadlineDays ?? 7) * 24 * 60 * 60 * 1000
      ).toISOString();

      if (selectedProblems.length > 0) {
        // Use custom problems assignment
        await assignProblemsToMentee({
          orgId: session.activeOrgId,
          bootcampId: session.activeBootcampId,
          bootcampEnrollmentId: selectedMenteeId,
          problemIds: selectedProblems.map(p => p.id),
          deadlineAt,
        });
        toast.success(`Assigned ${selectedProblems.length} problems to ${selectedMentee?.user.name}`);
      } else if (selectedGroupId) {
        // Use existing assignment group
        await assignToMentee({
          orgId: session.activeOrgId,
          bootcampId: session.activeBootcampId,
          assignmentGroupId: selectedGroupId,
          bootcampEnrollmentId: selectedMenteeId,
          deadlineAt,
        });
        toast.success(`Assigned "${selectedGroup?.title}" to ${selectedMentee?.user.name}`);
      }

      Alert.alert('Assigned!', 'Task has been assigned.', [
        { text: 'Done', onPress: () => navigation.goBack() },
      ]);
    } catch (e: any) {
      Alert.alert('Error', e.message);
    } finally {
      setIsAssigning(false);
    }
  };

  const handleSelectProblems = (problems: Problem[]) => {
    setSelectedProblems(problems);
    // Clear any previously selected group because we are using custom problems
    setSelectedGroupId(null);
  };

  // Convert mentees to dropdown items
  const menteeItems = mentees.map(m => ({ id: m.id, label: m.user.name }));

  // Convert assignment groups to dropdown items (for optional group selection)
  const groupItems = assignmentGroups.map(g => ({ id: g.id, label: g.title }));

  return (
    <ScreenWrapper scrollable padded>
      <TouchableOpacity onPress={() => navigation.goBack()} style={styles.back}>
        <Text style={styles.backText}>← Back</Text>
      </TouchableOpacity>
      <Text style={styles.title}>Assign Task</Text>

      {/* Step 1: Select Mentee */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>1. Select Mentee</Text>
        <Dropdown
          items={menteeItems}
          selectedId={selectedMenteeId}
          onSelect={setSelectedMenteeId}
          placeholder="Choose a mentee..."
        />
      </View>

      {/* Step 2: Select from Master Tasklist OR Custom Problems */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>2. Choose Tasks</Text>
        <Text style={styles.sectionHint}>
          You can either pick a pre‑made assignment group or select individual problems.
        </Text>

        {/* Option A: Existing assignment group */}
        <Text style={styles.subsectionLabel}>Option A: Assignment Group</Text>
        {isLoadingGroups ? (
          <SkeletonCard />
        ) : (
          <Dropdown
            items={groupItems}
            selectedId={selectedGroupId}
            onSelect={(id) => {
              setSelectedGroupId(id);
              setSelectedProblems([]); // clear custom problems
            }}
            placeholder="Select an assignment group (optional)"
          />
        )}

        {/* Option B: Custom problems from question bank */}
        <View style={styles.customSection}>
          <Text style={styles.subsectionLabel}>Option B: Custom Problems</Text>
          <Button
            label="Select from Question Bank"
            onPress={() => setShowBankSelector(true)}
            variant="outlined"
            size="md"
            fullWidth
          />
          {selectedProblems.length > 0 && (
            <View style={styles.selectedProblems}>
              <Text style={styles.selectedCount}>
                Selected: {selectedProblems.length} problem(s)
              </Text>
              {selectedProblems.slice(0, 3).map(p => (
                <Text key={p.id} style={styles.problemName} numberOfLines={1}>
                  • {p.title}
                </Text>
              ))}
              {selectedProblems.length > 3 && (
                <Text style={styles.problemName}>+{selectedProblems.length - 3} more</Text>
              )}
            </View>
          )}
        </View>
      </View>

      {/* Summary and Assign Button */}
      {selectedMentee && (selectedGroup || selectedProblems.length > 0) && (
        <View style={styles.summary}>
          <Text style={styles.summaryTitle}>Assignment Preview</Text>
          <Text style={styles.summaryText}>
            Assigning{' '}
            <Text style={styles.summaryHighlight}>
              {selectedProblems.length > 0
                ? `${selectedProblems.length} custom problem(s)`
                : selectedGroup?.title}
            </Text>{' '}
            to{' '}
            <Text style={styles.summaryHighlight}>{selectedMentee.user.name}</Text>
          </Text>
          <Text style={styles.summaryDeadline}>
            Deadline: {selectedGroup?.deadlineDays ?? 7} days from today
          </Text>
        </View>
      )}

      <Button
        label="Assign Now"
        onPress={handleAssign}
        loading={isAssigning}
        disabled={!selectedMenteeId || (!selectedGroupId && selectedProblems.length === 0)}
        fullWidth
        size="lg"
        style={styles.assignBtn}
      />

      <QuestionBankSelector
        visible={showBankSelector}
        onClose={() => setShowBankSelector(false)}
        onSelectProblems={handleSelectProblems}
      />
    </ScreenWrapper>
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
  subsectionLabel: {
    ...Typography.label,
    color: Colors.textSecondary,
    marginTop: Spacing.md,
    marginBottom: Spacing.sm,
  },
  customSection: {
    marginTop: Spacing.md,
  },
  selectedProblems: {
    marginTop: Spacing.md,
    backgroundColor: Colors.surfaceElevated,
    borderRadius: BorderRadius.md,
    padding: Spacing.sm,
  },
  selectedCount: {
    ...Typography.label,
    color: Colors.primary,
    marginBottom: Spacing.xs,
  },
  problemName: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginVertical: 2,
  },
  summary: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.xl,
    ...Shadow.md,
  },
  summaryTitle: {
    ...Typography.label,
    color: Colors.textSecondary,
    marginBottom: Spacing.sm,
  },
  summaryText: {
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  summaryHighlight: {
    fontWeight: 'bold',
    color: Colors.primary,
  },
  summaryDeadline: {
    ...Typography.caption,
    color: Colors.textDisabled,
  },
  assignBtn: { marginBottom: Spacing['2xl'] },
});