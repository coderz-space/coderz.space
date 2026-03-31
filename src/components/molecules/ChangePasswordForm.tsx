import React, { useState } from 'react';
import { View, Text, StyleSheet, Alert } from 'react-native';
import Input from '../atoms/Input';
import Button from '../atoms/Button';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import { useAuthStore } from '../../store/authStore';
import { toast } from '../../utils/toast';

export default function ChangePasswordForm() {
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [showCurrent, setShowCurrent] = useState(false);
  const [showNew, setShowNew] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);

  const { changePassword } = useAuthStore();

  const validate = (): boolean => {
    if (!currentPassword) {
      toast.error('Current password is required');
      return false;
    }
    if (newPassword.length < 8) {
      toast.error('New password must be at least 8 characters');
      return false;
    }
    if (!/[A-Za-z]/.test(newPassword) || !/[0-9]/.test(newPassword)) {
      toast.error(
        'New password must contain at least one letter and one number',
      );
      return false;
    }
    if (newPassword !== confirmPassword) {
      toast.error('Passwords do not match');
      return false;
    }
    return true;
  };

  const handleSubmit = async () => {
    if (!validate()) return;
    setLoading(true);
    try {
      await changePassword({ currentPassword, newPassword });
      toast.success('Password changed successfully');
      // Clear form
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (err: any) {
      toast.error(err.message || 'Failed to change password');
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Change Password</Text>
      <Input
        label="Current Password"
        placeholder="Enter current password"
        value={currentPassword}
        onChangeText={setCurrentPassword}
        secureTextEntry={!showCurrent}
        rightIcon={
          <Text style={styles.eyeIcon}>{showCurrent ? '👁️' : '👁️‍🗨️'}</Text>
        }
        onRightIconPress={() => setShowCurrent(!showCurrent)}
      />
      <Input
        label="New Password"
        placeholder="Min 8 chars, 1 letter + 1 number"
        value={newPassword}
        onChangeText={setNewPassword}
        secureTextEntry={!showNew}
        rightIcon={<Text style={styles.eyeIcon}>{showNew ? '👁️' : '👁️‍🗨️'}</Text>}
        onRightIconPress={() => setShowNew(!showNew)}
      />
      <Input
        label="Confirm New Password"
        placeholder="Re-enter new password"
        value={confirmPassword}
        onChangeText={setConfirmPassword}
        secureTextEntry={!showConfirm}
        rightIcon={
          <Text style={styles.eyeIcon}>{showConfirm ? '👁️' : '👁️‍🗨️'}</Text>
        }
        onRightIconPress={() => setShowConfirm(!showConfirm)}
      />
      <Button
        label="Update Password"
        onPress={handleSubmit}
        loading={loading}
        variant="outlined"
        fullWidth
        size="md"
        style={styles.button}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius.xl,
    padding: Spacing.base,
    marginBottom: Spacing.xl,
    ...Shadow.md,
  },
  title: {
    ...Typography.headingSmall,
    color: Colors.textPrimary,
    marginBottom: Spacing.md,
  },
  eyeIcon: {
    fontSize: 16,
  },
  button: {
    marginTop: Spacing.sm,
  },
});
