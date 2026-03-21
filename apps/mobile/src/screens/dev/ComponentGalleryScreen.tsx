import React, { useState } from 'react';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import { Colors, Typography, Spacing } from '../../theme';
import Button from '../../components/atoms/Button';
import Input from '../../components/atoms/Input';
import Badge from '../../components/atoms/Badge';
import TaskCard from '../../components/molecules/TaskCard';
import ProgressRing from '../../components/molecules/ProgressRing';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import { Task } from '../../types';

// Mock task for gallery
const MOCK_TASK: Task = {
  id: '1',
  title: 'Implement Binary Search Tree with AVL Rotations',
  description: 'Build a self-balancing BST and explain the rotation logic with comments.',
  status: 'in_progress',
  difficulty: 'hard',
  dueDate: new Date(Date.now() + 86400000 * 2).toISOString(),
  hasDoubt: true,
  doubtDescription: 'Confused about left-right rotation case',
  assignedAt: new Date().toISOString(),
  tags: ['dsa', 'trees', 'algorithms'],
  points: 50,
};

const MOCK_TASK_2: Task = {
  ...MOCK_TASK,
  id: '2',
  title: 'Build REST API with Go Fiber',
  status: 'completed',
  difficulty: 'medium',
  hasDoubt: false,
  points: 30,
};

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <View style={styles.section}>
      <Text style={styles.sectionTitle}>{title}</Text>
      <View style={styles.sectionDivider} />
      {children}
    </View>
  );
}

export default function ComponentGalleryScreen() {
  const [inputValue, setInputValue] = useState('');
  const [passwordValue, setPasswordValue] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  return (
    <ScreenWrapper scrollable>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.headerLabel}>CODERZ SPACE</Text>
        <Text style={styles.headerTitle}>Component Gallery</Text>
        <Text style={styles.headerSub}>Phase 2 — Design System</Text>
      </View>

      {/* BUTTONS */}
      <Section title="Buttons">
        <Button label="Primary Action" onPress={() => {}} variant="primary" fullWidth />
        <View style={styles.spacer} />
        <Button label="Secondary Action" onPress={() => {}} variant="secondary" fullWidth />
        <View style={styles.spacer} />
        <Button label="Outlined Action" onPress={() => {}} variant="outlined" fullWidth />
        <View style={styles.spacer} />
        <Button label="Ghost Action" onPress={() => {}} variant="ghost" fullWidth />
        <View style={styles.spacer} />
        <Button label="Danger Action" onPress={() => {}} variant="danger" fullWidth />
        <View style={styles.spacer} />
        <View style={styles.row}>
          <Button label="Small" onPress={() => {}} size="sm" style={styles.flex} />
          <View style={styles.rowGap} />
          <Button label="Medium" onPress={() => {}} size="md" style={styles.flex} />
          <View style={styles.rowGap} />
          <Button label="Large" onPress={() => {}} size="lg" style={styles.flex} />
        </View>
        <View style={styles.spacer} />
        <Button label="Loading State" onPress={() => {}} loading={true} fullWidth />
        <View style={styles.spacer} />
        <Button label="Disabled State" onPress={() => {}} disabled={true} fullWidth />
      </Section>

      {/* INPUTS */}
      <Section title="Inputs">
        <Input
          label="Email Address"
          placeholder="you@coderz.space"
          value={inputValue}
          onChangeText={setInputValue}
          keyboardType="email-address"
          autoCapitalize="none"
        />
        <Input
          label="Password"
          placeholder="Enter your password"
          value={passwordValue}
          onChangeText={setPasswordValue}
          secureTextEntry={!showPassword}
          hint="Must be at least 8 characters"
        />
        <Input
          label="With Error"
          placeholder="Enter something"
          value=""
          onChangeText={() => {}}
          error="This field is required"
        />
      </Section>

      {/* BADGES */}
      <Section title="Badges">
        <View style={styles.badgeRow}>
          <Badge label="Pending" variant="pending" dot style={styles.badgeSpacing} />
          <Badge label="In Progress" variant="in_progress" dot style={styles.badgeSpacing} />
          <Badge label="Completed" variant="completed" dot style={styles.badgeSpacing} />
          <Badge label="Review" variant="review" dot style={styles.badgeSpacing} />
          <Badge label="Doubt" variant="doubt" dot style={styles.badgeSpacing} />
        </View>
      </Section>

      {/* PROGRESS RINGS */}
      <Section title="Progress Rings">
        <View style={styles.ringsRow}>
          <ProgressRing
            progress={75}
            size={100}
            label="Weekly"
            sublabel="6/8 tasks"
          />
          <ProgressRing
            progress={45}
            size={100}
            label="Monthly"
            sublabel="18/40 tasks"
            color={Colors.info}
          />
          <ProgressRing
            progress={100}
            size={100}
            label="Done!"
            sublabel="All tasks"
            color={Colors.success}
          />
        </View>
      </Section>

      {/* TASK CARDS */}
      <Section title="Task Cards">
        <TaskCard task={MOCK_TASK} onPress={(t) => console.log('Task pressed:', t.id)} />
        <TaskCard task={MOCK_TASK_2} onPress={(t) => console.log('Task pressed:', t.id)} />
      </Section>

      <View style={styles.footer}>
        <Text style={styles.footerText}>🔥 Coderz Space — Phase 2 Complete</Text>
      </View>
    </ScreenWrapper>
  );
}

const styles = StyleSheet.create({
  header: {
    paddingVertical: Spacing['3xl'],
    alignItems: 'center',
  },
  headerLabel: {
    ...Typography.labelSmall,
    color: Colors.primary,
    marginBottom: Spacing.sm,
  },
  headerTitle: {
    ...Typography.displayMedium,
    color: Colors.textPrimary,
    marginBottom: Spacing.xs,
  },
  headerSub: {
    ...Typography.bodyMedium,
    color: Colors.textSecondary,
  },
  section: {
    marginBottom: Spacing['2xl'],
  },
  sectionTitle: {
    ...Typography.headingSmall,
    color: Colors.textSecondary,
    marginBottom: Spacing.sm,
  },
  sectionDivider: {
    height: 1,
    backgroundColor: Colors.divider,
    marginBottom: Spacing.base,
  },
  spacer: {
    height: Spacing.sm,
  },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  rowGap: {
    width: Spacing.sm,
  },
  flex: {
    flex: 1,
  },
  badgeRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: Spacing.sm,
  },
  badgeSpacing: {
    marginRight: Spacing.sm,
    marginBottom: Spacing.sm,
  },
  ringsRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    paddingVertical: Spacing.base,
  },
  footer: {
    alignItems: 'center',
    paddingVertical: Spacing['3xl'],
  },
  footerText: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
});