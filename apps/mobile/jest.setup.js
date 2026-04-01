import mockAsyncStorage from '@react-native-async-storage/async-storage/jest/async-storage-mock';
jest.mock('@react-native-async-storage/async-storage', () => mockAsyncStorage);
import 'react-native-gesture-handler/jestSetup';

jest.mock('react-native-reanimated', () => {
  return {
    useSharedValue: jest.fn(() => ({ value: 0 })),
    useAnimatedStyle: jest.fn(() => ({})),
    withSpring: jest.fn((val) => val),
    withTiming: jest.fn((val) => val),
    createAnimatedComponent: jest.fn((Component) => Component),
    default: {
      View: 'View',
      Text: 'Text',
      ScrollView: 'ScrollView',
    },
    View: 'View',
    Text: 'Text',
    ScrollView: 'ScrollView',
    Extrapolate: { CLAMP: 'clamp' },
    interpolate: jest.fn(),
  };
});
