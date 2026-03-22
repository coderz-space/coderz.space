import React, { useState } from 'react';
import {
  View, Text, StyleSheet, ScrollView, Alert, TouchableOpacity,
} from 'react-native';
import { useRoute, RouteProp, useNavigation } from '@react-navigation/native';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Button from '../../components/atoms/Button';
import Input from '../../components/atoms/Input';
import Badge from '../../components/atoms/Badge';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useMenteeStore } from '../../store/menteeStore';
import { useAuthStore } from '../../store/authStore';
import { MenteeStackParamList, MenteeStatus, AssignmentProblem } from '../../types';

type Route = RouteProp<MenteeStackParamList, 'ProblemDetail'>;

// Maps from wireframe status options to our MenteeStatus type
const STATUS_OPTIONS: { label: string; value: MenteeStatus; variant: any }[] = [
  { label: 'Not Started', value: 'not_started', variant: 'pending' },
  { label: 'Discussion Needed', value: 'discussion_needed', variant: 'doubt' },
  { label: 'Revision Needed', value: 'revision_needed', variant: 'in_progress' },
  { label: 'Completed', value: 'completed', variant: 'completed' },
];

export default function ProblemDetailScreen() {
  const route = useRoute<Route>();
  const navigation = useNavigation();
  const { assignmentProblemId, problemTitle, isCompleted } = route.params;
  const { session } = useAuthStore();
  const { activeAssignments, updateProblemProgress, raiseDoubt } = useMenteeStore();

  // Find the problem across all active assignments
  let foundAP: AssignmentProblem | undefined;
  let foundAssignmentId = '';
  for (const a of activeAssignments) {
    const p = a.problems.find((p) => p.id === assignmentProblemId);
    if (p) { foundAP = p; foundAssignmentId = a.id; break; }
  }

  const [status, setStatus] = useState<MenteeStatus>(foundAP?.menteeStatus ?? 'not_started');
  const [solutionLink, setSolutionLink] = useState(foundAP?.solutionLink ?? '');
  const [notes, setNotes] = useState(foundAP?.notes ?? '');
  const [remarkForSelf, setRemarkForSelf] = useState(foundAP?.remarkForSelf ?? '');
  const [remarkForMentor, setRemarkForMentor] = useState(foundAP?.remarkForMentor ?? '');
  const [doubtMessage, setDoubtMessage] = useState('');
  const [showDoubtInput, setShowDoubtInput] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [isRaisingDoubt, setIsRaisingDoubt] = useState(false);

  if (!foundAP || !session) {
    return (
      <ScreenWrapper padded>
        <Text style={{ color: Colors.textSecondary, marginTop: Spacing.xl }}>
          Problem not found.
        </Text>
      </ScreenWrapper>
    );
  }

  const handleSave = async () => {
    setIsSaving(true);
    try {
      await updateProblemProgress({
        orgId: session.activeOrgId,
        bootcampId: session.activeBootcampId,
        assignmentId: foundAssignmentId,
        assignmentProblemId,
        status,
        solutionLink,
        notes,
        remarkForSelf,
        remarkForMentor,
      });
      Alert.alert('Saved!', 'Your progress has been updated.', [
        { text: 'OK', onPress: () => navigation.goBack() },
      ]);
    } catch (e: any) {
      Alert.alert('Error', e.message);
    } finally {
      setIsSaving(false);
    }
  };

  const handleRaiseDoubt = async () => {
    if (!doubtMessage.trim()) {
      Alert.alert('Required', 'Please describe your doubt.');
      return;
    }
    setIsRaisingDoubt(true);
    try {
      await raiseDoubt({
        orgId: session.activeOrgId,
        bootcampId: session.activeBootcampId,
        assignmentId: foundAssignmentId,
        assignmentProblemId,
        message: doubtMessage,
      });
      setShowDoubtInput(false);
      setDoubtMessage('');
      Alert.alert('Doubt Raised', 'Your mentor will be notified.');
    } catch (e: any) {
      Alert.alert('Error', e.message);
    } finally {
      setIsRaisingDoubt(false);
    }
  };

  const problem = foundAP.problem;

  return (
    <ScreenWrapper scrollable padded avoidKeyboard>
      {/* Back */}
      <TouchableOpacity onPress={() => navigation.goBack()} style={styles.back}>
        <Text style={styles.backText}>← Back</Text>
      </TouchableOpacity>

      {/* Problem Info */}
      <View style={styles.problemInfo}>
        <View style={styles.titleRow}>
          <Text style={styles.title} numberOfLines={3}>{problem.title}</Text>
          {problem.externalLink && (
            <Badge label="LC" variant="info" style={styles.lcBadge} />
          )}
        </View>
        <View style={styles.tagRow}>
          <Badge label={problem.difficulty} variant={
            problem.difficulty === 'easy' ? 'completed'
              : problem.difficulty === 'medium' ? 'in_progress'
              : 'doubt'
          } />
          {problem.tags.map((t) => (
            <Text key={t.id} style={styles.tag}>{t.name}</Text>
          ))}
        </View>
        {problem.description ? (
          <Text style={styles.problemDesc} numberOfLines={3}>
            {problem.description}
          </Text>
        ) : null}
      </View>

      {/* Status Picker */}
      <View style={styles.section}>
        <Text style={styles.sectionTitle}>Status</Text>
        <View style={styles.statusGrid}>
          {STATUS_OPTIONS.map((opt) => (
            <TouchableOpacity
              key={opt.value}
              style={[
                styles.statusOption,
                status === opt.value && styles.statusOptionActive,
              ]}
              onPress={() => setStatus(opt.value)}
            >
              <View style={[
                styles.statusDot,
                { backgroundColor: status === opt.value ? Colors.primary : Colors.surfaceBorder }
              ]} />
              <Text style={[
                styles.statusLabel,
                status === opt.value && styles.statusLabelActive,
              ]}>
                {opt.label}
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      </View>

      {/* Attach Link */}
      <View style={styles.section}>
        <Input
          label="Attach Solution Link"
          placeholder="https://github.com/you/solution"
          value={solutionLink}
          onChangeText={setSolutionLink}
          autoCapitalize="none"
          keyboardType="url"
        />
      </View>

      {/* Write Note */}
      <View style={styles.section}>
        <Input
          label="Write Note"
          placeholder="Your approach, observations..."
          value={notes}
          onChangeText={setNotes}
          multiline
          numberOfLines={3}
          style={styles.textArea}
        />
      </View>

      {/* Remark for Self */}
      <View style={styles.section}>
        <Input
          label="Remark for Self"
          placeholder="What do you want to remember?"
          value={remarkForSelf}
          onChangeText={setRemarkForSelf}
          multiline
          numberOfLines={2}
          style={styles.textArea}
        />
      </View>

      {/* Remark for Mentor */}
      <View style={styles.section}>
        <Input
          label="Remark for Mentor"
          placeholder="Ask your mentor something specific..."
          value={remarkForMentor}
          onChangeText={setRemarkForMentor}
          multiline
          numberOfLines={2}
          style={styles.textArea}
        />
      </View>

      {/* Existing Doubt */}
      {foundAP.doubt && (
        <View style={[styles.section, styles.existingDoubt]}>
          <Text style={styles.existingDoubtLabel}>
            {foundAP.doubt.resolved ? '✅ Doubt Resolved' : '🔴 Doubt Raised'}
          </Text>
          <Text style={styles.existingDoubtMsg}>{foundAP.doubt.message}</Text>
        </View>
      )}

      {/* Raise Doubt */}
      {!foundAP.doubt && (
        <View style={styles.section}>
          {showDoubtInput ? (
            <>
              <Input
                label="Describe your doubt"
                placeholder="What are you confused about?"
                value={doubtMessage}
                onChangeText={setDoubtMessage}
                multiline
                numberOfLines={3}
                style={styles.textArea}
              />
              <View style={styles.doubtActions}>
                <Button
                  label="Cancel"
                  onPress={() => setShowDoubtInput(false)}
                  variant="ghost"
                  size="sm"
                  style={styles.flex}
                />
                <Button
                  label="Submit Doubt"
                  onPress={handleRaiseDoubt}
                  loading={isRaisingDoubt}
                  variant="outlined"
                  size="sm"
                  style={styles.flex}
                />
              </View>
            </>
          ) : (
            <Button
              label="Raise a Doubt"
              onPress={() => setShowDoubtInput(true)}
              variant="outlined"
              fullWidth
              size="sm"
            />
          )}
        </View>
      )}

      {/* Save */}
      <Button
        label="Save Progress"
        onPress={handleSave}
        loading={isSaving}
        fullWidth
        size="lg"
        style={styles.saveBtn}
      />
    </ScreenWrapper>
  );
}

const styles = StyleSheet.create({
  back: { paddingTop: Spacing.lg, marginBottom: Spacing.md },
  backText: { ...Typography.label, color: Colors.primary },
  problemInfo: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.lg,
    ...Shadow.sm,
  },
  titleRow: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    marginBottom: Spacing.sm,
  },
  title: {
    ...Typography.headingMedium,
    color: Colors.textPrimary,
    flex: 1,
    marginRight: Spacing.sm,
  },
  lcBadge: {},
  tagRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    alignItems: 'center',
    gap: Spacing.sm,
    marginBottom: Spacing.sm,
  },
  tag: {
    ...Typography.caption,
    color: Colors.textDisabled,
    backgroundColor: Colors.surfaceElevated,
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 4,
  },
  problemDesc: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    marginTop: Spacing.sm,
  },
  section: {
    marginBottom: Spacing.base,
  },
  sectionTitle: {
    ...Typography.label,
    color: Colors.textSecondary,
    marginBottom: Spacing.sm,
  },
  statusGrid: {
    gap: Spacing.sm,
  },
  statusOption: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.lg,
    padding: Spacing.md,
    borderWidth: 1.5,
    borderColor: Colors.surfaceBorder,
  },
  statusOptionActive: {
    borderColor: Colors.primary,
    backgroundColor: Colors.primaryMuted,
  },
  statusDot: {
    width: 10,
    height: 10,
    borderRadius: 5,
    marginRight: Spacing.sm,
  },
  statusLabel: {
    ...Typography.label,
    color: Colors.textSecondary,
  },
  statusLabelActive: {
    color: Colors.primary,
    fontWeight: '700',
  },
  textArea: {
    minHeight: 80,
    textAlignVertical: 'top',
  },
  existingDoubt: {
    backgroundColor: Colors.errorMuted,
    borderRadius: BorderRadius.lg,
    padding: Spacing.base,
    borderWidth: 1,
    borderColor: Colors.error,
  },
  existingDoubtLabel: {
    ...Typography.label,
    color: Colors.error,
    marginBottom: Spacing.xs,
  },
  existingDoubtMsg: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
  doubtActions: {
    flexDirection: 'row',
    gap: Spacing.sm,
    marginTop: Spacing.sm,
  },
  flex: { flex: 1 },
  saveBtn: {
    marginBottom: Spacing['3xl'],
  },
});