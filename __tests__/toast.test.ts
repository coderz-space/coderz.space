import { toast } from '../src/utils/toast';

describe('toast utility', () => {
  let mockListener: jest.Mock;

  beforeEach(() => {
    mockListener = jest.fn();
  });

  afterEach(() => {
    // Clear listeners to avoid leaking between tests
    // Since we don't have a public clear method, we rely on unsubscribe
    jest.clearAllMocks();
  });

  it('subscribes and receives a show event', () => {
    const unsubscribe = toast.subscribe(mockListener);
    toast.show('Hello World', 'info', 2000);

    expect(mockListener).toHaveBeenCalledTimes(1);
    expect(mockListener).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Hello World',
        type: 'info',
        duration: 2000,
        id: expect.any(String)
      })
    );

    unsubscribe();
  });

  it('helper methods work correctly', () => {
    const unsubscribe = toast.subscribe(mockListener);

    toast.success('Success message');
    expect(mockListener).toHaveBeenLastCalledWith(
      expect.objectContaining({ message: 'Success message', type: 'success' })
    );

    toast.error('Error message');
    expect(mockListener).toHaveBeenLastCalledWith(
      expect.objectContaining({ message: 'Error message', type: 'error' })
    );

    toast.warning('Warning message');
    expect(mockListener).toHaveBeenLastCalledWith(
      expect.objectContaining({ message: 'Warning message', type: 'warning' })
    );

    toast.info('Info message');
    expect(mockListener).toHaveBeenLastCalledWith(
      expect.objectContaining({ message: 'Info message', type: 'info' })
    );

    unsubscribe();
  });

  it('unsubscribes correctly', () => {
    const unsubscribe = toast.subscribe(mockListener);
    unsubscribe();

    toast.show('Test');
    expect(mockListener).not.toHaveBeenCalled();
  });
});
