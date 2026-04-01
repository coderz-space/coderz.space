import { TextStyle } from 'react-native';

export const FontFamily = {
  // Use system fonts that look great - swap with custom fonts later
  heading: 'System',     // Replace with e.g. 'Rajdhani-Bold' if you add custom fonts
  body: 'System',
  mono: 'Courier New',
} as const;

export const FontSize = {
  xs: 11,
  sm: 13,
  base: 15,
  md: 17,
  lg: 20,
  xl: 24,
  '2xl': 28,
  '3xl': 34,
  '4xl': 42,
} as const;

export const FontWeight = {
  regular: '400' as TextStyle['fontWeight'],
  medium: '500' as TextStyle['fontWeight'],
  semibold: '600' as TextStyle['fontWeight'],
  bold: '700' as TextStyle['fontWeight'],
  extrabold: '800' as TextStyle['fontWeight'],
} as const;

export const LineHeight = {
  tight: 1.2,
  normal: 1.5,
  relaxed: 1.75,
} as const;

export const Typography = {
  displayLarge: {
    fontSize: FontSize['3xl'],
    fontWeight: FontWeight.extrabold,
    letterSpacing: -0.5,
  } as TextStyle,

  displayMedium: {
    fontSize: FontSize['2xl'],
    fontWeight: FontWeight.bold,
    letterSpacing: -0.3,
  } as TextStyle,

  headingLarge: {
    fontSize: FontSize.xl,
    fontWeight: FontWeight.bold,
    letterSpacing: -0.2,
  } as TextStyle,

  headingMedium: {
    fontSize: FontSize.lg,
    fontWeight: FontWeight.semibold,
  } as TextStyle,

  headingSmall: {
    fontSize: FontSize.md,
    fontWeight: FontWeight.semibold,
  } as TextStyle,

  bodyLarge: {
    fontSize: FontSize.md,
    fontWeight: FontWeight.regular,
    lineHeight: FontSize.md * 1.5,
  } as TextStyle,

  bodyMedium: {
    fontSize: FontSize.base,
    fontWeight: FontWeight.regular,
    lineHeight: FontSize.base * 1.5,
  } as TextStyle,

  bodySmall: {
    fontSize: FontSize.sm,
    fontWeight: FontWeight.regular,
    lineHeight: FontSize.sm * 1.5,
  } as TextStyle,

  label: {
    fontSize: FontSize.sm,
    fontWeight: FontWeight.medium,
    letterSpacing: 0.3,
  } as TextStyle,

  labelSmall: {
    fontSize: FontSize.xs,
    fontWeight: FontWeight.semibold,
    letterSpacing: 0.8,
    textTransform: 'uppercase',
  } as TextStyle,

  caption: {
    fontSize: FontSize.xs,
    fontWeight: FontWeight.regular,
  } as TextStyle,

  mono: {
    fontFamily: FontFamily.mono,
    fontSize: FontSize.sm,
  } as TextStyle,
} as const;