import React, { useEffect, useState } from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { useNavigation, useRoute, RouteProp } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import MenteeRow from '../../components/molecules/MenteeRow';
import { SkeletonMenteeRow } from '../../components/atoms/SkeletonLoader';
import { Colors, Typography, Spacing } from '../../theme';
import { useMentorStore } from '../../store/mentorStore';
import { useAuthStore } from '../../store/authStore';
import { MentorStackParamList } from '../../types';

type Nav = NativeStackNavigationProp<MentorStackParamList>;

export default function MenteeListScreen() {
  const navigation = useNavigation<Nav>();
  const { session } = useAuthStore();
  const { mentees, isLoadingMentees, fetchMentees } = useMentorStore();

  useEffect(() => {
    if (!session) return;
    fetchMentees({ orgId: session.activeOrgId, bootcampId: session.activeBootcampId });
  }, []);

  return (
    <ScreenWrapper scrollable padded>
      <TouchableOpacity onPress={() => navigation.goBack()} style={styles.back}>
        <Text style={styles.backText}>← Back</Text>
      </TouchableOpacity>
      <Text style={styles.title}>Select Mentee</Text>
      <Text style={styles.subtitle}>Choose a mentee to assign tasks to</Text>

      {isLoadingMentees ? (
        <><SkeletonMenteeRow /><SkeletonMenteeRow /></>
      ) : (
        mentees.map((m) => (
          <MenteeRow
            key={m.id}
            member={m}
            onPress={() =>
              navigation.navigate('AssignTask', { menteeEnrollmentId: m.id })
            }
          />
        ))
      )}
    </ScreenWrapper>
  );
}

const styles = StyleSheet.create({
  back: { paddingTop: Spacing.lg, marginBottom: Spacing.md },
  backText: { ...Typography.label, color: Colors.primary },
  title: { ...Typography.displayMedium, color: Colors.textPrimary, marginBottom: Spacing.xs },
  subtitle: { ...Typography.bodySmall, color: Colors.textSecondary, marginBottom: Spacing.xl },
});