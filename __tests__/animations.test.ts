import { SpringPresets, TimingPresets } from '../src/utils/animations';

describe('Animation Presets', () => {
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
