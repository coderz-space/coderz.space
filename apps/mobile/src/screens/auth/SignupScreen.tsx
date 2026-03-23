import React, { useState } from 'react';
import {
  View, Text, StyleSheet, TouchableOpacity, Alert,
} from 'react-native';
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withSpring,
  withDelay,
} from 'react-native-reanimated';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Button from '../../components/atoms/Button';
import Input from '../../components/atoms/Input';
import { Colors, Typography, Spacing, BorderRadius } from '../../theme';
import { authMock } from '../../services/api/mock/authMock';
// import { authLive } from '../../services/api/live/authLive'; // swap when ready

interface Props {
  onNavigateToLogin: () => void;
}

export default function SignupScreen({ onNavigateToLogin }: Props) {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  const formScale = useSharedValue(0.96);
  const formOpacity = useSharedValue(0);

  React.useEffect(() => {
    formScale.value = withDelay(100, withSpring(1, { damping: 16, stiffness: 150 }));
    formOpacity.value = withDelay(100, withSpring(1, { damping: 20 }));
  }, []);

  const animStyle = useAnimatedStyle(() => ({
    transform: [{ scale: formScale.value }],
    opacity: formOpacity.value,
  }));

  const validate = (): boolean => {
    const e: Record<string, string> = {};
    if (!name.trim() || name.trim().length < 2) e.name = 'Name must be at least 2 characters';
    if (!email.trim() || !/\S+@\S+\.\S+/.test(email)) e.email = 'Valid email required';
    if (password.length < 8) e.password = 'Min 8 characters required';
    if (!/(?=.*[a-zA-Z])(?=.*[0-9])/.test(password))
      e.password = 'Must include at least 1 letter and 1 number';
    if (password !== confirmPassword) e.confirmPassword = 'Passwords do not match';
    setErrors(e);
    return Object.keys(e).length === 0;
  };

  const handleSignup = async () => {
    if (!validate()) return;
    setLoading(true);
    try {
      // In mock, signup just shows success — real flow goes to email verify
        await new Promise<void>((resolve) => setTimeout(resolve, 1200));
      setSuccess(true);
    } catch (err: any) {
      Alert.alert('Signup Failed', err.message ?? 'Something went wrong');
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <ScreenWrapper padded>
        <View style={styles.successContainer}>
          <View style={styles.successIcon}>
            <Text style={styles.successIconText}>✓</Text>
          </View>
          <Text style={styles.successTitle}>Account Created!</Text>
          <Text style={styles.successSub}>
            Your account is ready. Wait for your mentor to add you to a bootcamp.
          </Text>
          <Button
            label="Go to Login"
            onPress={onNavigateToLogin}
            fullWidth
            size="lg"
            style={styles.successBtn}
          />
        </View>
      </ScreenWrapper>
    );
  }

  return (
    <ScreenWrapper scrollable padded avoidKeyboard>
      <View style={styles.container}>
        {/* Logo */}
        <View style={styles.logoSection}>
          <View style={styles.logoMark}>
            <Text style={styles.logoIcon}>{'</>'}</Text>
          </View>
          <Text style={styles.brandName}>CODERZ SPACE</Text>
          <Text style={styles.tagline}>Create your account</Text>
        </View>

        <Animated.View style={[styles.form, animStyle]}>
          <Text style={styles.formTitle}>Sign Up</Text>

          <Input
            label="Full Name"
            placeholder="Arjun Sharma"
            value={name}
            onChangeText={setName}
            autoCapitalize="words"
            error={errors.name}
          />
          <Input
            label="Email"
            placeholder="you@coderz.space"
            value={email}
            onChangeText={setEmail}
            keyboardType="email-address"
            autoCapitalize="none"
            error={errors.email}
          />
          <Input
            label="Password"
            placeholder="Min 8 chars, 1 letter + 1 number"
            value={password}
            onChangeText={setPassword}
            secureTextEntry
            error={errors.password}
          />
          <Input
            label="Confirm Password"
            placeholder="Repeat your password"
            value={confirmPassword}
            onChangeText={setConfirmPassword}
            secureTextEntry
            error={errors.confirmPassword}
          />

          <Button
            label="Create Account"
            onPress={handleSignup}
            loading={loading}
            fullWidth
            size="lg"
            style={styles.signupBtn}
          />
        </Animated.View>

        <TouchableOpacity onPress={onNavigateToLogin} style={styles.loginLink}>
          <Text style={styles.loginLinkText}>
            Already have an account?{' '}
            <Text style={styles.loginLinkBold}>Sign In</Text>
          </Text>
        </TouchableOpacity>
      </View>
    </ScreenWrapper>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingVertical: Spacing['2xl'],
  },
  logoSection: {
    alignItems: 'center',
    marginBottom: Spacing['2xl'],
  },
  logoMark: {
    width: 64,
    height: 64,
    borderRadius: BorderRadius.xl,
    backgroundColor: Colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: Spacing.md,
    shadowColor: Colors.primary,
    shadowOffset: { width: 0, height: 6 },
    shadowOpacity: 0.4,
    shadowRadius: 14,
    elevation: 8,
  },
  logoIcon: {
    fontSize: 22,
    fontWeight: '800',
    color: Colors.textInverse,
    fontFamily: 'Courier New',
  },
  brandName: {
    ...Typography.headingMedium,
    color: Colors.textPrimary,
    letterSpacing: 3,
    marginBottom: Spacing.xs,
  },
  tagline: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
  form: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius['2xl'],
    padding: Spacing.xl,
    marginBottom: Spacing.lg,
  },
  formTitle: {
    ...Typography.headingMedium,
    color: Colors.textPrimary,
    marginBottom: Spacing.xl,
  },
  signupBtn: {
    marginTop: Spacing.sm,
  },
  loginLink: {
    alignItems: 'center',
    paddingVertical: Spacing.base,
  },
  loginLinkText: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
  loginLinkBold: {
    color: Colors.primary,
    fontWeight: '700',
  },
  successContainer: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: Spacing['2xl'],
  },
  successIcon: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: Colors.successMuted,
    borderWidth: 2,
    borderColor: Colors.success,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: Spacing.xl,
  },
  successIconText: {
    fontSize: 32,
    color: Colors.success,
    fontWeight: '800',
  },
  successTitle: {
    ...Typography.displayMedium,
    color: Colors.textPrimary,
    marginBottom: Spacing.md,
    textAlign: 'center',
  },
  successSub: {
    ...Typography.bodyMedium,
    color: Colors.textSecondary,
    textAlign: 'center',
    marginBottom: Spacing['2xl'],
  },
  successBtn: {},
});