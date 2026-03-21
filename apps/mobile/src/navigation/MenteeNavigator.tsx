import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { MenteeStackParamList } from '../types';
import PlaceholderScreen from '../screens/PlaceholderScreen.tsx';
import { Colors } from '../theme';

const Stack = createNativeStackNavigator<MenteeStackParamList>();

export default function MenteeNavigator() {
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
        options={{ title: 'My Dashboard', headerShown: false }}
        initialParams={{ title: 'Mentee Dashboard' }}
      />
      <Stack.Screen
        name="TaskDetail"
        component={PlaceholderScreen}
        options={{ title: 'Task Detail' }}
        initialParams={{ taskId: 'task-1', title: 'Task Detail' }}
      />
      <Stack.Screen
        name="Leaderboard"
        component={PlaceholderScreen}
        options={{ title: 'Leaderboard' }}
        initialParams={{ title: 'Leaderboard' }}
      />
    </Stack.Navigator>
  );
}
