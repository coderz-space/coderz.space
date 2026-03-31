import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  Modal,
  FlatList,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native';
import { Colors, Typography, Spacing, BorderRadius, Shadow } from '../../theme';
import Button from '../atoms/Button';
import { Problem } from '../../types';
import { getQuestionBanks, getProblemsByBank } from '../../services/api/mock/banksMock.ts'; // we'll create this mock

interface QuestionBankSelectorProps {
  visible: boolean;
  onClose: () => void;
  onSelectProblems: (problems: Problem[]) => void;
}

interface Bank {
  id: string;
  name: string;
}

export default function QuestionBankSelector({
  visible,
  onClose,
  onSelectProblems,
}: QuestionBankSelectorProps) {
  const [banks, setBanks] = useState<Bank[]>([]);
  const [selectedBankId, setSelectedBankId] = useState<string | null>(null);
  const [problems, setProblems] = useState<Problem[]>([]);
  const [selectedProblemIds, setSelectedProblemIds] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(false);
  const [step, setStep] = useState<'banks' | 'problems'>('banks');

  // Load banks on mount
  useEffect(() => {
    if (visible) {
      setBanks(getQuestionBanks());
      setStep('banks');
      setSelectedBankId(null);
      setProblems([]);
      setSelectedProblemIds(new Set());
    }
  }, [visible]);

  const handleBankSelect = async (bankId: string) => {
    setSelectedBankId(bankId);
    setLoading(true);
    try {
      const problemsData = await getProblemsByBank(bankId);
      setProblems(problemsData);
      setStep('problems');
    } catch (error) {
      console.error('Failed to load problems', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleProblem = (problemId: string) => {
    setSelectedProblemIds(prev => {
      const newSet = new Set(prev);
      if (newSet.has(problemId)) {
        newSet.delete(problemId);
      } else {
        newSet.add(problemId);
      }
      return newSet;
    });
  };

  const handleDone = () => {
    const selectedProblems = problems.filter(p => selectedProblemIds.has(p.id));
    onSelectProblems(selectedProblems);
    onClose();
  };

  return (
    <Modal
      animationType="slide"
      transparent
      visible={visible}
      onRequestClose={onClose}
    >
      <View style={styles.modalOverlay}>
        <View style={styles.modalContent}>
          {/* Header */}
          <View style={styles.header}>
            <Text style={styles.title}>
              {step === 'banks' ? 'Select Question Bank' : 'Select Problems'}
            </Text>
            <TouchableOpacity onPress={onClose}>
              <Text style={styles.closeBtn}>✕</Text>
            </TouchableOpacity>
          </View>

          {step === 'banks' && (
            <FlatList
              data={banks}
              keyExtractor={item => item.id}
              renderItem={({ item }) => (
                <TouchableOpacity
                  style={styles.bankItem}
                  onPress={() => handleBankSelect(item.id)}
                >
                  <Text style={styles.bankName}>{item.name}</Text>
                  <Text style={styles.arrowRight}>→</Text>
                </TouchableOpacity>
              )}
            />
          )}

          {step === 'problems' && (
            <>
              {loading ? (
                <ActivityIndicator size="large" color={Colors.primary} style={styles.loader} />
              ) : (
                <>
                  <FlatList
                    data={problems}
                    keyExtractor={item => item.id}
                    renderItem={({ item }) => (
                      <TouchableOpacity
                        style={styles.problemItem}
                        onPress={() => toggleProblem(item.id)}
                      >
                        <View style={styles.checkbox}>
                          {selectedProblemIds.has(item.id) && (
                            <Text style={styles.checkmark}>✓</Text>
                          )}
                        </View>
                        <Text style={styles.problemTitle}>{item.title}</Text>
                      </TouchableOpacity>
                    )}
                  />
                  <Button
                    label={`Done (${selectedProblemIds.size} selected)`}
                    onPress={handleDone}
                    fullWidth
                    size="md"
                    style={styles.doneBtn}
                  />
                </>
              )}
            </>
          )}
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  modalOverlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.6)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  modalContent: {
    backgroundColor: Colors.surface,
    borderRadius: BorderRadius['2xl'],
    width: '90%',
    maxHeight: '80%',
    padding: Spacing.lg,
    ...Shadow.lg,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: Spacing.lg,
    paddingBottom: Spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: Colors.divider,
  },
  title: {
    ...Typography.headingMedium,
    color: Colors.textPrimary,
  },
  closeBtn: {
    ...Typography.bodyLarge,
    color: Colors.textSecondary,
    padding: Spacing.sm,
  },
  bankItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: Spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: Colors.divider,
  },
  bankName: {
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
  },
  arrowRight: {
    ...Typography.bodyMedium,
    color: Colors.textSecondary,
  },
  loader: {
    marginVertical: Spacing.xl,
  },
  problemItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: Spacing.sm,
    borderBottomWidth: 1,
    borderBottomColor: Colors.divider,
  },
  checkbox: {
    width: 24,
    height: 24,
    borderRadius: 4,
    borderWidth: 1.5,
    borderColor: Colors.primary,
    marginRight: Spacing.md,
    justifyContent: 'center',
    alignItems: 'center',
  },
  checkmark: {
    color: Colors.primary,
    fontSize: 14,
    fontWeight: 'bold',
  },
  problemTitle: {
    ...Typography.bodyMedium,
    color: Colors.textPrimary,
    flex: 1,
  },
  doneBtn: {
    marginTop: Spacing.xl,
  },
});