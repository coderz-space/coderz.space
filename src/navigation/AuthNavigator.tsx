import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { AuthStackParamList } from '../types';
import PlaceholderScreen from '../screens/PlaceholderScreen.tsx';

const Stack = createNativeStackNavigator<AuthStackParamList>();

export default function AuthNavigator() {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false }}>
      <Stack.Screen
        name="Login"
        component={PlaceholderScreen}
        initialParams={{ title: 'Login Screen' }}
      />
    </Stack.Navigator>
  );
}
