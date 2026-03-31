import React, { useState } from 'react';
import {
  View, Text, StyleSheet, Animated, TouchableOpacity, Alert,
} from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';
import ScreenWrapper from '../../components/layout/ScreenWrapper';
import Button from '../../components/atoms/Button';
import Input from '../../components/atoms/Input';
import { Colors, Typography, Spacing, BorderRadius } from '../../theme';
import { authMock } from '../../services/api/mock/authMock';
import { useAuthStore } from '../../store/authStore';
import { AppSession } from '../../types';

// ✅ NEW: Props added
interface Props {
  onNavigateToSignup?: () => void;
}

// Enrollment IDs are hardcoded for mock; real values come from getMe() response
const MOCK_ENROLLMENT_IDS: Record<string, string> = {
  'user-mentee-1': 'enrollment-mentee-1',
  'user-mentor-1': 'enrollment-mentor-1',
};
const MOCK_ORG_MEMBER_IDS: Record<string, string> = {
  'user-mentee-1': 'orgmember-mentee-1',
  'user-mentor-1': 'orgmember-mentor-1',
};

export default function LoginScreen({ onNavigateToSignup }: Props) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [emailError, setEmailError] = useState('');
  const [passwordError, setPasswordError] = useState('');

  const { setSession, isLoading, setLoading } = useAuthStore();

  const validate = (): boolean => {
    let valid = true;
    setEmailError('');
    setPasswordError('');
    if (!email.trim()) { setEmailError('Email is required'); valid = false; }
    else if (!/\S+@\S+\.\S+/.test(email)) { setEmailError('Enter a valid email'); valid = false; }
    if (!password) { setPasswordError('Password is required'); valid = false; }
    else if (password.length < 4) { setPasswordError('Password too short'); valid = false; }
    return valid;
  };

  const handleLogin = async () => {
    if (!validate()) return;
    setLoading(true);
    try {
      const res = await authMock.login({ email, password });

      const session: AppSession = {
        user: res.user,
        accessToken: res.accessToken,
        orgRole: res.orgRole,
        bootcampRole: res.bootcampRole,
        activeOrgId: res.activeOrgId ?? 'org-1',
        activeBootcampId: res.activeBootcampId ?? 'bootcamp-1',
        orgMemberId: MOCK_ORG_MEMBER_IDS[res.user.id] ?? 'orgmember-1',
        bootcampEnrollmentId: MOCK_ENROLLMENT_IDS[res.user.id] ?? 'enrollment-1',
      };

      await AsyncStorage.setItem('@session', JSON.stringify(session));
      setSession(session);
    } catch (err: any) {
      Alert.alert('Login Failed', err.message ?? 'Something went wrong');
      setLoading(false);
    }
  };

  return (
    <ScreenWrapper avoidKeyboard padded>
      <View style={styles.container}>
        {/* Logo */}
        <View style={styles.logoSection}>
          <View style={styles.logoMark}>
            <Text style={styles.logoIcon}>{'</>'}</Text>
          </View>
          <Text style={styles.brandName}>CODERZ SPACE</Text>
          <Text style={styles.tagline}>Your bootcamp. Your progress.</Text>
        </View>

        {/* Form */}
        <View style={styles.form}>
          <Text style={styles.formTitle}>Sign In</Text>

          <Input
            label="Email"
            placeholder="you@coderz.space"
            value={email}
            onChangeText={setEmail}
            keyboardType="email-address"
            autoCapitalize="none"
            autoCorrect={false}
            error={emailError}
          />

          <Input
            label="Password"
            placeholder="Enter your password"
            value={password}
            onChangeText={setPassword}
            secureTextEntry={!showPassword}
            error={passwordError}
            rightIcon={
              <Text style={styles.eyeIcon}>{showPassword ? '🙈' : '👁️'}</Text>
            }
            onRightIconPress={() => setShowPassword((v) => !v)}
          />

          <Button
            label="Sign In"
            onPress={handleLogin}
            loading={isLoading}
            fullWidth
            size="lg"
            style={styles.loginBtn}
          />
        </View>

        {/* Hint */}
        <View style={styles.hint}>
          <Text style={styles.hintText}>
            Try{' '}
            <Text
              style={styles.hintLink}
              onPress={() => { setEmail('mentor@test.com'); setPassword('pass123'); }}
            >
              mentor@test.com
            </Text>
            {' or '}
            <Text
              style={styles.hintLink}
              onPress={() => { setEmail('mentee@test.com'); setPassword('pass123'); }}
            >
              mentee@test.com
            </Text>
          </Text>
        </View>

        {/* ✅ NEW: Signup Link */}
        <TouchableOpacity onPress={onNavigateToSignup} style={styles.signupLink}>
          <Text style={styles.signupLinkText}>
            New here?{' '}
            <Text style={styles.signupLinkBold}>Create an account</Text>
          </Text>
        </TouchableOpacity>
      </View>
    </ScreenWrapper>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    paddingVertical: Spacing['2xl'],
  },
  logoSection: {
    alignItems: 'center',
    marginBottom: Spacing['3xl'],
  },
  logoMark: {
    width: 72,
    height: 72,
    borderRadius: BorderRadius.xl,
    backgroundColor: Colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: Spacing.base,
    shadowColor: Colors.primary,
    shadowOffset: { width: 0, height: 6 },
    shadowOpacity: 0.45,
    shadowRadius: 16,
    elevation: 10,
  },
  logoIcon: {
    fontSize: 26,
    fontWeight: '800',
    color: Colors.textInverse,
    fontFamily: 'Courier New',
  },
  brandName: {
    ...Typography.headingLarge,
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
  loginBtn: {
    marginTop: Spacing.sm,
  },
  eyeIcon: {
    fontSize: 16,
  },
  hint: {
    alignItems: 'center',
    paddingHorizontal: Spacing.base,
  },
  hintText: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
    textAlign: 'center',
  },
  hintLink: {
    color: Colors.primary,
    fontWeight: '600',
  },

  // ✅ NEW STYLES
  signupLink: {
    alignItems: 'center',
    paddingTop: Spacing.base,
  },
  signupLinkText: {
    ...Typography.bodySmall,
    color: Colors.textSecondary,
  },
  signupLinkBold: {
    color: Colors.primary,
    fontWeight: '700',
  },
});