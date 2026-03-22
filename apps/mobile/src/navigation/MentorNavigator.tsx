import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { MentorStackParamList } from '../types';
import { Colors } from '../theme';

import MentorDashboardScreen from '../screens/mentor/DashboardScreen';
import MenteeListScreen from '../screens/mentor/MenteeListScreen';
import AssignTaskScreen from '../screens/mentor/AssignTaskScreen';
import PlaceholderScreen from '../screens/PlaceholderScreen';

const Stack = createNativeStackNavigator<MentorStackParamList>();

export default function MentorNavigator() {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: Colors.background },
        animation: 'slide_from_right',
      }}
    >
      <Stack.Screen name="Dashboard" component={MentorDashboardScreen} />
      <Stack.Screen name="MenteeList" component={MenteeListScreen} />
      <Stack.Screen name="AssignTask" component={AssignTaskScreen} />
      <Stack.Screen name="MenteeProgress" component={PlaceholderScreen} />
      <Stack.Screen name="QuestionBank" component={PlaceholderScreen} />
    </Stack.Navigator>
  );
}
