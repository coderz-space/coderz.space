import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { MentorStackParamList } from '../types';
import { Colors } from '../theme';

import MentorTabNavigator from './MentorTabNavigator';
import MenteeListScreen from '../screens/mentor/MenteeListScreen';
import AssignTaskScreen from '../screens/mentor/AssignTaskScreen';
import MenteeProgressScreen from '../screens/mentor/MenteeProgressScreen';
import QuestionBankScreen from '../screens/mentor/QuestionBankScreen';

type MentorRootStack = MentorStackParamList & { Tabs: undefined };
const Stack = createNativeStackNavigator<MentorRootStack>();

export default function MentorNavigator() {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: Colors.background },
        animation: 'slide_from_right',
      }}
    >
      <Stack.Screen name="Tabs" component={MentorTabNavigator} />
      <Stack.Screen name="MenteeList" component={MenteeListScreen} />
      <Stack.Screen name="AssignTask" component={AssignTaskScreen} />
      <Stack.Screen name="MenteeProgress" component={MenteeProgressScreen} />
      <Stack.Screen name="QuestionBank" component={QuestionBankScreen} />
    </Stack.Navigator>
  );
}