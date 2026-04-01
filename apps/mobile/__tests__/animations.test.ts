import { withSpring, withTiming, withDelay, withSequence } from 'react-native-reanimated';
import { SpringPresets, TimingPresets, springIn } from '../src/utils/animations';

jest.mock('react-native-reanimated', () => ({
  withSpring: jest.fn((val, config) => ({ type: 'spring', val, config })),
  withTiming: jest.fn((val, config) => ({ type: 'timing', val, config })),
  withDelay: jest.fn((delay, anim) => ({ type: 'delay', delay, anim })),
  withSequence: jest.fn((...anims) => ({ type: 'sequence', anims })),
  Easing: {
    out: jest.fn((e) => `out(${e})`),
    inOut: jest.fn((e) => `inOut(${e})`),
    quad: 'quad',
    cubic: 'cubic',
  },
}));

describe('Animation Presets', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  test('SpringPresets should have correct configurations', () => {
    expect(SpringPresets.snappy).toEqual({ damping: 18, stiffness: 200, mass: 0.8 });
    expect(SpringPresets.bouncy).toEqual({ damping: 10, stiffness: 120, mass: 1 });
    expect(SpringPresets.gentle).toEqual({ damping: 20, stiffness: 80, mass: 1 });
    expect(SpringPresets.stiff).toEqual({ damping: 25, stiffness: 300, mass: 0.6 });
  });

  test('TimingPresets should have correct configurations', () => {
    expect(TimingPresets.fast).toEqual(expect.objectContaining({ duration: 150, easing: expect.any(Function) }));
    expect(TimingPresets.normal).toEqual(expect.objectContaining({ duration: 250, easing: expect.any(Function) }));
    expect(TimingPresets.slow).toEqual(expect.objectContaining({ duration: 400, easing: expect.any(Function) }));
  });
});

describe('Animation Functions', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('springIn', () => {
    test('should call withSpring with default snappy preset', () => {
      const result = springIn(100);
      expect(withSpring).toHaveBeenCalledWith(100, SpringPresets.snappy);
      expect(result).toEqual({ type: 'spring', val: 100, config: SpringPresets.snappy });
    });

    test('should call withSpring with provided config', () => {
      const result = springIn(50, SpringPresets.gentle);
      expect(withSpring).toHaveBeenCalledWith(50, SpringPresets.gentle);
      expect(result).toEqual({ type: 'spring', val: 50, config: SpringPresets.gentle });
    });
  });

  describe('fadeInDelayed', () => {
    test('should call withDelay and withTiming with correct parameters', () => {
      const targetValue = 1;
      const delayMs = 300;

      fadeInDelayed(targetValue, delayMs);

      expect(withTiming).toHaveBeenCalledWith(targetValue, TimingPresets.normal);
      // withTiming mock returns the value itself, so withDelay gets the targetValue as second arg
      expect(withDelay).toHaveBeenCalledWith(delayMs, targetValue);
    });
  });
});
