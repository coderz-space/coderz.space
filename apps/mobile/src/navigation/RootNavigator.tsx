import React, { useEffect } from 'react';
import { View, ActivityIndicator } from 'react-native';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { useAuthStore } from '../store/authStore';
import { RootStackParamList } from '../types';
import { Colors } from '../theme';

import MentorNavigator from './MentorNavigator';
import MenteeNavigator from './MenteeNavigator';
import AuthNavigator from './AuthNavigator';

const Stack = createNativeStackNavigator<RootStackParamList>();

export default function RootNavigator() {
  const { isAuthenticated, session, isBootstrapping, bootstrapAuth } = useAuthStore();

  useEffect(() => {
    bootstrapAuth();
  }, []);

  if (isBootstrapping) {
    return (
      <View style={{ flex: 1, backgroundColor: Colors.background, alignItems: 'center', justifyContent: 'center' }}>
        <ActivityIndicator color={Colors.primary} size="large" />
      </View>
    );
  }

  const isMentor = session?.bootcampRole === 'mentor' || session?.orgRole === 'mentor' || session?.orgRole === 'admin';

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
        ) : isMentor ? (
          <Stack.Screen name="MentorApp" component={MentorNavigator} />
        ) : (
          <Stack.Screen name="MenteeApp" component={MenteeNavigator} />
        )}
      </Stack.Navigator>
    </NavigationContainer>
  );
}