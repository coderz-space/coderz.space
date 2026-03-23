import React, { useEffect, useRef, useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  Animated,
  Platform,
} from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { toast, ToastMessage, ToastType } from '../../utils/toast';

const TYPE_CONFIG: Record<ToastType, { bg: string; border: string; icon: string }> = {
  success: { bg: Colors.successMuted, border: Colors.success, icon: '✓' },
  error: { bg: Colors.errorMuted, border: Colors.error, icon: '✕' },
  warning: { bg: Colors.warningMuted, border: Colors.warning, icon: '!' },
  info: { bg: Colors.primaryMuted, border: Colors.primary, icon: 'i' },
};

function ToastItem({
  item,
  onDone,
}: {
  item: ToastMessage;
  onDone: (id: string) => void;
}) {
  const opacity = useRef(new Animated.Value(0)).current;
  const translateY = useRef(new Animated.Value(-20)).current;
  const config = TYPE_CONFIG[item.type];

  useEffect(() => {
    Animated.parallel([
      Animated.spring(opacity, { toValue: 1, useNativeDriver: true, speed: 30 }),
      Animated.spring(translateY, { toValue: 0, useNativeDriver: true, speed: 30 }),
    ]).start();

    const timer = setTimeout(() => {
      Animated.parallel([
        Animated.timing(opacity, { toValue: 0, duration: 300, useNativeDriver: true }),
        Animated.timing(translateY, { toValue: -20, duration: 300, useNativeDriver: true }),
      ]).start(() => onDone(item.id));
    }, item.duration ?? 3000);

    return () => clearTimeout(timer);
  }, []);

  return (
    <Animated.View
      style={[
        styles.toast,
        { backgroundColor: config.bg, borderColor: config.border },
        { opacity, transform: [{ translateY }] },
        Shadow.md,
      ]}
    >
      <View style={[styles.iconCircle, { borderColor: config.border }]}>
        <Text style={[styles.icon, { color: config.border }]}>{config.icon}</Text>
      </View>
      <Text style={styles.message} numberOfLines={2}>
        {item.message}
      </Text>
    </Animated.View>
  );
}

export default function ToastContainer() {
  const insets = useSafeAreaInsets();
  const [toasts, setToasts] = useState<ToastMessage[]>([]);

  useEffect(() => {
    const unsub = toast.subscribe((t) => {
      setToasts((prev) => [...prev.slice(-2), t]); // max 3 toasts
    });
    return unsub;
  }, []);

  const handleDone = (id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  };

  if (toasts.length === 0) return null;

  return (
    <View
      style={[
        styles.container,
        { top: insets.top + (Platform.OS === 'android' ? 16 : 8) },
      ]}
      pointerEvents="none"
    >
      {toasts.map((t) => (
        <ToastItem key={t.id} item={t} onDone={handleDone} />
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    position: 'absolute',
    left: Spacing.base,
    right: Spacing.base,
    zIndex: 9999,
  },
  toast: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 1,
    borderRadius: BorderRadius.lg,
    padding: Spacing.md,
    marginBottom: Spacing.sm,
  },
  iconCircle: {
    width: 24,
    height: 24,
    borderRadius: 12,
    borderWidth: 1.5,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: Spacing.sm,
  },
  icon: {
    fontSize: 12,
    fontWeight: '800',
  },
  message: {
    ...Typography.bodySmall,
    color: Colors.textPrimary,
    flex: 1,
    fontWeight: '500',
  },
});