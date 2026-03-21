// App.tsx
import React from 'react';
import { SafeAreaProvider } from 'react-native-safe-area-context';

// For testing Phase 1 & 2, temporarily show the gallery:
//import ComponentGalleryScreen from './src/screens/dev/ComponentGalleryScreen';

// For production wiring, use this instead:
import RootNavigator from './src/navigation/RootNavigator';

export default function App() {
  return (
    <SafeAreaProvider>
      <RootNavigator/>
      {/* Switch to <RootNavigator /> when done testing */}
    </SafeAreaProvider>
  );
}