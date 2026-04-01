describe('toast utility', () => {
  let toast: typeof import('../src/utils/toast').toast;
  let mockListener: jest.Mock;

  beforeEach(async () => {
    jest.resetModules();
    const mod = await import('../src/utils/toast');
    toast = mod.toast;
    mockListener = jest.fn();
  });

  it('subscribes and receives a show event', () => {
    const unsubscribe = toast.subscribe(mockListener);
    toast.show('Hello World', 'info', 2000);

    expect(mockListener).toHaveBeenCalledTimes(1);
    expect(mockListener).toHaveBeenCalledWith(
      expect.objectContaining({ message: 'Hello World', type: 'info', duration: 2000, id: expect.any(String) })
    );
    unsubscribe();
  });

  it('helper methods work correctly', () => {
    const unsubscribe = toast.subscribe(mockListener);

    toast.success('Success message');
    expect(mockListener).toHaveBeenLastCalledWith(expect.objectContaining({ type: 'success' }));

    toast.error('Error message');
    expect(mockListener).toHaveBeenLastCalledWith(expect.objectContaining({ type: 'error' }));

    toast.warning('Warning message');
    expect(mockListener).toHaveBeenLastCalledWith(expect.objectContaining({ type: 'warning' }));

    toast.info('Info message');
    expect(mockListener).toHaveBeenLastCalledWith(expect.objectContaining({ type: 'info' }));
    unsubscribe();
  });

  it('unsubscribes correctly', () => {
    const unsubscribe = toast.subscribe(mockListener);
    unsubscribe();
    toast.show('Test');
    expect(mockListener).not.toHaveBeenCalled();
  });
});
