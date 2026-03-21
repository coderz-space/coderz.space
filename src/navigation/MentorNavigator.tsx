import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { MentorStackParamList } from '../types';
import PlaceholderScreen from '../screens/PlaceholderScreen.tsx';
import { Colors } from '../theme';

const Stack = createNativeStackNavigator<MentorStackParamList>();

export default function MentorNavigator() {
  return (
    <Stack.Navigator
      screenOptions={{
        headerStyle: { backgroundColor: Colors.surface },
        headerTintColor: Colors.textPrimary,
        headerTitleStyle: { fontWeight: '700', color: Colors.textPrimary },
        headerShadowVisible: false,
        contentStyle: { backgroundColor: Colors.background },
      }}
    >
      <Stack.Screen
        name="Dashboard"
        component={PlaceholderScreen}
        options={{ title: 'Mentor Dashboard', headerShown: false }}
        initialParams={{ title: 'Mentor Dashboard' }}
      />
      <Stack.Screen
        name="MenteeProgress"
        component={PlaceholderScreen}
        options={{ title: 'Mentee Progress' }}
        initialParams={{ menteeId: 'mentee-1', title: 'Mentee Progress' }}
      />
      <Stack.Screen
        name="AssignTask"
        component={PlaceholderScreen}
        options={{ title: 'Assign Task' }}
        initialParams={{ menteeId: 'mentee-1', title: 'Assign Task' }}
      />
      <Stack.Screen
        name="QuestionBank"
        component={PlaceholderScreen}
        options={{ title: 'Question Bank' }}
        initialParams={{ title: 'Question Bank' }}
      />
    </Stack.Navigator>
  );
}
