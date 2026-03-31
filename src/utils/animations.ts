import {
  withSpring,
  withTiming,
  withDelay,
  withSequence,
  Easing,
  SpringConfig,
  TimingConfig,
} from 'react-native-reanimated';
export const SpringPresets = {
  snappy: { damping: 18, stiffness: 200, mass: 0.8 } as SpringConfig,
  bouncy: { damping: 10, stiffness: 120, mass: 1 } as SpringConfig,
  gentle: { damping: 20, stiffness: 80, mass: 1 } as SpringConfig,
  stiff: { damping: 25, stiffness: 300, mass: 0.6 } as SpringConfig,
};
export const TimingPresets = {
  fast: { duration: 150, easing: Easing.out(Easing.quad) } as TimingConfig,
  normal: { duration: 250, easing: Easing.out(Easing.cubic) } as TimingConfig,
  slow: { duration: 400, easing: Easing.inOut(Easing.cubic) } as TimingConfig,
}; // Animate a value in with spring 
export const springIn = (toValue: number, config = SpringPresets.snappy) => withSpring(toValue, config);
// Fade in with delay 
export const fadeInDelayed = (toValue: number, delayMs: number) => withDelay(delayMs, withTiming(toValue, TimingPresets.normal)); // Pulse animation (scale up then back) export const pulse = () => withSequence( withSpring(1.08, SpringPresets.stiff), withSpring(1.0, SpringPresets.snappy), ); // Shake animation for errors export const shake = () => withSequence( withTiming(-8, { duration: 60 }), withTiming(8, { duration: 60 }), withTiming(-6, { duration: 60 }), withTiming(6, { duration: 60 }), withTiming(0, { duration: 60 }), ); // Staggered list item entry export const staggeredEntry = (index: number) => withDelay( index * 60, withSpring(1, SpringPresets.gentle), );
