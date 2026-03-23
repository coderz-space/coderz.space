// App.tsx
import React from 'react';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import RootNavigator from './src/navigation/RootNavigator';
import ToastContainer from './src/components/atoms/Toast';

export default function App() {
  return (
    <SafeAreaProvider>
      <RootNavigator />
      <ToastContainer />
    </SafeAreaProvider>
  );
}