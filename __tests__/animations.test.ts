import { SpringPresets, TimingPresets } from '../src/utils/animations';

describe('Animation Presets', () => {
  test('SpringPresets should have correct configurations', () => {
    expect(SpringPresets.snappy).toEqual({ damping: 18, stiffness: 200, mass: 0.8 });
    expect(SpringPresets.bouncy).toEqual({ damping: 10, stiffness: 120, mass: 1 });
    expect(SpringPresets.gentle).toEqual({ damping: 20, stiffness: 80, mass: 1 });
    expect(SpringPresets.stiff).toEqual({ damping: 25, stiffness: 300, mass: 0.6 });
  });

  test('TimingPresets should have correct configurations', () => {
    expect(TimingPresets.fast.duration).toBe(150);
    expect(TimingPresets.normal.duration).toBe(250);
    expect(TimingPresets.slow.duration).toBe(400);

    expect(TimingPresets.fast.easing).toBeDefined();
    expect(TimingPresets.normal.easing).toBeDefined();
    expect(TimingPresets.slow.easing).toBeDefined();
  });
});
