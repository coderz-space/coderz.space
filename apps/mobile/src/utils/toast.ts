import { useRef, useCallback } from 'react';
import { Animated } from 'react-native';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface ToastMessage {
  id: string;
  type: ToastType;
  message: string;
  duration?: number;
}

// Simple event emitter for toast — no context needed
type ToastListener = (toast: ToastMessage) => void;
const listeners: ToastListener[] = [];

export const toast = {
  show: (message: string, type: ToastType = 'info', duration = 3000) => {
    const t: ToastMessage = { id: Date.now().toString(), type, message, duration };
    listeners.forEach((l) => l(t));
  },
  success: (message: string) => toast.show(message, 'success'),
  error: (message: string) => toast.show(message, 'error'),
  warning: (message: string) => toast.show(message, 'warning'),
  info: (message: string) => toast.show(message, 'info'),
  subscribe: (listener: ToastListener) => {
    listeners.push(listener);
    return () => {
      const idx = listeners.indexOf(listener);
      if (idx > -1) listeners.splice(idx, 1);
    };
  },
};