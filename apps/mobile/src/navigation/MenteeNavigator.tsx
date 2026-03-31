import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { MenteeStackParamList } from '../types';
import { Colors } from '../theme';

import MenteeTabNavigator from './MenteeTabNavigator';
import AssignmentDetailScreen from '../screens/mentee/AssignmentDetailScreen';
import ProblemDetailScreen from '../screens/mentee/ProblemDetailScreen';
import CompletedScreen from '../screens/mentee/CompletedScreen';

type MenteeRootStack = MenteeStackParamList & { Tabs: undefined };
const Stack = createNativeStackNavigator<MenteeRootStack>();

export default function MenteeNavigator() {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: Colors.background },
        animation: 'slide_from_right',
      }}
    >
      <Stack.Screen name="Tabs" component={MenteeTabNavigator} />
      <Stack.Screen name="AssignmentDetail" component={AssignmentDetailScreen} />
      <Stack.Screen name="ProblemDetail" component={ProblemDetailScreen} />
      <Stack.Screen name="CompletedProblems" component={CompletedScreen} />
    </Stack.Navigator>
  );
}