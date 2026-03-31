import React, { useState } from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { AuthStackParamList } from '../types';
import LoginScreen from '../screens/auth/LoginScreen';
import SignupScreen from '../screens/auth/SignupScreen';

const Stack = createNativeStackNavigator<AuthStackParamList & { Signup: undefined }>();

export default function AuthNavigator() {
  return (
    <Stack.Navigator
      screenOptions={{ headerShown: false, animation: 'slide_from_bottom' }}
    >
      <Stack.Screen name="Login" component={LoginScreenWrapper} />
      <Stack.Screen name="Signup" component={SignupScreenWrapper} />
    </Stack.Navigator>
  );
}

// Wrappers to pass navigation callbacks cleanly
function LoginScreenWrapper({ navigation }: any) {
  return <LoginScreen onNavigateToSignup={() => navigation.navigate('Signup')} />;
}
function SignupScreenWrapper({ navigation }: any) {
  return <SignupScreen onNavigateToLogin={() => navigation.navigate('Login')} />;
}