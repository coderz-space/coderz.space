import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { useAuthStore } from '../store/authStore';
import { RootStackParamList } from '../types';
import { Colors } from '../theme';

// Navigators (we'll stub these below)
import MentorNavigator from './MentorNavigator.tsx';
import MenteeNavigator from './MenteeNavigator.tsx';
import AuthNavigator from './AuthNavigator.tsx';

const Stack = createNativeStackNavigator<RootStackParamList>();

export default function RootNavigator() {
  const { isAuthenticated, user } = useAuthStore();

  return (
    <NavigationContainer
      theme={{
        dark: true,
        colors: {
          primary: Colors.primary,
          background: Colors.background,
          card: Colors.surface,
          text: Colors.textPrimary,
          border: Colors.surfaceBorder,
          notification: Colors.primary,
        },
        fonts: {
          regular: { fontFamily: 'System', fontWeight: '400' },
          medium: { fontFamily: 'System', fontWeight: '500' },
          bold: { fontFamily: 'System', fontWeight: '700' },
          heavy: { fontFamily: 'System', fontWeight: '800' },
        },
      }}
    >
      <Stack.Navigator screenOptions={{ headerShown: false }}>
        {!isAuthenticated ? (
          <Stack.Screen name="Auth" component={AuthNavigator} />
        ) : user?.role === 'mentor' ? (
          <Stack.Screen name="MentorApp" component={MentorNavigator} />
        ) : (
          <Stack.Screen name="MenteeApp" component={MenteeNavigator} />
        )}
      </Stack.Navigator>
    </NavigationContainer>
  );
}