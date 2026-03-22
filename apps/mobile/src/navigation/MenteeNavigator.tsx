import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { MenteeStackParamList } from '../types';
import { Colors } from '../theme';

import MenteeDashboardScreen from '../screens/mentee/DashboardScreen';
import AssignmentDetailScreen from '../screens/mentee/AssignmentDetailScreen';
import ProblemDetailScreen from '../screens/mentee/ProblemDetailsScreen';
import CompletedScreen from '../screens/mentee/CompletedScreen';

const Stack = createNativeStackNavigator<MenteeStackParamList>();

export default function MenteeNavigator() {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: Colors.background },
        animation: 'slide_from_right',
      }}
    >
      <Stack.Screen name="Dashboard" component={MenteeDashboardScreen} />
      <Stack.Screen name="AssignmentDetail" component={AssignmentDetailScreen} />
      <Stack.Screen name="ProblemDetail" component={ProblemDetailScreen} />
      <Stack.Screen name="CompletedProblems" component={CompletedScreen} />
    </Stack.Navigator>
  );
}