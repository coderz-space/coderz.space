export const Colors = {
  // Backgrounds
  background: '#121212',
  surface: '#1E1E1E',
  surfaceElevated: '#2A2A2A',
  surfaceBorder: '#333333',

  // Orange Accent Palette
  primary: '#FF6B00',
  primaryLight: '#FF8C38',
  primaryDark: '#CC5500',
  primaryMuted: '#FF6B0020', // 12% opacity orange for subtle tints

  // Text
  textPrimary: '#F0F0F0',
  textSecondary: '#9A9A9A',
  textDisabled: '#555555',
  textInverse: '#121212',

  // Status Colors
  success: '#4CAF50',
  successMuted: '#4CAF5020',
  warning: '#FFC107',
  warningMuted: '#FFC10720',
  error: '#F44336',
  errorMuted: '#F4433620',
  info: '#2196F3',

  // Badge specific
  badgeDoubt: '#F44336',
  badgeCompleted: '#4CAF50',
  badgePending: '#FF6B00',
  badgeReview: '#2196F3',

  // Utility
  white: '#FFFFFF',
  black: '#000000',
  transparent: 'transparent',

  // Dividers
  divider: '#2C2C2C',
} as const;

export type ColorKeys = keyof typeof Colors;